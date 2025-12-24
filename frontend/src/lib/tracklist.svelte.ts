import type { Track } from '$lib/upload';

export const trackList = $state({ tracks: [], loading: true, nowPlaying: "", currentTrack: {}, currentArtist: "", currentAlbum: ""});

export async function fetchTracks() {
    trackList.loading = true;
    trackList.currentArtist = "";
    const response = await fetch("http://localhost:7001/tracks");
    const data = await response.json();
    trackList.tracks = data.tracks;
    trackList.loading = false;
}

export async function fetchArtistTracks(artistID: string) {
    trackList.loading = true;
    const response = await fetch(`http://localhost:7001/tracks?artist=${artistID}`);
    const data = await response.json();
    trackList.tracks = data.tracks;
    trackList.loading = false;
}

export async function fetchAlbumTracks(albumID: string) {
    trackList.loading = true;
    const response = await fetch(`http://localhost:7001/tracks?album=${albumID}`);
    const data = await response.json();
    trackList.tracks = data.tracks;
    trackList.loading = false;
}
