import type { Track } from '$lib/upload';
import { StreamBuffer } from '$lib/stream';
import type { StreamChunk } from '$lib/stream';
import {audioContext, getAudioContext} from './stream';

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
    offset: 0,
    timeElapsed: 0,
    volume: 1,
    seeking: false,
    countInterval: 0,
    muted: false,
});

function countElapsed() {
    playback.timeElapsed = (audioContext.currentTime - playback.startedAt + playback.offset);
    if (playback.timeElapsed == playback.duration) {
        clearInterval(playback.countInterval);
    }
}

export function setGain(value: number) {
    if (playback.gain) {
        playback.gain.gain.setValueAtTime(value, audioContext.currentTime);
        playback.volume = value;
        localStorage.setItem("volume", value.toString());
    }
}

export async function fadeOut(duration: number = 1) {
    if (playback.node && playback.gain) {
        playback.gain.gain.linearRampToValueAtTime(0, audioContext.currentTime + duration);
        playback.node.stop(audioContext.currentTime + duration);
    }
}

export async function fadeIn(duration: number = 1) {
    if (playback.node && playback.gain) {
        playback.gain.gain.linearRampToValueAtTime(playback.volume, audioContext.currentTime + duration);
    }
}

export async function startPlaybackAt(offset: number = 0) {
    playback.offset = offset;
    playback.playing = false;
    if (playback.countInterval) {
        clearInterval(playback.countInterval);
    }
    if (playback.node) { 
        playback.node.onended = null;
        playback.node.stop();
        playback.node = undefined; 
    };
    getAudioContext()
    if (playback.track && playback.audio && audioContext) {
        // gain
        playback.gain = audioContext.createGain();
        playback.gain.gain.setValueAtTime(playback.volume, audioContext.currentTime);

        // buffer
        playback.node = audioContext?.createBufferSource();
        playback.node.connect(playback.gain);
        playback.node.buffer = playback.audio;
        playback.duration = playback.audio.duration; 
        playback.timeElapsed = offset;

        // playback
        playback.gain.connect(audioContext.destination);
        playback.node.start(audioContext.currentTime, offset);
        audioContext.resume();
        playback.playing = true;
        playback.startedAt = audioContext.currentTime;
        playback.countInterval = setInterval(countElapsed);

        playback.node.onended = () => {
            clearInterval(playback.countInterval);
        };
    }
}

export async function restartPlayback() {
    playback.playing = false;
    if (playback.countInterval) {
        clearInterval(playback.countInterval);
    }
    if (playback.node) { 
        playback.node.stop();
        playback.node = undefined; 
    };
    getAudioContext()
    if (playback.track && playback.audio && audioContext) {
        // gain
        playback.gain = audioContext.createGain();
        playback.gain.gain.setValueAtTime(playback.volume, audioContext.currentTime);

        // buffer
        playback.node = audioContext?.createBufferSource();
        playback.node.connect(playback.gain);
        playback.node.buffer = playback.audio;
        playback.duration = playback.audio.duration; 
        playback.offset = 0;
        playback.timeElapsed = 0;

        // playback
        playback.gain.connect(audioContext.destination);
        playback.node.start();
        audioContext.resume();
        playback.playing = true;
        playback.startedAt = audioContext.currentTime;
        playback.countInterval = setInterval(countElapsed);

        playback.node.onended = () => {
            clearInterval(playback.countInterval);
        };
    }
}

export function togglePlayback() {
    getAudioContext()
    if (audioContext) {
        audioContext.state === "suspended" ? audioContext.resume() : audioContext.suspend();
        audioContext.state === "suspended" ? playback.playing = true : playback.playing = false;
    }
}

export async function startPlayback() {
    playback.playing = false;
    if (playback.countInterval) {
        clearInterval(playback.countInterval);
    }
    if (playback.node) { 
        playback.node.onended = null;
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
        playback.audio = await audioContext?.decodeAudioData(playable.buffer as ArrayBuffer);

        // gain
        playback.gain = audioContext.createGain();
        playback.gain.gain.setValueAtTime(playback.volume, audioContext.currentTime);

        // buffer
        playback.node = audioContext?.createBufferSource();
        playback.node.connect(playback.gain);
        playback.node.buffer = playback.audio;
        playback.duration = playback.audio.duration; 
        playback.offset = 0;
        playback.timeElapsed = 0;

        // playback
        playback.gain.connect(audioContext.destination);
        playback.node.start();
        audioContext.resume();
        playback.playing = true;
        playback.startedAt = audioContext.currentTime;
        playback.countInterval = setInterval(countElapsed);

        playback.node.onended = () => {
            clearInterval(playback.countInterval);
        };
    } else {
        console.error("Could not fetch audio context for playback!");
    }
}
