import { currentTime, bufferedTime, isPlaying, isSeeking, buffering } from '$lib/audio_store';
import type { AudioStreamRequest } from '$lib/gen/stream_pb';

import { create, toBinary, fromBinary } from "@bufbuild/protobuf";
import { BinaryReader } from "@bufbuild/protobuf/wire";
import { AudioStreamRequestSchema, AudioStreamChunkSchema } from '$lib/gen/stream_pb';
import type { AudioStreamChunk } from '$lib/gen/stream_pb';

import type { ClientReadableStream } from 'grpc-web';
import type { Unsubscriber } from 'svelte/store';
import type { Track } from '$lib/audio_store';

type Frame = {
    frame: Uint8Array;
    seq: number;
};


type FrameMessage = {
    type: string;
    frame: Frame | undefined;
}


/* Globals */
/* ==================== */
// Global audio context
let context: AudioContext | null = null;

// Stream control
let currentNode: AudioBufferSourceNode | null = null;
let currentGain: GainNode | null = null;
let timeInterval: number = 0;
let bufferInterval: number = 0;
let buffTimeInterval: number = 0;
let playing = false;
let seeking = false;
let timeElapsed: number = 0;
let startTime: number = 0;
let volume: number = 1;
let cachedVolume: number = 1;
let unsubTime: Unsubscriber;
let unsubPlay: Unsubscriber;
let trackDuration: number = 0;

// Audio data
let playbackBuffer: AudioBuffer | null = null;
let decodeQueue: Array<Frame> = [];

// Chunks
let chunkBuffer: Uint8Array = new Uint8Array(0);
let workingFrame: Uint8Array = new Uint8Array(0);

/* Stream worker for multithreading */
/* ==================== */
let streamWorker: Worker;
let dequeueFail: number = 0;
if (typeof window !== 'undefined') {
    streamWorker = new Worker(new URL('$lib/stream_worker.ts', import.meta.url), {type: 'module'});

    // Switch on message types
    streamWorker.onmessage = async (event: MessageEvent<FrameMessage>) => {
        const { type, frame } = event.data;
        switch(type) {
            case 'queue_ok':
                break;
            case 'dequeue_ok':
                if (frame) {
                    decodeQueue.push(frame);
                    decodeQueue.sort((a, b) => a.seq - b.seq);
                }
                break;
            case 'dequeue_fail':
                console.warn('Failed to dequeue a frame!');
                dequeueFail += 1;
                break;
            case 'empty':
                console.debug('Frame queue has been emptied.');
                break;
            case 'undefined':
                console.error('Worker responded with blank message.');
                break;
        }
    }
}


/* Audio buffer handling below */
/* ==================== */

/**
 * Toggles playback of audio stored in playbackBuffer.
 * @function
 * @async
 * @returns Promise<void>
 */
export async function togglePlayback(): Promise<void> {
    unsubPlay = isPlaying.subscribe((value: boolean) => playing = value);
    unsubTime = currentTime.subscribe((value: number) => timeElapsed = value);
    if (playing) {
        playBuffer(timeElapsed);
    } else {
        if (context) {
            currentTime.set(context.currentTime - startTime);
        }
        if (currentNode) {
            currentNode.onended = null;
            currentNode.stop(0);
            currentNode.disconnect();
        }
        clearInterval(timeInterval);
        currentNode?.stop(0);
    }
}

/**
 * Updates the buffered time. 
 * @returns void
 */
export function updateBufferedTime(): void {
    if (!context || !playing) return;
    // set the current time here every 1 second.
    if (playbackBuffer) bufferedTime.set(playbackBuffer.duration);
}

/**
 * Updates the elapsed time of playback. 
 * @returns void
 */
export function updateCurrentTime(): void {
    if (!context || !playing) return;
    // set the current time here every 1 second.
    currentTime.set(context.currentTime - startTime);
}

/**
 * Sets the volume of the playback via gain nodes.
 * @param v The new volume.
 */
export function setVolume(v: number): void {
    volume = v;
    if (currentGain) {
        currentGain.gain.setValueAtTime(v, 0);
    }
}

/**
 * Seeks to a certain time in the playback buffer.
 * @param time The time to seek to.
 * @async
 */
export async function seek(time: number): Promise<void> {
    if (!context || !playbackBuffer) return;

    // Flag for avoiding playback when currentNode is stopped.
    seeking = true;
    isSeeking.set(seeking);

    // Kill the source node.
    if (currentNode) {
        currentNode.onended = null;
        currentNode.stop(0);
        currentNode.disconnect();
    }

    // Pause time updates 
    clearInterval(timeInterval);
    timeElapsed = time;
    currentTime.set(time);

    // reset the flag
    seeking = false;
    isSeeking.set(seeking);

    // begin playback again if required
    if (playing) playBuffer(timeElapsed);
}

function bufferComplete(): boolean {
    const tolerance = 1;
    if (!playbackBuffer) return false;
    const bufferDuration: number = Math.round(playbackBuffer.duration);
    // 100 - 100 = 0
    // 101 - 100 = 1
    // 99 - 100 = -1
    return bufferDuration >= trackDuration - tolerance;
}

/**
 * Decodes queued frames and updates the playback buffer with the resulting audio.
 * @function
 * @async
 */
export async function bufferAudio(): Promise<void> {
    if (!context) context = createAudioContext();
    if (bufferComplete()) {
        clearInterval(bufferInterval); 
        clearInterval(buffTimeInterval);
        return;
    }
    if (dequeueFail === 100) {
        console.error("Failed to dequeue a new frame after 100 retries. Cancelling buffering!");
        clearInterval(bufferInterval); 
        clearInterval(buffTimeInterval);
        return;
    }

    // Use the stream worker to fetch more frames
    let msg: FrameMessage = {type: 'dequeue', frame: undefined};
    while (decodeQueue.length < 5) {
        if (bufferComplete()) {
            clearInterval(bufferInterval); 
            clearInterval(buffTimeInterval);
            break;
        }
        streamWorker.postMessage(msg);
        await new Promise(resolve => setTimeout(resolve, 100));
    }

    // combine some frames together
    let data: Uint8Array = new Uint8Array(0);
    while (decodeQueue.length > 0) {
        let frame: Uint8Array | undefined = decodeQueue.shift()?.frame;
        if (frame) {
            data = concatArrays(data, frame);
        }
    }

    // decode the combined frames
    try {
        let audio: AudioBuffer = await context.decodeAudioData(data.buffer as ArrayBuffer);
        playbackBuffer = playbackBuffer ? concatAudioBuffers(playbackBuffer, audio) : audio;
    } catch (error) {
        if (error instanceof DOMException) {
            console.error("Failed to decode an audio buffer... discarding it.");
        }
    }
}

async function bufferReady(): Promise<void> {
    while (!playbackBuffer) {
        await new Promise(resolve => setTimeout(resolve, 1000));
    }
}

/**
 * Recursively plays back the buffered audio.
 * @function
 * @async 
 * @param offset Offsets the time to begin playback in the buffer.
 */
async function playBuffer(offset: number = 0): Promise<void> {
    if (!context) context = createAudioContext();
    if (!playbackBuffer) {
        await bufferReady();
        buffering.set(false);
    }

    // cleanup
    if (currentNode) {
        currentNode.onended = null;
        currentNode.stop(0);
        currentNode.disconnect();
    }

    // source -> gain -> filter
    // filter
    const filterNode = context.createBiquadFilter();
    filterNode.type = 'lowpass';
    filterNode.connect(context.destination);

    // gain
    const gainNode = context.createGain();
    gainNode.connect(filterNode);
    gainNode.gain.setValueAtTime(seeking ? 0 : volume, 0);
    currentGain = gainNode;

    // buffer
    const sourceNode = context.createBufferSource();
    sourceNode.buffer = playbackBuffer;
    sourceNode.connect(gainNode);
    currentNode = sourceNode;

    filterNode.frequency.setValueAtTime(10000, 0);

    // start playback
    startTime = context.currentTime - offset;
    timeInterval = setInterval(updateCurrentTime, 500);
    buffTimeInterval = setInterval(updateBufferedTime, 500);
    sourceNode.start(0, offset);

    sourceNode.onended = () => {
        if (seeking) return;
        if (playing) playBuffer(sourceNode.buffer?.duration);
        else {
    console.log(encodeBinary(AudioStreamRequestSchema, streamRequest));
            clearInterval(timeInterval);
            sourceNode.stop();
            sourceNode.disconnect();
            gainNode.disconnect();
            filterNode.disconnect();
        }
    };
}


/* Frame handling below */
/* ==================== */

/**
 * Enqueues frames via the stream worker.
 * @function
 * @param frame The frame to buffer.
 * @param seq The sequence ID for this frame.
 */
function bufferFrames(frame: Uint8Array, seq: number): Promise<void> {
    let msg: FrameMessage = {type: 'queue', frame: {frame, seq}}; 
    streamWorker.postMessage(msg);
    console.log("Buffered a frame...");
}

/**
 * Convenience function for resetting stream state.
 * @function
 */
export function resetStream(): void {
    playing = false;
    isPlaying.set(false);
    if (currentNode) {
        currentNode.stop(0);
        currentNode.disconnect();
        currentNode = null;
    }
    if (currentGain) {
        currentGain.disconnect();
        currentGain = null;
    }
    if (context) {
        context.close();
    }
    context = null;
    playbackBuffer = null;
    bufferedTime.set(0);
    clearInterval(buffTimeInterval);

    clearInterval(timeInterval);
    timeElapsed = 0;
    startTime = 0;
    currentTime.set(0);
    streamWorker.postMessage({type: 'empty', frame: undefined});
    decodeQueue = []; 
    chunkBuffer = new Uint8Array(0);
    workingFrame = new Uint8Array(0);
    playing = true;
    isPlaying.set(true);
}

/**
 * Fetches a new stream.
 * @function
 * @param host The host/port to receive the stream from.
 * @param track.path The path for the file to stream on the server.
 * @param sessionId A session ID to associate with the stream.
 */
export async function fetchStream(host: string, track: Track, sessionId: string): void {
    resetStream();

    // TODO: consider turning this into a client class or something
    const streamRequest = create(AudioStreamRequestSchema, {
        fileName: track.path,
        sessionId: sessionId,
    });
    const binary = toBinary(AudioStreamRequestSchema, streamRequest);
    console.log(`Created stream request: ${binary}`);

    const response = await fetch(`http://localhost:8080/streamAudio/StreamAudio`, {
        method: 'POST',
        headers: {
            "Content-Type": "application/lute-grpc",
        },
        body: binary,
    });

    console.log("Received response, starting stream...");
    trackDuration = track.duration;
    bufferInterval = setInterval(bufferAudio, 1000);

    const reader = response.body.getReader();
    let grpcBuffer = new Uint8Array(0);
    let adtsBuffer = new Uint8Array(0);
    let audioBuffer = new Uint8Array(0);
    
    // gather chunks
    while(true) {
        const { value, done } = await reader.read();
        if (done && grpcBuffer.length === 0) {
            console.log("Stream finished!");
            break;
        }

        // add the chunk to the grpcBuffer
        let raw = new Uint8Array(value);
        grpcBuffer = concatArrays(grpcBuffer, raw);

        let grpcLength = grpcFrameLength(grpcBuffer);
        if (grpcLength < 1 || grpcBuffer.length < 5 + grpcLength) continue;

        // add adts frame to the adts buffer and remove it from the grpc buffer
        let chunk = fromBinary(AudioStreamChunkSchema, stripGrpcFrame(grpcBuffer.slice(0, 5 + grpcLength)));
        console.log(chunk);
        grpcBuffer = grpcBuffer.slice(5 + grpcLength);
        adtsBuffer = concatArrays(adtsBuffer, chunk.data);

        // check for complete adts frames
        while (true) {
            let adtsLength = adtsFrameLength(adtsBuffer);
            if (adtsLength < 1 || adtsLength > adtsBuffer.length) break;
            let audioChunk: Uint8Array = adtsBuffer.slice(0, adtsLength);
            audioBuffer = concatArrays(audioBuffer, audioChunk);
            adtsBuffer = adtsBuffer.slice(adtsLength);
        }
        bufferFrames(audioBuffer);
        audioBuffer = new Uint8Array(0);
    }
}

export function grpcFrameLength(data: Uint8Array): number {
    if (data.length < 5) return -1; // too small to be complete
    const length =
        (data[1] << 24) | // the top byte into MSB
        (data[2] << 16) | // the next MSB
        (data[3] << 8)  | // and the next
        data[4];          // this is the LSB
    if (data.length >= 5 + length) {
        return length;
    } else {
        return -1;
    }
}

export function stripGrpcFrame(data: Uinut8Array): Uint8Array {
    // create the length integer by bitshifting into
    // a big-endian 32 bit integer
    const length =
        (data[1] << 24) | // the top byte into MSB
        (data[2] << 16) | // the next MSB
        (data[3] << 8)  | // and the next
        data[4];          // this is the LSB

    // return the contents
    return data.slice(5, 5 + length);
}

/**
 * Concatenates two AudioBuffers into a single AudioBuffer.
 * @function  
 * @param x The audio buffer to append to
 * @param y The audio buffer to append
 * @returns AudioBuffer
 */
function concatAudioBuffers(x: AudioBuffer, y: AudioBuffer): AudioBuffer | null {
    if (!context) {
        console.error("Unable to concatenate AudioBuffers as the AudioContext was not initialised!");
        return null;
    };

    let channels = Math.min(x.numberOfChannels, y.numberOfChannels); // use the smaller of both to avoid overflow
    let length = x.length + y.length;
    const combinedData: AudioBuffer = context.createBuffer(channels, length, x.sampleRate);

    for (let i = 0; i < channels; i++) {
        let channel = combinedData.getChannelData(i);
        channel.set(x.getChannelData(i), 0);
        channel.set(y.getChannelData(i), x.length);
    }

    return combinedData;
}

function adtsFrameLength(data: Uint8Array): number {
    if (data.length < 7) return -1;
    const length =
        ((data[3] & 0x03) << 11 |   // lower 2 bits of byte 3
        (data[4] << 3) |            // all 8 bits of byte 4
        ((data[5] & 0xE0) >> 5));    // the top 3 bits of byte 5
    return length >= 7 && length <= data.length ? length : -1;
}

/**
 * Determines if the data in the array contains a valid ADTS header.
 * @function
 * @param data {Uint8Array} The data to operate on.
 * @returns {number}        The index that the ADTS header begins at. 
 */
function containsADTSHeader(data: Uint8Array): number { // abaduser: potentially add a bool here to track chunks with no ADTS header?
    for (let i = 0; i < data.length - 7; i++) {
        // Identify the sync word - 12 bits
        if (data[i] === 0xFF && (data[i + 1] & 0xF0) === 0xF0) {
            // Check the ID (MPEG version) bits - 2 bits
            // Should be 0 for MPEG-4.
            if ((data[i + 1] & 0x08) >> 3 !== 0x00) {
                continue;
            }
            // layer bit - 1 bit
            // Always 0 for AAC. We don't use layers.
            if ((data[i + 1] & 0x06) >> 1 !== 0x00) {
                continue;
            }
            // protection_absent bit - 1 bit
            // 0 = error protection present (CRC check)
            // 1 = error protection absent
            if ((data[i + 1] & 0x01) !== 0x01) {
                // i think most of the time it should be present, so...
                continue;
                // that said, make sure we revisit this
            }
            // profile - 2 bits
            // expecting 01, for AAC-LC (low-complexity)
            if ((data[i + 2] & 0xC0) >> 6 !== 0x01) {
                continue;
            }
            // sampling_frequency - 4 bits
            // expecting this to be 0x04 (0100)
            if ((data[i + 2] & 0x3C) >> 2 !== 0x04) {
                continue;
            }
            // privacy bit - 1 bit
            // should always be 0
            if ((data[i + 2] & 0x02) >> 1 !== 0x00) {
                continue;
            }
            // channel configuration - 4 bits 
            // i expect this to be 0x04 (0100), but it's split across 4 bits.
            // this means the audio is stereo.
            if ((data[i + 2] & 0x01) !== 0x00 || (data[i + 3] & 0xE0) >> 5 !== 0x04) {
                continue;
            }
            // copy bit - 1 bit
            // should be 1 if it's a copy or 0 if it isn't
            if ((data[i + 3] & 0x20) >> 5 === 1) {
                continue;
            }
            // home bit, should be 0.
            if ((data[i + 3] & 0x20) >> 5 !== 0) {
                continue;
            }
            return i;
        }
    }
    return -1;
}


/**
 * Concatenates two byte arrays and returns the result.
 * @function
 * @param x {Uint8Array}    The first array.
 * @param y {Uint8Array}    The second array, to be appended.
 * @returns {Uint8Array}    An array with x and y's data combined.
 */
function concatArrays(x: Uint8Array, y: Uint8Array): Uint8Array {
    const combinedData: Uint8Array = new Uint8Array(x.length + y.length);
    combinedData.set(x);
    combinedData.set(y, x.length);
    return combinedData;
}


/**
 * Creates and sets the global audio context for the stream.
 * @function
 * @returns AudioContext
 */
function createAudioContext(): AudioContext {
    console.debug("Creating new AudioContext...");
    const newContext = new AudioContext({
        "sampleRate": 44100,
        "latencyHint": "playback",
    }); 
    return newContext;
}
