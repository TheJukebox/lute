import { writable } from 'svelte/store';
import { fetchStream, togglePlayback, seek } from '$lib/stream_handler';

export type Track = {
    path: string;
    title: string;
    artist: string;
    album: string;
    duration: number;
}

// Writable stores
export const isPlaying = writable<boolean>(false);
export const isSeeking = writable<boolean>(false);
export const currentTime = writable<number>(0);
export const currentTrack = writable<Track> ({
    path: '',
    title: '--',
    artist: '--',
    album: '--',
    duration: 0,
});

export function startStream(track: Track): void {
    currentTrack.set(track);
    fetchStream('http://127.0.0.1:8080', track, 'test-session');
    isPlaying.set(true);
}

export function seekToTime(time: number = 0) {
    seek(time);
}
