import { browser } from '$app/environment';

let audioContext: AudioContext;

export function getAudioContext() {
    if (!browser) return null;

    if (!audioContext) {
        audioContext = new AudioContext();
    }
    return audioContext;
}

