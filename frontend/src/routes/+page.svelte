<script lang="ts">
    import { onMount } from 'svelte';
    import Upload from '$lib/components/Upload.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import Nav from '$lib/components/Nav.svelte';
    import ArtistBox from '$lib/components/ArtistBox.svelte';
    import {artistList, fetchArtists} from '$lib/tracklist.svelte';


    onMount(() => {
        fetchArtists();
    });
</script>

<Nav />
<main class="flex flex-col items-center bg-gray-950 min-h-screen min-w-screen">
    <div class="w-full">
        <Upload/>
    </div>
    <div class="p-4">
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
