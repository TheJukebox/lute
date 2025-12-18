import { browser } from '$app/environment';

let audioContext: AudioContext;

export interface StreamChunk {
    ChunkSize: number;
    Sequence: number;
    Data: string;
}

export function getAudioContext() {
    if (!browser) return null;

    if (!audioContext) {
        audioContext = new AudioContext();
    }
    return audioContext;
}

export class StreamBuffer {
    head: number;
    chunks: Array<StreamChunk>;
    
    constructor() {
        this.head = 0
        this.chunks = new Array(0);
    }

    push(data: StreamChunk) {
        this.chunks.push(data);
        this.chunks.sort((a, b) => a.Sequence - b.Sequence);
    }

    next() {
        return this.chunks.shift()
    }
}
