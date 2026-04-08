<script lang="ts">
    import { ListPlus, ListRestart, CircleCheck } from "@lucide/svelte";
    import { tick } from 'svelte';

    const playlistList = $state({playlists: []});

    let { showPlaylistModal = $bindable(), addTrackId = $bindable() } = $props();
    let dialog = $state();
    let nameInput = $state();
    let newPlaylistName = $state('New playlist');
    let showNameInput = $state(false);

    $effect(() => {
        fetchPlaylists();
    });

    async function addToPlaylist(playlistId: string) {
        const response: Response = await fetch(
            "http://localhost:7001/playlists",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    "type": "add",
                    "id": playlistId,
                    "track_id": addTrackId,
                }),
            },
        );
        if (!response.ok) {
            console.error('Failed to add track to playlist');
            return;
        }
        console.log(`Added track '${addTrackId}' to playlist '${playlistId}'.`);
    }

    async function createPlaylist(e, newPlaylistName) {
        e.stopPropagation();
        showNameInput = false;
        const response: Response = await fetch(
            "http://localhost:7001/playlists",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    "type": "create",
                    "title": newPlaylistName, 
                    "description": "test playlist",
                }),
            },
        );
        if (!response.ok) {
            console.error('Failed to create new playlist');
            return;
        }
        fetchPlaylists();
    }

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

    async function popNameInput() {
        showNameInput = true;
        await tick();
        nameInput.focus();
    };

    function closeModal() {
        showPlaylistModal = false;
        showNameInput = false;
    }
</script>

{#if showPlaylistModal }
<dialog
    bind:this={dialog}
    onclose={closeModal}
    onclick={(e) => { if (e.target === dialog) closeModal(); }}
    class="w-full h-full fixed inset-0 backdrop-blur-sm bg-black/30 flex items-center justify-center "
>
    <div class="w-200 bg-gray-800 border-1 border-gray-700 rounded-sm shadow-sm p-5 text-white">
        <div class="flex gap-5">
            {#if !showNameInput}
            <button data-sveltekit-reload class="flex gap-1 p-2 hover:text-lime-400 hover:cursor-pointer transition rounded-sm bg-gray-700 border-gray-600 border-1" onclick={popNameInput}><ListPlus/> New playlist</button>
            {/if}
            {#if showNameInput}
            <div class="flex gap-2">
                <input bind:this={nameInput} class="w-full h-full rounded-sm bg-gray-700 text-lime-400" type="text" bind:value={newPlaylistName}/>
                <button class="text-lime-800 hover:text-lime-400 hover:cursor-pointer transition" onclick={(e) => createPlaylist(e, newPlaylistName)}><CircleCheck/></button>
            </div>
            {/if}
        </div>
        <div class="flex flex-col p-2 gap-1 items-start">
            {#each playlistList.playlists as playlist}
                <button class="hover:text-lime-400 transition" onclick={() => addToPlaylist(playlist.id)}>{playlist.title}</button>
            {/each}
        </div>
    </div>
</dialog>
{/if}
