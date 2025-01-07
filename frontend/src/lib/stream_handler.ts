import '$lib/gen/stream_grpc_web_pb';
import type { ClientReadableStream } from 'grpc-web';

// Global audio context
let context: AudioContext | null = null;

// Stream control
let playing: boolean = false;
let currentNode: AudioBufferSourceNode | null = null;
let currentTime: number = 0;

// Audio data
let audioData: Uint8Array = new Uint8Array(0);
let audioBuffer: AudioBuffer | null = null;
let sequence: number = 1;

// Chunks
let chunkBuffer: Uint8Array = new Uint8Array(0);
let frame: Uint8Array = new Uint8Array(0);


function createAudioContext(): AudioContext {
    console.debug("Creating new AudioContext...");
    context = new AudioContext({
        "sampleRate": 44100,
        "latencyHint": "playback",
    });
    return context;
}

async function enqueueAudioData(seq: number, data: Uint8Array): Promise<void> {
    return new Promise(resolve => {
        setInterval(() => {
            if (seq === sequence) {
                sequence += seq + 1;
                audioData = concatArrays(audioData, data);
                resolve();
            }
        });
    });
}

async function dataAvailable(): Promise<void> {
    return new Promise(resolve => {
        setInterval(() => {
            if (audioData.length > 0) resolve();
        });
    });
}

async function decodeBuffer(): Promise<void> {
    if (!context) context = createAudioContext();
    await dataAvailable()
    console.debug("Decoding buffer: ", audioData);
    audioBuffer = await context.decodeAudioData(audioData.subarray().buffer as ArrayBuffer);
    audioData = new Uint8Array(0);
};

async function bufferReady(): Promise<void> {
    return new Promise(resolve => {
        setInterval(() => {
            if (audioBuffer && audioBuffer.length > 0) {
                resolve();
            }
        });
    });
}


async function playBuffer(offset: number = 0): Promise<void> {
    console.log("Waiting for buffer...");
    await bufferReady();
    if (!context) context = createAudioContext();

    if (!audioBuffer) {
        return;
    }
    const next: AudioBuffer = audioBuffer;
    console.log("Playing from buffer: ", next);

    const source: AudioBufferSourceNode = context.createBufferSource();
    const duration = next.duration;

    source.buffer = next;
    source.connect(context.destination);

    source.start(0, 0, duration);
}

export async function togglePlayback(): Promise<void> {
    playing = !playing;
    if (playing) {
        console.debug("Toggling playback on.");
        decodeBuffer();
        playBuffer();
     } else {
        console.debug("Toggling playback off.");
        if (currentNode) {
            if (currentNode.buffer && audioBuffer) {
                audioBuffer = concatAudioBuffers(currentNode.buffer, audioBuffer);
            } else if (currentNode.buffer) {
                audioBuffer = currentNode.buffer;
            }
            currentNode.stop();
        }
     }
}

export async function fetchStream(host: string, path: string, sessionId: string): Promise<void> {   
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

    audioStream.on('data', async (response: proto.stream.AudioStreamChunk) => {
        // @ts-ignore: this won't be a string
        let chunk: Uint8Array = response.getData();
        let seq: number = response.getSequence();

        // Append the new chunk to the buffer of chunks.
        chunk = concatArrays(chunkBuffer, chunk);

        let ADTSindex: number = containsADTSHeader(chunk);
        
        // Slice the data just before the detected frame and add it to the previously detected one.
        if (ADTSindex > 0 && frame.length > 0) {
            const data = concatArrays(frame, chunk.slice(0, ADTSindex));
            await enqueueAudioData(seq, data);
            frame = chunk.slice(ADTSindex); // set the new frame
        // The new frame is at the start, so we should have a complete frame already buffered.
        } else if (ADTSindex === 0 && frame.length > 0) {
            await enqueueAudioData(seq, frame);
            frame = chunk;
        // The first frame.
        } else if (ADTSindex === 0) {
            frame = chunk;
        } else if (ADTSindex > 0) {
            frame = chunk.slice(ADTSindex);
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

    let channels = Math.min(x.numberOfChannels, y.numberOfChannels);
    let length = x.length + y.length;
    const combinedData = context?.createBuffer(channels, length, x.sampleRate);

    for (let i = 0; i < channels; i++) {
        let channel = combinedData?.getChannelData(i);
        channel?.set(x.getChannelData(i), 0);
        channel?.set(y.getChannelData(i), x.length);
    }

    return combinedData;
}


/**
 * Determines if the data in the array contains a valid ADTS header.
 * @param data {Uint8Array} The data to operate on.
 * @returns {number}        The index that the ADTS header begins at. 
 */
function containsADTSHeader(data: Uint8Array): number {
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