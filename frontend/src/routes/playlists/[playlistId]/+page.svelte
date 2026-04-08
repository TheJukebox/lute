<script lang="ts">
    import { onMount } from 'svelte';
    import type { PageProps } from './$types';
    import Nav from '$lib/components/Nav.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import PlaylistModal from '$lib/components/PlaylistModal.svelte';
    import { fetchPlaylistTracks, fetchPlaylist, trackList } from '$lib/tracklist.svelte';

    let { params }: PageProps = $props();
    let playlist = $state(null);

    $effect(() => {
        const playlistId = params.playlistId;
        console.log(`Fetching tracks for playlist '${playlistId}'.`);
        fetchPlaylistTracks(playlistId);
        (async () => {
            playlist = await fetchPlaylist(playlistId);
        })();
    });
</script>

<Nav/>
<main class="flex flex-col gap-4 p-3 items-center min-w-screen min-h-screen bg-gray-950">
    <div class="lg:w-1/2 md:w-3/4">
        <div class="flex gap-10 p-3">
            <div class="rounded-sm bg-gray-800 border-1 border-gray-700 w-50 h-50"></div>
            <div>
                <h1 class="text-white text-2xl">{playlist?.title ?? "..."}</h1>
                <h2 class="text-gray-400 text-xl">{playlist?.description ?? "..."}</h2>
            </div> 
        </div>
        <div class="flex justify-center">
            <TrackList />
        </div>
    </div>
    <PlaylistModal/>
</main>
