import '$lib/gen/stream_grpc_web_pb';
import { sequence } from '@sveltejs/kit/hooks';

let context: AudioContext | null = null;
const audioQueue: AudioBuffer[] = [];
let dataQueue: Uint8Array[] = [];

function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export function playFromBuffer() {
    if (context === null) {
        // potentially need to set sampleRate option: https://webaudio.github.io/web-audio-api/#dom-baseaudiocontext-samplerate
        // by default it should use the target device's preferred sample rate
        console.log("Created new context...");
        context = new AudioContext({"latencyHint": "playback", "sampleRate": 44100}); 
    }
    const source: AudioBufferSourceNode = context.createBufferSource();
    const next: AudioBuffer | null | undefined = audioQueue.shift();
    if (next === null || next === undefined) {
        return
    }
    
    source.buffer = next; 
    // @ts-ignore - context will not be null, even though ts thinks it could be
    source.connect(context.destination);
    source.start(0, 0, next.duration);
}

export async function decodeAudio(data: ArrayBuffer) : Promise<AudioBuffer> {
    // we defer the creation of the AudioContext to avoid "autoplay" policy.
    if (!context) {
        // potentially need to set sampleRate option: https://webaudio.github.io/web-audio-api/#dom-baseaudiocontext-samplerate
        // by default it should use the target device's preferred sample rate
        context = new AudioContext(); 
    }
    return await context.decodeAudioData(data);
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
    console.log(service);
    const audioStream = service.streamAudio(request, null);
    //console.info(`(${host}) (${sessionId}) Opened stream!`);
    audioStream.on('data', async (response: proto.stream.AudioStreamChunk) => {
        // @ts-ignore - for some reason it thinks it's going to be a string
        // after assignment...
        const data: Uint8Array = response.getData(); 
        const frameAt = containsADTSFrame(data)
        if (containsADTSFrame(data) !== -1) {
            console.log("FRAME AT: ", frameAt);
        }

    });

    audioStream.on('end', async () => {
    })
}

function containsADTSFrame(data: Uint8Array) {
    for (let i = 0; i < data.length - 7; i++) {
        if (data[i] === 0xFF && (data[i + 1] & 0xF0) === 0xF0) {
            // check if the ID = 0 (MPEG-4).
            if ((data[i + 1] & 0x08) >> 3 !== 0) {
                break
            }
            // layer
            if ((data[i + 1] & 0x06) >> 1 !== 0) {
                break
            }
            // protection_absent
            if ((data[i + 1] & 0x01) !== 0x01) {
                // i think most of the time it should be present, so...
                break
                // that said, make sure we revisit this
            }
            // sampling frequency
            if ((data[i + 2] & 0xF0) >> 4 !== 0x05) {
                break;
            }
            // private_bit
            if ((data[i + 2] & 0x08) >> 3 !== 0) {
                // sometimes, very rarely, this bit can be 1
                break;
                // revisit it
            }
            // channel configuration - i'd actually expect this to be 2,
            // but the file i'm working with is weird.
            if ((data[i + 2] & 0x07) !== 0) {
                break;
            }
            // it's a copy!
            if ((data[i + 3] & 0xC0) >> 6 === 1) {
                break;
            }
            if ((data[i + 3] & 0x20) >> 5 !== 0) {
                break;
            }
            return i;
        }
    }
    return -1;
}