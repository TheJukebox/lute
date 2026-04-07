<script lang="ts">
    import { onMount } from 'svelte';
    import { ListPlus, ListRestart } from "@lucide/svelte";
    const playlistList = $state({playlists: []});

    async function createPlaylist() {
        const response: Response = await fetch(
            "http://localhost:7001/playlists",
            {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    "type": "create",
                    "title": "new playlist",
                    //"id": "661bd448-ce59-43ef-9111-67b89eb5383d",
                    //"track_id": "6ab52a20-5876-4a5c-a00d-d71f835616eb",
                    "description": "test playlist",
                }),
            },
        );
        if (!response.ok) {
            console.error('Failed to create new playlist');
            return;
        }
        console.log(await response.json());
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

    onMount(async () => {
        fetchPlaylists();
    });
</script>

<div class="w-1/2 h-full bg-gray-800 border-1 border-gray-700 rounded-sm shadow-sm p-5 text-white">
    <div class="flex gap-5">
        <button data-sveltekit-reload class="flex gap-1 p-2 hover:text-lime-400 hover:cursor-pointer transition rounded-sm bg-gray-700 border-gray-600 border-1" onclick={createPlaylist}><ListPlus/> New playlist</button>
        <button class="flex gap-1 p-2 hover:text-lime-400 hover:cursor-pointer transition rounded-sm bg-gray-700 border-gray-600 border-1" onclick={fetchPlaylists}><ListRestart/>Refresh playlists</button>
    </div>
    <div class="flex flex-col p-2 gap-1">
        {#each playlistList.playlists as playlist}
            <a data-sveltekit-reload class="hover:text-lime-400 transition" href="/playlists?playlist={playlist.id}">{playlist.title}</a>
        {/each}
    </div>
</div>
