<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext } from '$lib/stream';

    let ws: WebSocket;
    let status: string = 'Disconnected';

    onMount(() => {
        let audioContext = getAudioContext()
        let stream: Uint8Array[] = []
        ws = new WebSocket('ws://172.31.204.147:7001/stream');

        ws.onopen = () => {
            status = 'Connected';
        };

        ws.onmessage = async (e) => {
            stream.push(new Uint8Array(await e.data.arrayBuffer())); 
        };

        ws.onclose = async (e) => {
            status = `Disconnected: ${e}`;
            console.log(stream);
            let out: Uint8Array = new Uint8Array(0);
            for (const buf of stream) {
                let tmp = new Uint8Array(out.byteLength + buf.byteLength);
                tmp.set(out, 0);
                tmp.set(buf, out.byteLength);
                out = tmp; 
            }
            console.log(out);
            let audio: AudioBuffer = audioContext?.decodeAudioData(out.buffer as ArrayBuffer)
            console.log(audio);
            const source = audioContext?.createBufferSource();
            source.buffer = await audio;
            source.connect(audioContext?.destination);
            source.start();
        };

        ws.onerror = (e) => {
            status = `Error: ${e}`;
        };

        const play = () => audioContext?.state === 'suspended' && audioContext?.resume();
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
