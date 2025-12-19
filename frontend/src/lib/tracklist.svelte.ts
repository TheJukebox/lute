import type { Track } from '$lib/upload';

export const trackList = $state({ tracks: [], loading: true, nowPlaying: "", currentTrack: {}});

export async function fetchTracks() {
    trackList.loading = true;
    const response = await fetch("http://localhost:7001/tracks");
    const data = await response.json();
    trackList.tracks = data.tracks;
    trackList.loading = false;
}
