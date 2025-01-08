import { writable } from 'svelte/store';
import { fetchStream, togglePlayback } from '$lib/stream_handler';

export interface Track {
    path: string;
    title: string;
    artist: string;
    album: string;
    paused: boolean;
    duration: number;
}

export const currentTime = writable<number>(0);

export const playing = writable<Track> ({
    path: '',
    title: '--',
    artist: '--',
    album: '--',
    paused: true,
    duration: 0,
});

export function startStream(path: string, title: string, artist: string, album: string, duration: number): void {
    playing.set({
        path,
        title,
        artist,
        album,
        paused: false,
        duration,
    });
    fetchStream('http://127.0.0.1:8080', path, 'test-session');
    togglePlayback();
}

export function toggleStream(): void {
    togglePlayback();
}