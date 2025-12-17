<script lang="ts">
    import { onMount } from 'svelte';
    import { getAudioContext } from '$lib/stream';

    let ws: WebSocket;
    let status: string = 'Disconnected';

    onMount(() => {
        audioContext = getAudioContext()
        ws = new WebSocket('ws://localhost:7001/ws');

        ws.onopen = () => {
            status = 'Connected';
        };

        ws.onmessage = (e) => {
            const message: string = e.data;
            status = `Message received: ${message}`;
        };

        ws.onclose = (e) => {
            status = `Disconnected: ${e}`;
        };

        ws.onerror = (e) => {
            status = `Error: ${e}`;
        };

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
