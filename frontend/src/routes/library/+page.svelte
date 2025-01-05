<script lang='ts'>
    import AudioPlayer from '../../AudioPlayer.svelte';
    import { startStream } from '../../audio_store';
    
    import type { Track } from '../../audio_store';

    type SongData = {
        id: number;
        title: string;
        artist: string;
        album: string;
        num: number;
        path: string;
    }

    // some fake song data
    let songs: SongData[] = [
        { id: 1, title: 'Something In The Way', artist: 'Nirvana', album: 'Nevermind', num: 1, path: "uploads/conmverted/SomethingInTheWay.aac" },
        { id: 2, title: 'song2', artist: 'artist2', album: 'album', num: 1, path: "" },
        { id: 3, title: 'song3', artist: 'artist3', album: 'album', num: 1, path: ""},
    ];
</script>
  

<div>
    <AudioPlayer></AudioPlayer>
</div>

<div class='title'>
    <h1>Your Library</h1>
</div>
<div class='library_container'>
    <table class='song_table'>
        <thead class='table_header'>
            <tr>
                <th>Title</th>
                <th>Length</th>
                <th>Artist</th>
                <th>Album</th>
                <th>Genre</th>
            </tr>
        </thead>
        <tbody>
            {#each songs as song}
                <tr class='table_entry' on:click={() => startStream("uploads/converted/SomethingInTheWay.aac", song.title, song.artist, song.album)}> 
                    <td>{song.title}</td>
                    <td>00:00</td>
                    <td>{song.artist}</td>
                    <td>{song.album}</td>
                    <td></td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>

<style>
	:root {
		--indigo-dye: #124e78ff;
		--goldenrod: #d5a021ff;
		--alabaster: #ede7d9ff;
		--alabaster-dark: rgb(234, 228, 215);
		--viridian: #6a8e7fff;
		--viridian-dark: rgb(81, 121, 104);	
		--viridian-darker: rgb(61, 91, 78);	
		--bright-pink-crayola: #ea526fff;
	}

    .title {
        margin-top: 12%;
        margin-left: 5%;
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
        color: var(--viridian-dark);
    }
    .library_container {
        font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
    }

    .song_table {
        width: 90%;
        border-spacing: 0;
        margin: auto;
    }

    .table_header {
        color: var(--viridian);
    }

    .table_header th {
        padding-left: 0.75rem;
        text-align: left;
        border-bottom: 2px solid var(--viridian);
    }

    .table_entry {
        background-color: var(--alabaster);
        transition: background-color 0.2s ease;
    }

    .table_entry:nth-child(even) {
        background-color: var(--alabaster-dark);
    }

    .table_entry:hover {
        color: var(--viridian);
        cursor: pointer;
    }

    .table_entry:active {
        color: var(--indigo-dye);
    }

    .table_entry td {
        padding: 0.75rem;
        text-align: left;
    }
</style>