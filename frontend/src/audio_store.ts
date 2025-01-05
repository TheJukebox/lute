import { writable } from 'svelte/store';
import { fetchStream, togglePlayback } from '$lib/audio_processing';

export interface Track {
    path: string;
    title: string;
    artist: string;
    album: string;
    paused: boolean;
}

export const playing = writable<Track> ({
    path: '',
    title: '--',
    artist: '--',
    album: '--',
    paused: true,
});

export function startStream(path: string, title: string, artist: string, album: string): void {
    playing.set({
        path,
        title,
        artist,
        album,
        paused: false
    });
    fetchStream('http://127.0.0.1:8080', path, 'test-session');
    togglePlayback();
}

export function toggleStream(): void {
    togglePlayback();
}