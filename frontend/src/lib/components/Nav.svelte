<script lang="ts">
    import Controls from '$lib/components/Controls.svelte';
    import { trackList } from '$lib/tracklist.svelte';
</script>

<div class="grid grid-cols-3 items-center w-full h-20 bg-gray-900 gap-2">
    <div class="flex gap-2 p-4">
        <a href="/" class="flex items-center justify-center self-center border-1 border-gray-600 hover:border-lime-400 rounded-md bg-gray-800 text-white hover:text-gray-900 w-28 h-10 hover:bg-lime-300 transition hover:cursor-pointer">Home</a>
        <a href="/" class="flex items-center justify-center self-center border-1 border-gray-600 hover:border-lime-400 rounded-md bg-gray-800 text-white hover:text-gray-900 w-28 h-10 hover:bg-lime-300 transition hover:cursor-pointer">Library</a>
        <a href="/" class="flex items-center justify-center self-center border-1 border-gray-600 hover:border-lime-400 rounded-md bg-gray-800 text-white hover:text-gray-900 w-28 h-10 hover:bg-lime-300 transition hover:cursor-pointer">Queue</a>
    </div>
    <div class="flex">
        <Controls class="justify-self-center grid-span-2" />
    </div>
</div>
