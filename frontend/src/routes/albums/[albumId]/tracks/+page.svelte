<script lang="ts">
    import { onMount } from 'svelte';
    import Controls from '$lib/components/Controls.svelte';
    import Nav from '$lib/components/Nav.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import PlaylistModal from '$lib/components/PlaylistModal.svelte';
    import { fetchAlbumTracks, trackList, fetchArtist } from '$lib/tracklist.svelte';

    let { params }: PageProps = $props();
    
    $effect(() => {
        const albumId = params.albumId;
        console.log(`Fetching album '${albumId}'`);
        fetchAlbumTracks(albumId);
    });
</script>

<Nav/>
<main class="flex flex-col gap-4 p-3 items-center min-w-screen min-h-screen bg-gray-950">
    <div class="lg:w-1/2 md:w-3/4">
        <div class="flex justify-center">
            <TrackList />
        </div>
    </div>
    <PlaylistModal/>
</main>
