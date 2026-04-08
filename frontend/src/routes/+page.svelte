<script lang="ts">
    import { onMount } from 'svelte';
    import Upload from '$lib/components/Upload.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import Nav from '$lib/components/Nav.svelte';
    import ArtistBox from '$lib/components/ArtistBox.svelte';
    import {artistList, fetchArtists} from '$lib/tracklist.svelte';

    const playlistList = $state({playlists: []});

    async function fetchPlaylists() {
        const response: Response = await fetch(
            "http://localhost:7001/playlists",
            {
                method: "GET",
            },
        );
        if (!response.ok) {
            console.error("failed to fetch playlists");
            return;
        }
        playlistList.playlists = await response.json();
        console.log($state.snapshot(playlistList.playlists));
    }

    $effect(() => {
        fetchArtists();
        fetchPlaylists();
    });
</script>

<Nav />
<main class="flex flex-col gap-10 items-center bg-gray-950 min-h-screen min-w-screen">
    <div class="w-full">
        <Upload/>
    </div>
    <div class="w-1/2">
        <h2 class="text-2xl text-white">Artists</h2>
    </div>
    <div>
        {#if artistList.loading }
            <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition">
            {#each { length: 8 } }
                <div>
                    <div class="animate-pulse flex md:w-65 md:h-65 w-40 h-40 bg-slate-200 rounded-lg shadow-sm aspect-square items-center justify-center">
                    </div>
                </div>
            {/each}
            </div>
        {:else }
            <div class="grid grid-cols-1 gap-5 sm:grid-cols-2 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 transition">
                {#each artistList.artists as artist}
                    <div>
                        <ArtistBox artist={artist} />
                    </div>
                {/each}
            </div>
        {/if}
    </div>
    <div class="w-1/2">
        <h2 class="text-2xl text-white">Playlists</h2>
    </div>
    <div class="flex flex-col gap-5 p-5 items-start w-1/2">
        {#each playlistList.playlists as playlist}
            <a class="flex gap-5 items-center text-xl text-lime-200 hover:text-lime-400 hover:cursor-pointer hover:scale-[1.05] transition" href="/playlists/{playlist.id}">
                <div class="w-20 h-20 bg-gray-700 border-1 border-gray-600 rounded-sm"></div>
                <div>{playlist.title}</div>
            </a>
        {/each}
    </div>
</main>
