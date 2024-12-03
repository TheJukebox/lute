import stream from '$lib/gen/stream_grpc_web_pb';

let context: AudioContext | null = null;
const queue: AudioBuffer[] = [];

function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export function playFromBuffer() {
    if (!context) {
        // potentially need to set sampleRate option: https://webaudio.github.io/web-audio-api/#dom-baseaudiocontext-samplerate
        // by default it should use the target device's preferred sample rate
        context = new AudioContext(); 
    }
    const source: AudioBufferSourceNode = context.createBufferSource();
    const next: AudioBuffer | null | undefined = queue.shift();
    if (next === null || next === undefined) {
        return
    }
    source.buffer = next; 
    // @ts-ignore - context will not be null, even though ts thinks it could be
    source.connect(context.destination);
    source.start();
    playFromBuffer();
    sleep(3000000);
}

export async function decodeAudio(data: ArrayBuffer) : Promise<AudioBuffer> {
    // we defer the creation of the AudioContext to avoid "autoplay" policy.
    if (!context) {
        // potentially need to set sampleRate option: https://webaudio.github.io/web-audio-api/#dom-baseaudiocontext-samplerate
        // by default it should use the target device's preferred sample rate
        context = new AudioContext(); 
    }
    //return await context.decodeAudioData(data);

    return await context.decodeAudioData(data);
}

export function fetchStream(host: string, path: string, sessionId: string) {
    let service = new stream.AudioStreamClient(
        host,
        null,
        {
            'use-fetch': true,
        }
    );
    let request = new stream.AudioStreamRequest();
    request.setFileName(path);
    request.setSessionId(sessionId);

    console.info(`(${host}) (${sessionId}) Requesting stream '${path}'...`)
    const audioStream = service.streamAudio(request);
    console.info(`(${host}) (${sessionId}) Opened stream!`);

    audioStream.on('data', async (response: stream.AudioStreamChunk) => {
        const encoder = new TextEncoder()
        // @ts-ignore - for some reason it thinks it's going to be a string
        // after assignment...
        const data: Uint8Array = response.getData(); 

        const blob: Blob = new Blob([data]); // convert to a Blob first, so we can create an ArrayBuffer
        const buff: ArrayBuffer = await blob.arrayBuffer();
        const decodedData = await decodeAudio(buff);
        queue.push(decodedData);
        playFromBuffer();
        console.debug(`(${host}) (${sessionId}) Buffering chunk #${response.getSequence()}...`);
    });

    audioStream.on('end'), async (end: any) => {
        console.info(`(${host}) (${sessionId}) Stream ended.`);
    }
}