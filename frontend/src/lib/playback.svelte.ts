import type { Track } from '$lib/upload';
import { StreamBuffer } from '$lib/stream';
import type { StreamChunk } from '$lib/stream';
import {audioContext, getAudioContext} from './stream';

let counting: number = 0;

export const playback = $state({ 
    playing: false,
    track: {},
    trackIndex: 0,
    buffer: new StreamBuffer(),
    audio: {},
    node: undefined,
    gain: undefined,
    duration: 0,
    startedAt: 0,
    timeElapsed: 0,
    volume: 1,
});

function countElapsed() {
    playback.timeElapsed = (audioContext.currentTime - playback.startedAt);
    if (playback.timeElapsed == playback.duration) {
        clearInterval(counting);
    }
}

export async function fadeOut(duration: number = 1) {
    if (playback.node) {
        playback.node.stop(duration);
    }
}

export async function restartPlayback() {
    playback.playing = false;
    if (counting) {
        clearInterval(counting);
    }
    if (playback.node) { 
        playback.node.stop();
        playback.node = undefined; 
    };
    getAudioContext()
    if (playback.track && playback.audio && audioContext) {
        playback.timeElapsed = 0;
        playback.startedAt = Date.now();
        playback.node = audioContext?.createBufferSource();
        playback.node.buffer = playback.audio;
        playback.node.connect(audioContext?.destination);
        playback.duration = playback.audio.duration; 
        playback.node.start();
        audioContext.resume();
        playback.playing = true;
        playback.startedAt = audioContext.currentTime;
        counting = setInterval(countElapsed);
    }
}

export function togglePlayback() {
    getAudioContext()
    if (audioContext) {
        audioContext.state === "suspended" ? audioContext.resume() : audioContext.suspend();
        audioContext.state === "suspended" ? playback.playing = true : playback.playing = false;
    }
}

export async function startPlayback(track: Track) {
    playback.playing = false;
    if (counting) {
        clearInterval(counting);
    }
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
        playback.audio = await audioContext?.decodeAudioData(playable.buffer as ArrayBuffer);

        // gain
        playback.gain = audioContext.createGain();
        playback.gain.gain.setValueAtTime(playback.volume, audioContext.currentTime);

        // buffer
        playback.node = audioContext?.createBufferSource();
        playback.node.connect(playback.gain);
        playback.node.buffer = playback.audio;
        playback.duration = playback.audio.duration; 

        // playback
        playback.gain.connect(audioContext.destination);
        playback.node.start();
        audioContext.resume();
        playback.playing = true;
        playback.startedAt = audioContext.currentTime;
        counting = setInterval(countElapsed);
    } else {
        console.error("Could not fetch audio context for playback!");
    }
}
