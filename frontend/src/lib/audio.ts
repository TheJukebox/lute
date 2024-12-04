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
        dataQueue = dataQueue.concat(data);
    });

    audioStream.on('end', async () => {
        console.log(dataQueue);
        let blob = new Blob(dataQueue);
        console.log(blob);
        let buff = await blob.arrayBuffer();
        console.log(buff);
        let audio = await decodeAudio(buff);
        console.log(audio);
        audioQueue.push(audio);
        playFromBuffer();
    })
}
