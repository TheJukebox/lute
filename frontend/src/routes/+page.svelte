<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext } from '$lib/stream';
    import { type StreamChunk, StreamBuffer } from '$lib/stream';

    let ws: WebSocket;
    let status: string = 'Disconnected';

    onMount(() => {
        const audioContext = getAudioContext();
        let stream: Uint8Array[] = [];
        const buffer: StreamBuffer = new StreamBuffer();
        let go: boolean = true;
        ws = new WebSocket('ws://172.31.204.147:7001/stream');

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
            let playable: Uint8Array = new Uint8Array(0);
            while (true) {
                const chunk: StreamChunk | undefined = buffer.next();
                console.debug(`Adding chunk ${chunk?.Sequence} to buffer!`);
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

        const play = async () => {
            audioContext?.state === 'suspended' && audioContext?.resume();
        }
        document.addEventListener('click', play, {once: true});
    });

    const sendMessage = () => {
       ws.send("Test"); 
    };
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
<p>Websocket status: {status}</p>
<button
    onclick={sendMessage}
    aria-label="doit"
    class="p-2 text-xl border rounded-full"
>
    TEST
</button>
