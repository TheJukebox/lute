import { writable } from 'svelte/store';
import { fetchStream, seek } from '$lib/stream_handler';

export type Track = {
    path: string;
    title: string;
    artist: string;
    album: string;
    duration: number;
}

export const library: Array<Track> = [
        { title: 'Something In The Way', artist: 'Nirvana', album: 'Nevermind', path: "uploads/converted/SomethingInTheWay.aac", duration: 235 },
        { title: 'Breezeblocks', artist: 'alt-j', album: 'alt-j', path: "uploads/converted/Breezeblocks.aac", duration: 226 },
        { title: 'R U Mine', artist: 'Arctic Monkeys', album: 'AM', path: "uploads/converted/R_U_Mine.aac", duration: 200 },
        { title: 'Lua', artist: 'Bright Eyes', album: 'I\'m Wide Awake, It\'s Morning', path: "uploads/converted/Lua.aac", duration: 262 },
];

// Writable stores
export const isPlaying = writable<boolean>(false);
export const isSeeking = writable<boolean>(false);
export const currentTime = writable<number>(0);
export const bufferedTime = writable<number>(0);
export const buffering = writable<boolean>(false);
export const currentTrack = writable<Track> ({
    path: '',
    title: '--',
    artist: '--',
    album: '--',
    duration: 0,
});

export function startStream(track: Track): void {
    buffering.set(true);
    currentTrack.set(track);
    fetchStream('http://127.0.0.1:8080', track, 'test-session');
    isPlaying.set(true);
}

export function seekToTime(time: number = 0) {
    seek(time);
}
