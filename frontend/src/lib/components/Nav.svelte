<script lang="ts">
    import Controls from '$lib/components/Controls.svelte';
    import { trackList } from '$lib/tracklist.svelte';
</script>

<div class="grid grid-cols-3 items-center w-full h-20 bg-gray-900 gap-2">
    <div class="flex gap-2 p-4">
        <a href="/" class="self-center rounded-full bg-gray-800 text-white w-30 h-10 hover:bg-sky-700 transition hover:cursor-pointer">Home</a>
        <a href="/" class="self-center rounded-full bg-gray-800 text-white w-30 h-10 hover:bg-sky-700 transition hover:cursor-pointer">Library</a>
        <a href="/" class="self-center rounded-full bg-gray-800 text-white w-30 h-10 hover:bg-sky-700 transition hover:cursor-pointer">Queue</a>
    </div>
    <div class="flex">
        <Controls class="justify-self-center grid-span-2" />
    </div>
</div>
