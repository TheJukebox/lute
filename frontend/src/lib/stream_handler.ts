import '$lib/gen/stream_grpc_web_pb';
import { currentTime, isPlaying, isSeeking } from '$lib/audio_store';

import type { ClientReadableStream } from 'grpc-web';


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
let playing = false;
let seeking = false;
let timeElapsed: number = 0;
let startTime: number = 0;
let volume: number = 1;

// Audio data
let playbackBuffer: AudioBuffer | null = null;
let decodeQueue: Array<Frame> = [];

// Chunks
let chunkBuffer: Uint8Array = new Uint8Array(0);
let workingFrame: Uint8Array = new Uint8Array(0);

/* Stream worker for multithreading */
/* ==================== */
let streamWorker: Worker;
if (typeof window !== 'undefined') {
    streamWorker = new Worker(new URL('$lib/stream_worker.ts', import.meta.url), {type: 'module'});

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
                console.error('Failed to dequeue a frame!');
                break;
            case 'undefined':
                console.error('Worker responded with blank message.');
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
    isPlaying.subscribe((value: boolean) => playing = value);
    currentTime.subscribe((value: number) => timeElapsed = value);
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
 * 
 * @returns 
 */
export function updateCurrentTime(): void {
    if (!context || !playing) return;
    // set the current time here every 1 second.
    currentTime.set(context.currentTime - startTime);
}


export function setVolume(v: number): void {
    volume = v;
    if (currentGain) {
        currentGain.gain.setValueAtTime(v, 0);
    }
}

export async function seek(time: number): Promise<void> {
    if (!context || !playbackBuffer) return;
    console.debug(`Seeking to ${time}`);
    seeking = true;
    isSeeking.set(seeking);

    if (currentNode) {
        currentNode.onended = null;
        currentNode.stop(0);
        currentNode.disconnect();
        clearInterval(timeInterval);
    }

    timeElapsed = time;
    currentTime.set(time);
    seeking = false;
    isSeeking.set(seeking);
    if (playing) playBuffer(timeElapsed);
}


export async function bufferAudio(): Promise<void> {
    if (!context) context = createAudioContext();

    // Use the stream worker to fetch more frames
    let msg: FrameMessage = {type: 'dequeue', frame: undefined};
    while (decodeQueue.length < 5) {
        streamWorker.postMessage(msg);
        await new Promise(resolve => setTimeout(resolve, 100));
    }

    // combine some frames together
    let data: Uint8Array = new Uint8Array(0);
    while (decodeQueue.length > 0) {
        data = concatArrays(data, decodeQueue.shift()?.frame);
    }

    // decode the combined frames
    let audio: AudioBuffer = await context.decodeAudioData(data.buffer);
    playbackBuffer = playbackBuffer ? concatAudioBuffers(playbackBuffer, audio) : audio;
}


async function playBuffer(offset: number = 0): Promise<void> {
    if (!context) context = createAudioContext();
    if (!playbackBuffer) {
        await bufferAudio();
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
    gainNode.gain.setValueAtTime(volume, 0);
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
    console.log("new buff", sourceNode.buffer);
    sourceNode.start(0, offset);
    // preload some more audio
    console.debug(playbackBuffer?.duration);
    if (playbackBuffer.duration < 215.0) bufferAudio();

    sourceNode.onended = () => {
        if (seeking) return;
        if (playing) playBuffer(sourceNode.buffer?.duration);
        else {
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
async function bufferFrame(frame: Uint8Array, seq: number): Promise<void> {
    let msg: FrameMessage = {type: 'queue', frame: {frame, seq}}; 
    streamWorker.postMessage(msg);
}


export function fetchStream(host: string, path: string, sessionId: string): void {   
    const service: proto.stream.AudioStreamClient = new proto.stream.AudioStreamClient(
        host,
        null,
        {},
    );

    const request: proto.stream.AudioStreamRequest = new proto.stream.AudioStreamRequest();
    request.setFileName(path);
    request.setSessionId(sessionId);

    console.info(`(${host}) (${sessionId}) Requestion stream for '${path}'...`);
    const audioStream: ClientReadableStream<proto.stream.AudioStreamChunk> = service.streamAudio(request, null);
    let frameCount: number = 0;

    audioStream.on('data', (response: proto.stream.AudioStreamChunk) => {
        // @ts-ignore: this won't be a string
        let chunk: Uint8Array = response.getData();
        let seq: number = response.getSequence();

        // Append the new chunk to the buffer of chunks.
        chunk = concatArrays(chunkBuffer, chunk);

        let ADTSindex: number = containsADTSHeader(chunk);
        
        // Slice the data just before the detected frame and add it to the previously detected one.
        if (ADTSindex > 0 && workingFrame.length > 0) {
            frameCount += 1;
            workingFrame = concatArrays(workingFrame, chunk.slice(0, ADTSindex));
            bufferFrame(workingFrame, seq);

            workingFrame = chunk.slice(ADTSindex); // set the new frame
        // The new frame is at the start, so we should have a complete frame already buffered.
        } else if (ADTSindex === 0 && workingFrame.length > 0) {
            frameCount += 1;
            bufferFrame(workingFrame, seq);
            workingFrame = chunk;
        // The first frame.
        } else if (ADTSindex === 0) {
            workingFrame = chunk;
        } else if (ADTSindex > 0) {
            workingFrame = chunk.slice(ADTSindex);
        // Finally, we should just throw it in the buffer if it hasn't been handled
        } else {
            chunkBuffer = concatArrays(chunkBuffer, chunk);
        }
    });

    audioStream.on('end', async () => {
        console.info(`(${host}) (${sessionId}) Stream complete.`);
    });
}


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


/**
 * Determines if the data in the array contains a valid ADTS header.
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
            if ((data[i] + 2 & 0x02) >> 1 !== 0x00) {
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




function createAudioContext(): AudioContext {
    console.debug("Creating new AudioContext...");
    context = new AudioContext({
        "sampleRate": 44100,
        "latencyHint": "playback",
    }); 
    startTime = context.currentTime;
    return context;
}