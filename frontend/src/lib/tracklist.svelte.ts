import type { Track } from '$lib/upload';

export const trackList = $state({ tracks: [], noTracks: true, loading: true, nowPlaying: "", currentTrack: {}, currentArtist: {}, currentAlbum: {} });
export const artistList = $state({ artists: [], loading: true });

export async function fetchPlaylistTracks(playlistId: string) {
    trackList.loading = true;
    trackList.tracks = [];
    const response = await fetch(`http://localhost:7001/playlists/${playlistId}`);
    const data = await response.json();
    console.log(data);
    trackList.tracks = data; 
    if (trackList.tracks.length === 0) {
        console.log("no tracks");
    }
    trackList.loading = false;
}

export async function fetchTracks() {
    trackList.loading = true;
    trackList.tracks = [];
    const response = await fetch("http://localhost:7001/tracks");
    const data = await response.json();
    trackList.tracks = data.tracks.sort((a, b) => {
        const albumSort = a.album.Title.localeCompare(b.album.Title);
        if (albumSort !== 0) return albumSort;
        return a.trackNumber - b.trackNumber 
    });
    trackList.loading = false;
}

export async function fetchArtist(artistID: string) {
    const response = await fetch(`http://localhost:7001/artists/${artistID}`);
    const data = await response.json();
    return data
}

export async function fetchArtists() {
    artistList.loading = true;
    const response = await fetch("http://localhost:7001/artists");
    const data = await response.json();
    artistList.artists = data.artists;
    artistList.loading = false;
}

export async function fetchArtistTracks(artistID: string) {
    trackList.loading = true;
    trackList.tracks = [];
    const response = await fetch(`http://localhost:7001/artists/${artistID}/tracks`);
    const data = await response.json();
    trackList.tracks = data.tracks.sort((a, b) => {
        const albumSort = a.album.Title.localeCompare(b.album.Title);
        if (albumSort !== 0) return albumSort;
        return a.trackNumber - b.trackNumber 
    });
    trackList.loading = false;
}

export async function fetchAlbumTracks(albumID: string) {
    trackList.loading = true;
    trackList.tracks = [];
    const response = await fetch(`http://localhost:7001/albums/${albumID}/tracks`);
    const data = await response.json();
    trackList.tracks = data.tracks.sort((a, b) => {
        const albumSort = a.album.Title.localeCompare(b.album.Title);
        if (albumSort !== 0) return albumSort;
        return a.trackNumber - b.trackNumber 
    });
    trackList.loading = false;
}
