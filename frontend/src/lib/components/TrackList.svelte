<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext, audioContext } from '$lib/stream';
    import { type StreamChunk, StreamBuffer } from '$lib/stream';
    
    const buffer: StreamBuffer = new StreamBuffer();
    let tracks = $state([]);
    let loading = $state(true);

    const fetchTracks = async () => {
        const response = await fetch("http://localhost:7001/tracks");
        const data = await response.json();
        tracks = data.tracks;
        loading = false;
    };

    async function playTrack(track: string) {
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

    onMount(() => {
        fetchTracks();
    });
</script>

<div class="p-4 shadow rounded-xl bg-lime-50 min-h-80">
    {#if loading}
        <div>
            Loading...
        </div>
    {:else}
        {#each tracks as track}
    <div 
        onclick={() => playTrack(track.Path)}
        class="p-2 bg-lime-100 border border-b border-lime-200 hover:cursor-pointer"
    >
                {track.Name}
            </div>
        {/each}
    {/if}
</div>
