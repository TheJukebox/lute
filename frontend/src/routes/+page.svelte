<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext, audioContext } from '$lib/stream';
    import { type StreamChunk, StreamBuffer } from '$lib/stream';
    import { Play } from '@lucide/svelte';
    import Upload from '$lib/components/Upload.svelte';

    let ws: WebSocket;
    let status: string = 'Disconnected';
    const buffer: StreamBuffer = new StreamBuffer();

    onMount(() => {
        ws = new WebSocket('ws://localhost:7001/stream?track=NECROPOLITE.mp3');

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
        };

        ws.onerror = (e) => {
            status = `Error: ${e}`;
        };

        const play = async () => {
            audioContext?.state === 'suspended' && audioContext?.resume();
        }
        document.getElementById("play")?.addEventListener('click', play, {once: true});
    });

    const sendMessage = () => {
       ws.send("Test"); 
    };

    const fetchTracks = async () => {
        const response = await fetch("http://localhost:7001/tracks");
        console.log(response.data);
    };

    const play = async () => {
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
    }
</script>
<main>
    <div class="p-4">
        <h1 class="text-xl font-bold">Welcome to SvelteKit</h1>
        <p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
        <p>Websocket status: {status}</p>
        <div class="p-5 flex flex-col max-w-48 gap-4 items-center">
            <button
                onclick={sendMessage}
                aria-label="send-message"
                class="p-2 text-xl border rounded-full hover:bg-blue-100 active:scale-[0.95] transition shadow cursor-pointer"
            >
                Send message to websocket
            </button>
            <button
                onclick={fetchTracks}
                aria-label="fetch-tracks"
                class="p-2 text-xl border rounded-full hover:bg-blue-100 active:scale-[0.95] transition shadow cursor-pointer"
            >
               Fetch Tracks 
            </button>


            <button
                aria-label="play"
                onclick={play}
                class="p-2 max-w-10 max-h-10 text-xl rounded-full bg-lime-400 hover:bg-lime-600 active:scale-[0.95] transition shadow cursor-pointer items-center"
            >
                <Play class="h-5 w-5" />
            </button>
        </div>
        <Upload />
    </div>
</main>
