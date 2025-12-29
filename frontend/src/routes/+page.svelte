<script lang="ts">
    import { onMount } from 'svelte';
    import Upload from '$lib/components/Upload.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import Controls from '$lib/components/Controls.svelte';
    import ArtistBox from '$lib/components/ArtistBox.svelte';
    import {artistList, fetchArtists} from '$lib/tracklist.svelte';


    onMount(() => {
        fetchArtists();
    });
</script>
<main class="flex flex-col items-center bg-green-50 min-h-screen min-w-screen">
    <div class="p-4 flex flex-col">
        <div class="flex justify-center">
            <Controls />
        </div>
        <div>
            <h1 class="ml-20 p-2 text-xl font-bold">Home</h1>
            <h2 class="ml-20 p-2 text-lg font-semibold">Artists</h2>
        </div>
        {#if artistList.loading }
            <div class="p-5 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition">
            {#each { length: 8 } }
                <div class="p-2">
                    <div class="animate-pulse flex md:w-65 md:h-65 w-40 h-40 bg-slate-200 rounded-lg shadow-sm aspect-square items-center justify-center">
                    </div>
                </div>
            {/each}
            </div>
        {:else }
            <div class="p-5 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition">
                {#each artistList.artists as artist}
                    <div class="p-2">
                        <ArtistBox artist={artist} />
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</main>
