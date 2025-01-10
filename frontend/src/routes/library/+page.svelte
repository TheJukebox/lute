<script lang='ts'>
    import AudioPlayer from '../../AudioPlayer.svelte';
    import { startStream } from '$lib/audio_store';
    
    import type { Track } from '$lib/audio_store';

    // some fake track data
    let tracks: Track[] = [
        { title: 'Something In The Way', artist: 'Nirvana', album: 'Nevermind', path: "uploads/converted/SomethingInTheWay.aac", duration: 235 },
        { title: 'Breezeblocks', artist: 'alt-j', album: 'alt-j', path: "uploads/converted/Breezeblocks.aac", duration: 226 },
    ];

	/**
	 * Takes a time in seconds and converts it to a string in the format MM:SS
	 * @function format
	 * @param time {number}	A length of time in seconds.
	 * @returns {string}	Formated time or '...'.
	 */
	function formatSeconds(time: number): string {
		if (isNaN(time)) return '...';

		const minutes: number = Math.floor(time / 60);
		const seconds: number = Math.floor(time % 60);

		return `${minutes}:${seconds < 10 ? `0${seconds}` : seconds}`;
	}
</script>
  

<div>
    <AudioPlayer></AudioPlayer>
</div>

<div class='title'>
    <h1>Your Library</h1>
</div>
<div class='library_container'>
    <table class='track_table'>
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
            {#each tracks as track}
                <tr class='table_entry' onclick={() => startStream(track)}> 
                    <td>{track.title}</td>
                    <td>{formatSeconds(track.duration)}</td>
                    <td>{track.artist}</td>
                    <td>{track.album}</td>
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

    .track_table {
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