<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext, audioContext } from '$lib/stream';
    import { type StreamChunk, StreamBuffer } from '$lib/stream';
    import { trackList, fetchTracks } from '$lib/tracklist.svelte.ts';
    
    const buffer: StreamBuffer = new StreamBuffer();

    async function playTrack(track: string, id: string) {
        trackList.nowPlaying = id;
        const ws: WebSocket = new WebSocket(
            `ws://localhost:7001/stream?track=${track}`
        );
        ws.onopen = () => {
            status = 'Connected';
        };

        ws.onmessage = async (e) => {
            const data: StreamChunk = JSON.parse(await e.data.text())
            buffer.push(data);
        };

        ws.onclose = async (e) => {
            status = `Disconnected: ${e}`;
            console.log(buffer);

            getAudioContext();
            audioContext?.state === 'suspended' && audioContext?.resume();
            let playable: Uint8Array = new Uint8Array(0);
            while (true) {
                const chunk: StreamChunk | undefined = buffer.next();
                if (chunk === undefined) {
                    break;
                }
                const data: Uint8Array = Uint8Array.from(atob(chunk.Data), c => c.charCodeAt(0));
                let temp = new Uint8Array(playable.byteLength + data.byteLength);
                temp.set(playable, 0);
                temp.set(data, playable.byteLength);
                playable = temp;
            }
            console.log(playable);
            const audio: AudioBuffer = await audioContext?.decodeAudioData(playable.buffer as ArrayBuffer);
            const source = audioContext?.createBufferSource();
            source.buffer = audio;
            source.connect(audioContext?.destination);
            source.start();
        };

        ws.onerror = (e) => {
            status = `Error: ${e}`;
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
                    onclick={() => playTrack(encodeURIComponent(track.Path), track.ID)}
                    type="button"
                >
                    <div>{track.Number}</div>
                    <div>{track.Name}</div>
                    <div>{track.Artist}</div>
                    <div>{track.Album}</div>
                </button>
            {/each}
        </div>
    {/if}
</div>
