<script lang="ts">
    import { onMount } from 'svelte';
    import Upload from '$lib/components/Upload.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import Controls from '$lib/components/Controls.svelte';
</script>
<main class="flex flex-col items-center bg-green-50 min-h-screen min-w-screen">
    <div class="p-4 flex flex-col gap-4 w-full h-full">
        <div class="flex justify-center">
            <Controls />
        </div>
        <Upload />
        <TrackList />
    </div>
</main>
