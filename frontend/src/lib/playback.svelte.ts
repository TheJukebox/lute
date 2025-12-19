import type { Track } from '$lib/upload';
import { StreamBuffer } from '$lib/stream';
import type { StreamChunk } from '$lib/stream';
import {audioContext, getAudioContext} from './stream';

export const playback = $state({ playing: false, track: {}, buffer: new StreamBuffer(), node: undefined});

export async function togglePlayback() {
    getAudioContext()
    if (audioContext) {
        audioContext.state === "suspended" ? audioContext.resume() : audioContext.suspend();
        audioContext.state === "suspended" ? playback.playing = true : playback.playing = false;
    }
}

export async function startPlayback(track: Track) {
    playback.playing = false;
    if (playback.node) { 
        playback.node.stop();
        playback.node = undefined; 
    };
    getAudioContext()
    if (audioContext) {
        let playable: Uint8Array = new Uint8Array(0);
        while (true) {
            const chunk: StreamChunk | undefined = playback.buffer.next();
            if (chunk === undefined) {
                break;
            }
            const data: Uint8Array = Uint8Array.from(atob(chunk.Data), c => c.charCodeAt(0));
            let temp = new Uint8Array(playable.byteLength + data.byteLength);
            temp.set(playable, 0);
            temp.set(data, playable.byteLength);
            playable = temp;
        }
        console.log(playable);
        const audio: AudioBuffer = await audioContext?.decodeAudioData(playable.buffer as ArrayBuffer);
        playback.node = audioContext?.createBufferSource();
        playback.node.buffer = audio;
        playback.node.connect(audioContext?.destination);
        playback.node.start();
        audioContext.resume();
        playback.playing = true;
    } else {
        console.error("Could not fetch audio context for playback!");
    }
}
