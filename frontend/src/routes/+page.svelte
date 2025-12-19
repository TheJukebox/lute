<script lang="ts">
    import { onMount } from 'svelte';
    import { Play } from '@lucide/svelte';
    import Upload from '$lib/components/Upload.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
</script>
<main>
    <div class="p-4 flex flex-col gap-4">
        <h1 class="text-xl font-bold">Welcome to SvelteKit</h1>
        <p>Visit <a href="https://svelte.dev/docs/kit">svelte.dev/docs/kit</a> to read the documentation</p>
        <Upload />
        <TrackList />
    </div>
</main>
