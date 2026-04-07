<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext, audioContext } from '$lib/stream';
    import { type Track } from '$lib/upload';
    import { type StreamChunk, StreamBuffer } from '$lib/stream';
    import { trackList, fetchTracks, fetchArtistTracks, fetchAlbumTracks } from '$lib/tracklist.svelte';
    import { startPlayback, playback, stopPlayback } from '$lib/playback.svelte'; 
    import { LoaderCircle, CirclePlus } from '@lucide/svelte';

    async function playTrack(track: Track) {
        trackList.nowPlaying = track.id;
        trackList.currentTrack = track; 
        const ws: WebSocket = new WebSocket(
            `ws://localhost:7001/stream?track=${encodeURIComponent(track.path)}`
        );
        ws.onopen = () => {
            console.debug("Websocket open.")
        };

        ws.onmessage = async (e) => {
            const data: StreamChunk = JSON.parse(await e.data.text())
            playback.buffer.push(data);
        };

        ws.onclose = async (e) => {
            console.debug("Websocket closed.")
            startPlayback(track);
        };

        ws.onerror = (e) => {
            console.error("Websocket error: ", e)
        };
    };

    function isPlaying(id: string): boolean {
        if (trackList.nowPlaying === id) { return true; } 
        return false;
    };

    function trackStyle(id: string): string {
        if (isPlaying(id)) {
            return "hover:cursor-pointer text-lime-200 hover:text-lime-400 border hover:text-gray-700 transition";
        }
        return "hover:cursor-pointer text-white hover:text-lime-400 transition";
    };

    onMount(() => {
    });

</script>

<div class="p-2 rounded-sm min-w-full min-h-80 flex flex-col overflow-x-clip h-full max-h-320 transition scroll-smooth">
        <div class="grid grid-cols-[0.1fr_1fr_1fr_1fr_0.1fr] gap-4 px-4 py-2 text-left min-w-full text-white bg-gray-800">
            <div>No.</div>
            <div>Track</div>
            <div>Artist</div>
            <div>Album</div>
            <div></div>
        </div>
            {#each trackList.tracks as track}
            <button 
                class={`grid grid-cols-[0.1fr_1fr_1fr_1fr_0.1fr] gap-4 px-4 py-2 text-left min-w-full text-gray-600 bg-gray-700 ${trackStyle(track.id)}`}
                onclick={() => {
                    stopPlayback();
                    playTrack(track);
                }}
                type="button"
            >
                <div>{track.trackNumber}</div>
                <div>{track.title}</div>
                <div>
                    {track.artist.Name}
                </div>
                <div>
                    {track.album.Title}
                </div>
                <div class="flex items-center justify-center self-center bg-gray-600 hover:bg-lime-400 aspect-square w-5 h-5 rounded-full text-xl hover:text-lime-200 hover:scale-[1.05]" onclick={(e) => { e.preventDefault(); e.stopPropagation(); console.log("add to playlist");}}><CirclePlus/></div>
            </button>
        {/each}
</div>
