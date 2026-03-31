<script lang="ts">
    import { onMount } from 'svelte';
    import { Play, Pause, SkipForward, SkipBack, Volume, Volume1, Volume2, VolumeX } from '@lucide/svelte';
    import { trackList } from '$lib/tracklist.svelte';
    import { togglePlayback, restartPlayback, playback, setGain, startPlaybackAt } from '$lib/playback.svelte';
    import { audioContext } from '$lib/stream';

    let dragging: boolean = false;
    let volumeSlider: HTMLDivElement;
    let seekSlider: HTMLDivElement;
    let offset: number = 0;

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

    function seek(e: MouseEvent) {
        if (!seekSlider) return;  

        const rect = seekSlider.getBoundingClientRect();
        const x = Math.max(0, Math.min(e.clientX - rect.left, rect.width));
        // we need to calculate the offset in seconds based on our progress across the bar
        offset = (x / rect.width) * playback.duration;
        if (playback.node) {
            playback.node.stop();
            playback.node = undefined;
        }
        if (playback.countInterval) {
            clearInterval(playback.countInterval);
        };
        playback.timeElapsed = offset;
        playback.startedAt = audioContext.currentTime - offset;
    };

    function setVolume(e: MouseEvent) {
        if (!volumeSlider) return;
        const rect = volumeSlider.getBoundingClientRect();
        const x = Math.max(0, Math.min(e.clientX - rect.left, rect.width));
        playback.volume = x / rect.width;
        playback.gain?.gain.setValueAtTime(playback.volume, audioContext?.currentTime);
        localStorage.setItem("volume", playback.volume.toString());
    };

    function mouseDown(e: MouseEvent) {
        console.log(seekSlider.contains(e.target));
        if (volumeSlider.contains(e.target)) {
            dragging = true;
            setVolume(e);
        }
        if (seekSlider.contains(e.target)) {
            playback.seeking = true;
            seek(e);
        }
    };

    function mouseMove(e: MouseEvent) {
        if (dragging) setVolume(e);
        if (playback.seeking) seek(e);
    };

    function mouseUp(e: MouseEvent) {
        dragging = false;
        if (playback.seeking) {
            playback.seeking = false;
            startPlaybackAt(offset);
        }
    };

    function toggleMute() {
        if (playback.muted) {
            playback.muted = false;
            playback.volume = 0.25;
            setGain(0.25);
        } else {
            playback.muted = true;
            playback.volume = 0;
            setGain(0);
        };
    };

    onMount(() => {
        playback.volume = localStorage.getItem("volume") as number || 1;
    });

</script>

<svelte:window
    on:mousemove={mouseMove}
    on:mouseup={mouseUp}
/>

<div class="flex rounded-lg w-200 h-18 bg-gray-800">
    <!-- Album Art -->
    <div class="rounded-lg w-18 h-18 bg-gray-700 aspect-square">
    </div>
    <!-- Seekbar -->
    <div class="w-full p-2 gap-2">
        <div class="flex text-white text-sm gap-2">
            {#if trackList.currentTrack.title }
                <p>{trackList.currentTrack.title}</p>
                <p>|</p><a class="text-lime-300" href="/tracks?album={trackList.currentTrack.artist.ID}">{trackList.currentTrack.album.Title}</a>
                <p>-</p><a class="text-lime-300" href="/tracks?artist={trackList.currentTrack.artist.ID}">{trackList.currentTrack.artist.Name}</a>
            {:else }
                <p>Happy luting!</p>
            {/if}
        </div>
        <div bind:this={seekSlider} 
            onmousedown={mouseDown}
            role="slider"
            aria-valuenow={playback.timeElapsed}
            aria-valuemin={0}
            aria-valuemax={1}
            tabindex="0"
            class="group h-2 rounded-full bg-slate-400 w-full shadow-sm/100"
        >
            <div class="relative h-2 rounded-full bg-lime-100" style={`width: ${(playback.timeElapsed / playback.duration) * 100}%`}>
                <div 
                    class="absolute -right-1 top-1/2 -translate-y-1/2 w-2 h-2 group-hover:scale-[2] bg-lime-400 rounded-full shadow-xs/40 group-hover:bg-lime-500 transition"
                    onmousedown={mouseDown}
                    role="slider"
                    aria-valuenow={playback.timeElapsed}
                    aria-valuemin={0}
                    aria-valuemax={1}
                    tabindex="0"
                >
                </div>
            </div>

        </div>
        <!-- Elapsed time -->
        <div class="px-4 py-1 float-right text-white text-sm">{formatTime(playback.timeElapsed)}/{formatTime(playback.duration) || "00:00:00"}</div>
        <!-- Volume -->
        <div class="mt-2 flex items-center gap-2">
            <button onclick={toggleMute} class="hover:cursor-pointer">
                {#if playback.volume <= 0 }
                    <VolumeX class="text-slate-300 fill-slate-300"/>
                {:else if playback.volume <= 0.25 }
                    <Volume class="text-slate-500 fill-slate-500"/>
                {:else if playback.volume <= 0.75 }
                    <Volume1 class="text-slate-500 fill-slate-500"/>
                {:else }
                    <Volume2 class="text-slate-500 fill-slate-500"/>
                {/if}
            </button>
            <div
                bind:this={volumeSlider}
                role="slider"
                aria-valuenow={playback.volume}
                aria-valuemin={0}
                aria-valuemax={1}
                tabindex="0"
                onmousedown={mouseDown}
                class="group h-2 rounded-full bg-slate-400 w-15 shadow-sm/100 hover:cursor-pointer"
            >
                <div class={`h-2 rounded-full bg-lime-100 relative ${dragging ? "" : "transition-all"}`} style={`width: ${playback.volume * 100}%`}>
                    <div 
                        class="absolute -right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 group-hover:scale-[1.25] bg-lime-400 rounded-full shadow-xs/40 group-hover:bg-lime-500 transition"
                        onmousedown={mouseDown}
                        role="slider"
                        aria-valuenow={playback.volume}
                        aria-valuemin={0}
                        aria-valuemax={1}
                        tabindex="0"
                    ></div>
                </div>
            </div>
        </div>
        <!-- Playback controls -->
        <div class="relative flex items-center w-full">
            <div class="absolute bottom-1/2 left-1/2 transform -translate-x-1/2 flex items-center">
                <button onclick={() => startPlaybackAt(0)} class="active:scale-[0.95] transition">
                    <SkipBack class="text-slate-500 fill-slate-500  hover:text-lime-500 hover:fill-lime-500 transition cursor-pointer" />
                </button>
                <button onclick={togglePlayback} class="mx-4 active:scale-[0.95] transition">
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
</div>
