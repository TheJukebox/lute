import '$lib/gen/stream_grpc_web_pb';
import { currentTime } from '../audio_store';
import stream from '$lib/gen/stream_grpc_web_pb';

import type { ClientReadableStream } from 'grpc-web';


type Frame = {
    frame: Uint8Array;
    seq: number;
};

let streamWorker: Worker;
if (typeof window !== 'undefined') {
    let nextSeq = 1;
    let framesBuffer: Array<{frames: ArrayBuffer, seq: number}> = [];
    streamWorker = new Worker(new URL('$lib/stream_worker.ts', import.meta.url), {type: 'module'});

    streamWorker.onmessage = async (event: MessageEvent<{frames: Uint8Array, seq: number}>) => {
        let frames: ArrayBuffer = event.data.frames.buffer as ArrayBuffer;
        let seq: number = event.data.seq;
        framesBuffer.push({frames: frames, seq: seq});
        framesBuffer.sort((a, b) => a.seq - b.seq);
        let x = framesBuffer.shift();
        decodeBuffer(x?.frames, x?.seq);
    }
}


// Global audio context
let context: AudioContext | null = null;

// Stream control
let playing: boolean = false;
let currentNode: AudioBufferSourceNode | null = null;
let currentGain: GainNode | null = null;
let streamIntervalId: number = 0;
let timeIntervalId: number = 0;

// Audio data
let playbackBuffer: AudioBuffer;

// Chunks
let chunkBuffer: Uint8Array = new Uint8Array(0);
let workingFrame: Uint8Array = new Uint8Array(0);

function createAudioContext(): AudioContext {
    console.debug("Creating new AudioContext...");
    context = new AudioContext({
        "sampleRate": 44100,
        "latencyHint": "playback",
    });
    playbackBuffer = context.createBuffer(2, 1, 44100);
    return context;
}

export async function updateCurrentTime(): Promise<void> {
    if (!context) return;
    currentTime.set(context.currentTime);
}

async function playBuffer(offset: number = 0): Promise<void> {
    if (!context) context = createAudioContext();

    clearInterval(streamIntervalId);

    const source: AudioBufferSourceNode = context.createBufferSource();
    const gainNode = context.createGain();
    const filter = context.createBiquadFilter();

    filter.type = 'lowpass';
    source.buffer = playbackBuffer;
    source.connect(gainNode);
    gainNode.connect(filter);
    filter.connect(context.destination);

    currentNode = source;
    currentGain = gainNode;

    filter.frequency.setValueAtTime(5000, 0);
    filter.frequency.setValueAtTime(10000, 0.1);

    console.debug(`Playing @ ${offset}: `, source.buffer);

    source.start(0, offset);

    source.onended = () => {
        if (playing && source.buffer) playBuffer(source.buffer.duration);
        currentTime.set(context ? context.currentTime : 0);
        source.disconnect();
        gainNode.disconnect(); 
    };
}


export async function togglePlayback(): Promise<void> {
    playing = !playing;
    let time = 0;
    currentTime.subscribe((value: number) => time = value);
    if (playing) {
        streamIntervalId = setInterval(playBuffer, 1, time);
        timeIntervalId = setInterval(updateCurrentTime, 1000);
    } else {
        currentGain?.gain.setTargetAtTime(0, time - 1, 1);
        currentNode?.stop(1);
    }
}


async function bufferFrame(frame: Uint8Array, seq: number): Promise<void> {
    let msg: Frame = {frame: frame, seq: seq}; 
    streamWorker.postMessage(msg);
}


async function decodeBuffer(data: ArrayBuffer, seq: number): Promise<void> {
    if (!context) context = createAudioContext();
    const audio: AudioBuffer = await context.decodeAudioData(data);
    playbackBuffer = concatAudioBuffers(playbackBuffer, audio);
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

    audioStream.on('data', (response: proto.stream.AudioStreamChunk) => {
        // @ts-ignore: this won't be a string
        let chunk: Uint8Array = response.getData();
        let seq: number = response.getSequence();

        // Append the new chunk to the buffer of chunks.
        chunk = concatArrays(chunkBuffer, chunk);

        let ADTSindex: number = containsADTSHeader(chunk);
        
        // Slice the data just before the detected frame and add it to the previously detected one.
        if (ADTSindex > 0 && workingFrame.length > 0) {
            workingFrame = concatArrays(workingFrame, chunk.slice(0, ADTSindex));
            bufferFrame(workingFrame, seq);

            workingFrame = chunk.slice(ADTSindex); // set the new frame
        // The new frame is at the start, so we should have a complete frame already buffered.
        } else if (ADTSindex === 0 && workingFrame.length > 0) {
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