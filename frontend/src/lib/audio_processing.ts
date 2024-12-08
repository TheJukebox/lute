import '$lib/gen/stream_grpc_web_pb';

let context: AudioContext | null = null;

// Buffered audio
const audioBuffer: AudioBuffer[] = [];
let nextBuffer: number = 1;
let playing: boolean = false;
let currentNode: AudioBufferSourceNode | null = null;

// Raw chunk queueing
let lastFrame: Uint8Array = new Uint8Array(0);
let chunkQueue: Uint8Array = new Uint8Array(0);


function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

async function playbackReady() {
    return new Promise(resolve => {
        setInterval(() => {
            if (audioBuffer.length > 0) {
                resolve(true);
            }
        }, 1);
    });

}

export async function togglePlayback(): Promise<void> {
    playing = !playing;
    if (playing) {
        console.log("Stream playing.");
        await playbackReady();
        playFromBuffer();
    } else {
        if (currentNode) {
            currentNode.stop();
        }
        console.log("Stream paused.");
    }
}

function playFromBuffer() {
    // if we aren't meant to be playing, bail out
    if (!playing) { 
        // we need to handle stopping better, otherwise we'll lose
        // a tiny bit of time whenever we pause.
        return;
    }

    if (context === null) {
        console.log("CREATED NEW CONTEXT");
        // potentially need to set sampleRate option: https://webaudio.github.io/web-audio-api/#dom-baseaudiocontext-samplerate
        // by default it should use the target device's preferred sample rate
        context = new AudioContext({"sampleRate": 44100, "latencyHint": "balanced"}); 
    }
    const source: AudioBufferSourceNode = context.createBufferSource();
    currentNode = source;
    const next: AudioBuffer | null | undefined = audioBuffer.shift();
    if (next === null || next === undefined) {
        return;
    }
    
    source.buffer = next; 
    // @ts-ignore - context will not be null, even though ts thinks it could be
    source.connect(context.destination);
    source.start();

    source.onended = () => {
        source.stop();
        source.disconnect();
        playFromBuffer();
    }
}


async function decodeAudio(data: Uint8Array) : Promise<AudioBuffer> {
    // we defer the creation of the AudioContext to avoid "autoplay" policy.
    if (!context) {
        console.log("CREATED NEW CONTEXT");
        // potentially need to set sampleRate option: https://webaudio.github.io/web-audio-api/#dom-baseaudiocontext-samplerate
        // by default it should use the target device's preferred sample rate
        context = new AudioContext({"sampleRate": 44100, "latencyHint": "balanced"}); 
    }
    const buff: ArrayBuffer = data.buffer as ArrayBuffer;
    return await context.decodeAudioData(buff);
}


function concatArrays(x: Uint8Array, y: Uint8Array) {
    const combinedData: Uint8Array = new Uint8Array(x.length + y.length);
    combinedData.set(x);
    combinedData.set(y, x.length);
    return combinedData;
}


async function waitForPrevious(seq: number) {
    return new Promise(resolve => {
        setInterval(() => {
            if (seq === nextBuffer) {
                nextBuffer = seq + 1;
                resolve(true);
            }
        }, 10);
    });
}


async function queueAudio(chunk: Uint8Array, seq: number): Promise<void> {
    let decoded: AudioBuffer = await decodeAudio(chunk);
    await waitForPrevious(seq);
    audioBuffer.push(decoded);
}


export function fetchStream(host: string, path: string, sessionId: string) {
    let service = new proto.stream.AudioStreamClient(
        host,
        null,
        {
            'use-fetch': true,
        }
    );
    let request = new proto.stream.AudioStreamRequest();
    request.setFileName(path);
    request.setSessionId(sessionId);

    console.info(`(${host}) (${sessionId}) Requesting stream '${path}'...`);
    const audioStream = service.streamAudio(request, null);

    audioStream.on('data', async (response: proto.stream.AudioStreamChunk) => {
        // @ts-ignore - for some reason it thinks it's going to be a string
        // after assignment...
        let chunk: Uint8Array = response.getData(); 
        let seq: number = response.getSequence();

        chunk = concatArrays(chunkQueue, chunk);

        let i: number = containsADTSHeader(chunk);
        if (i > 0 && lastFrame.length > 0) {
            const audio = concatArrays(lastFrame, chunk.slice(0,i));
            try {
                queueAudio(audio, seq);
            } catch (err) {
                console.debug("\tFailed to decode the frame!");
            }
            lastFrame = chunk.slice(i);
        } else if (i === 0 && lastFrame.length > 0) {
            try {
                queueAudio(lastFrame, seq);
            } catch (err) {
                console.debug("\tFailed to decode the frame!");
            }
            lastFrame = chunk;
        } else if (i === 0) {
            lastFrame = chunk;
        } else if (i > 0) {
            lastFrame = chunk.slice(i);
        } else {
            chunkQueue = concatArrays(chunkQueue, chunk);
        }
    });

    audioStream.on('end', async () => {
        console.info(`(${host}) (${sessionId}) Stream ended!`);
    });
}

function containsADTSHeader(data: Uint8Array) {
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