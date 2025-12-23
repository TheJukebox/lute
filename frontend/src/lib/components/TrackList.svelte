<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext, audioContext } from '$lib/stream';
    import { type Track } from '$lib/upload';
    import { type StreamChunk, StreamBuffer } from '$lib/stream';
    import { trackList, fetchTracks } from '$lib/tracklist.svelte';
    import { startPlayback, playback } from '$lib/playback.svelte'; 


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
            return "hover:bg-lime-200 hover:cursor-pointer text-lime-900 font-bold border border-sm border-lime-500 rounded-lg hover:text-green-700 shadow-sm";
        }
        return "hover:bg-lime-200 hover:cursor-pointer hover:scale-[1.001] text-lime-900 hover:font-bold rounded-lg hover:shadow";
    };

    onMount(() => {
        fetchTracks();
    });
</script>

<div class="p-4 shadow rounded-xl bg-lime-50 min-h-80 flex flex-col gap-2 overflow-x-clip h-full max-h-320 transition scroll-smooth">
        <div class="grid grid-cols-[0.1fr_1fr_1fr_1fr] gap-4 px-4 py-2 text-left min-w-full text-lime-600 font-semibold rounded-full bg-lime-200 shadow">
            <div>No.</div>
            <div>Track</div>
            <div>Artist</div>
            <div>Album</div>
        </div>
        {#if trackList.loading}
            <div class="flex flex-1 justify-center items-center max-w-full max-h-full text-xl">
                Loading...
            </div>
        {:else}
        <div class="shadow rounded-xl bg-lime-100 overflow-y-auto overflow-x-clip">
            {#each trackList.tracks as track}
                <button 
                    class={`grid grid-cols-[0.1fr_1fr_1fr_1fr] gap-4 px-4 py-2 text-left min-w-full transition ${trackStyle(track.ID)}`}
                    onclick={() => playTrack(track)}
                    type="button"
                >
                    <div>{track.trackNumber}</div>
                    <div>{track.title}</div>
                    <div>{track.artist.Name}</div>
                    <div>{track.album.Title}</div>
                </button>
            {/each}
        </div>
    {/if}
</div>
