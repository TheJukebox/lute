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
    
<div
    class="grid grid-cols-[auto_1fr] gap-4 p-2 shadow bg-lime-50 rounded-xl w-full max-w-5xl"
>
    <div class="flex p-2 w-32 h-32 sm:w-32 sm:h-32 md:w-48 md:h-48 shadow rounded-lg bg-lime-100 aspect-square">
        <div class="inset-shadow-xs/20 p-2 min-w-full min-h-full bg-lime-200 rounded-sm"></div>
    </div>
    <div class="flex flex-col gap-5 p-2 min-w-full shadow rounded-lg bg-lime-100 items-center">
        <div class="inset-shadow-sm/20 bg-lime-200 rounded-lg p-2 min-w-full min-h-30 transition">
            <div bind:this={seekSlider} 
                onmousedown={mouseDown}
                role="slider"
                aria-valuenow={playback.timeElapsed}
                aria-valuemin={0}
                aria-valuemax={1}
                tabindex="0"
                class="h-2 rounded-full bg-slate-400 w-full shadow-sm/100"
            >
                <div class="relative h-2 rounded-full bg-lime-100" style={`width: ${(playback.timeElapsed / playback.duration) * 100}%`}>
                    <div 
                        class="absolute -right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 hover:scale-[1.25] bg-lime-400 rounded-full shadow-xs/40 hover:bg-lime-500"
                        onmousedown={mouseDown}
                        role="slider"
                        aria-valuenow={playback.timeElapsed}
                        aria-valuemin={0}
                        aria-valuemax={1}
                        tabindex="0"
                    ></div>
                </div>
            </div>
            <div class="px-4 py-1 float-right text-lime-800 text-sm">{formatTime(playback.timeElapsed)}/{formatTime(playback.duration) || "00:00:00"}</div>
            <div class="flex flex-col items-center justify-center text-center w-full transition">
                <div class="grid grid-cols-3 gap-0 font-semibold w-full text-lime-800">
                    <span>{trackList.currentTrack.Name}</span>
                    <span >{trackList.currentTrack.Artist}</span> 
                    <span>{trackList.currentTrack.Album}</span>
                </div>
            </div>
        </div>
        <div class="relative flex items-center w-full">
            <div class="ml-15 flex items-center gap-2">
                <button onclick={toggleMute}>
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
                    class="h-2 rounded-full bg-slate-400 w-30 shadow-sm/100"
                >
                    <div class={`h-2 rounded-full bg-lime-100 relative ${dragging ? "" : "transition-all"}`} style={`width: ${playback.volume * 100}%`}>
                        <div 
                            class="absolute -right-2 top-1/2 -translate-y-1/2 w-3.5 h-3.5 hover:scale-[1.25] bg-lime-400 rounded-full shadow-xs/40 hover:bg-lime-500"
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
            <div class="absolute left-1/2 transform -translate-x-1/2 flex items-center gap-4">
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
