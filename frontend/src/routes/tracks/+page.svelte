<script lang="ts">
    import { onMount } from 'svelte';
    import Controls from '$lib/components/Controls.svelte';
    import TrackList from '$lib/components/TrackList.svelte';
    import { fetchArtistTracks, trackList, fetchArtist } from '$lib/tracklist.svelte';
    import ArtistBox from '$lib/components/ArtistBox.svelte';

    let artistID: string | null;
    let artist; 
    
    onMount(async () => {
        const params: URLSearchParams = new URLSearchParams(window.location.search);
        artistID = params.get("artist");
        await fetchArtistTracks(artistID);
        artist = await fetchArtist(artistID);
        console.log(artist)
    });

</script>

<main class="flex flex-col gap-4 p-3 items-center min-w-screen min-h-screen bg-emerald-50">
    <div class="lg:w-1/2 md:w-3/4">
        <div class="flex justify-center">
            <Controls />
        </div>
        <div class="flex justify-left">
            {#if artist }
                <ArtistBox artist={artist} />
            {:else }
                <div class="animate-pulse flex md:w-65 md:h-65 w-40 h-40 bg-slate-200 rounded-lg shadow-sm aspect-square items-center justify-center">
                </div>
            {/if}
        </div>
        <div class="flex justify-center">
            <TrackList />
        </div>
    </div>
</main>
