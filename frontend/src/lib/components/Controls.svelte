<script lang="ts">
    import { onMount } from 'svelte';
    import { Play, Pause, SkipForward, SkipBack } from '@lucide/svelte';
    import { trackList } from '$lib/tracklist.svelte';
    import { togglePlayback, playback } from '$lib/playback.svelte';

    function formatTime(seconds) {
        const h = Math.floor(seconds / 3600);
        const m = Math.floor((seconds % 3600) / 60);
        const s = Math.floor(seconds % 60);

        return [
            String(h).padStart(2, '0'),
            String(m).padStart(2, '0'),
            String(s).padStart(2, '0'),
        ].join(':');
    };
</script>
    
<div
    class="grid grid-cols-[auto_1fr] gap-4 p-2 shadow bg-lime-50 rounded-xl w-full max-w-5xl"
>
    <div class="flex p-2 w-32 h-32 sm:w-32 sm:h-32 md:w-48 md:h-48 shadow rounded-lg bg-lime-100 aspect-square">
        <div class="inset-shadow-xs/20 p-2 min-w-full min-h-full bg-lime-200 rounded-sm"></div>
    </div>
    <div class="flex flex-col gap-5 p-2 min-w-full shadow rounded-lg bg-lime-100 items-center">
        <div class="inset-shadow-sm/20 bg-lime-200 rounded-lg p-2 min-w-full min-h-30 transition">
            <div class="h-2 rounded-full bg-slate-400 w-full shadow-sm/100">
                <div class="h-2 rounded-full bg-lime-100" style={`width: ${(playback.timeElapsed / playback.duration) * 100}%`}></div>
            </div>
            <div class="px-4 py-1 float-right text-lime-800 text-sm">{formatTime(playback.timeElapsed)}/{formatTime(playback.duration) || "00:00:00"}</div>
            <div class="flex flex-col items-center justify-center text-center w-full">
                <div class="grid grid-cols-3 gap-0 font-semibold w-full text-lime-800">
                    <span>{trackList.currentTrack.Name}</span>
                    <span >{trackList.currentTrack.Artist}</span> 
                    <span>{trackList.currentTrack.Album}</span>
                </div>
            </div>
        </div>
        <div class="grid grid-cols-3 gap-4">
            <button>
                <SkipBack class="text-slate-500 fill-slate-500  hover:text-lime-500 hover:fill-lime-500 transition cursor-pointer" />
            </button>
            <button onclick={togglePlayback} class="active:scale-[0.95] transition">
            {#if playback.playing}
                <Pause class="text-slate-500 fill-slate-500  hover:text-lime-500 hover:fill-lime-500 transition cursor-pointer"/>
            {:else}
                <Play class="text-slate-500 fill-slate-500  hover:text-lime-500 hover:fill-lime-500 transition cursor-pointer"/>
            {/if}
            </button>
            <button>
                <SkipForward class="text-slate-500 fill-slate-500  hover:text-lime-500 hover:fill-lime-500 transition cursor-pointer"/>
            </button>
        </div>
    </div>
</div>
