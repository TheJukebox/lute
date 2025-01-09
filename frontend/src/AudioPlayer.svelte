<script lang='ts'>
	import { onDestroy } from 'svelte';
	import { currentTrack, currentTime, seekToTime, isPlaying } from '$lib/audio_store';
	import { setVolume } from '$lib/stream_handler';

	import type { Track } from '$lib/audio_store'
	import type { Unsubscriber } from 'svelte/store';
	import { togglePlayback } from '$lib/stream_handler';

	let mouseDown: boolean = false;	
	let currentVolume: number = 1;
	let cachedVolume: number = 1;
	let muted: boolean = false;

	let track: Track = {
		path: '',
		title: '',
		artist: '',
		album: '',
		paused: true,
		duration: 0,
	};

	let time = "0:00";
	let maxTime = "0:00";
	let fillPercentage = 0;

	const unsubTime = currentTime.subscribe((value) => {
		time = formatSeconds(value);
		fillPercentage = (value / track.duration) * 100;
	});

	const unsubTrack: Unsubscriber = currentTrack.subscribe((value) => {
		track = value;
		maxTime = formatSeconds(track.duration);
	});

	let playing: boolean = false;
	const unsubPlay: Unsubscriber = isPlaying.subscribe((value) => {
		playing = value;
		togglePlayback();
	});


	onDestroy(() => {
		unsubPlay();
		unsubTrack();
		unsubTime();
	});

	/**
	 * Toggles playback of the stream.
	 * @function
	 */
	function toggle(): void {
		playing = !playing;
		isPlaying.set(playing);
	}

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

	function clamp(min: number, max: number, x: number): number {
		return Math.min(Math.max(x, min), max);
	}

	/**
	 * Updates the fill of the seekbar element
	 * @function
	 * @param event	{MouseEvent}	The mouse event that triggered this function.
	 */
	function updateFill(event: MouseEvent): void {
		let seekbar: HTMLElement = document.querySelector(".seekbar") as HTMLElement;
		let bounds: DOMRect = seekbar.getBoundingClientRect();	
		let relativePos: number = event.clientX - bounds.left;

		let style = window.getComputedStyle(seekbar);
		let width = parseFloat(style.width);

		let percentage: number = clamp(0, 100, Math.floor((relativePos / width) * 100));
		let seekFill: HTMLElement = seekbar.querySelector(".seekbar span") as HTMLElement;
		seekFill.style.width = `${percentage}%`;
		seekFill.ariaValueNow = `${percentage}`;
		seekToTime((percentage / 100) * track.duration);
	}
	

	/**
	 * Convenience function for discarding mouse events.
	 * @function
	 * @param event
	 */
	function filterEvents(event: MouseEvent): void {
		if (mouseDown) {
			updateFill(event);
		}
	}

	function toggleVolume(): void {
		muted = !muted;
		if (muted) {
			cachedVolume = currentVolume;
			console.log(cachedVolume);
			currentVolume = 0;
			setVolume(0);
		} else {
			currentVolume = cachedVolume;
			setVolume(cachedVolume);
		}
	}

	function updateVolume(event: Event): void {
		const target = event.target as HTMLInputElement;
		currentVolume = parseFloat(target.value);
		setVolume(target.value);
	}

</script>
<svelte:window 
	onmouseup={() => mouseDown = false} 
	onmousemove={(event => filterEvents(event))}
></svelte:window>

<div class='banner'>
	<a href="/" class="logo">
		<img src="./assets/logo_lute.svg" class="logo" alt="Home"/>
		<span><h1>LUTE</h1></span>
	</a>

	<div class='player' class:playing>
		<audio>
		</audio>

		<div class='albumArt'>
			<span>?</span>
		</div>

		<div class='details'>
			<p class='title'>{track.title}</p>
			<p class='artist'>{track.artist}</p>
			<p class='album'>{track.album}</p>
		</div>

		<div class='time'>
			{time}/{maxTime}
		</div>

		<div class='seekbar' id="seekbar"
			onmousedown={() => mouseDown = true }
			onmouseup={() => mouseDown = false }
			onmousemove={(event) => filterEvents(event)}
			role="slider"
			aria-valuenow=0
			aria-valuemin=0
			aria-valuemax=100
			tabindex=0
		>
			<span class='seekbar' id="seekFill" style="width: {fillPercentage}%">
				<div class='playhead' id="playhead"></div>
			</span>
		</div>

		<div class='streamControl'>
			<button 
				class='previous'
				aria-label='previous'
				onclick={() => seekToTime()}	
			></button>
			<button 
				class='pause'
				onclick={toggle}
				aria-label={playing ? 'pause' : 'play'}
			></button>
			<button 
				class='next'
				aria-label='next'
			></button>
		</div>
		<div class="volume-control">
			<button class='mute' aria-label='mute' onclick={toggleVolume}></button>
			<input id="volume" class="volume-slider" type="range" min="0" max="1" step="0.01" value={currentVolume} oninput={updateVolume}/>
		</div>
	</div>
</div>


<style>
	:root {
		--indigo-dye: #124e78ff;
		--goldenrod: #d5a021ff;
		--alabaster: #ede7d9ff;
		--viridian: #6a8e7fff;
		--viridian-dark: rgb(81, 121, 104);	
		--viridian-darker: rgb(61, 91, 78);	
		--bright-pink-crayola: #ea526fff;
	}

	.banner {
		background-color: var(--viridian);
		position: fixed;
		top: 0;
		left: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		z-index: 1;
		width: 100%;
		height: 100px;
		box-shadow: 0 0px 10px rgba(0, 0, 0, 0.5);
	}

	.logo img{
		position: absolute;
		left: 0;
		top: 5%;
		height: 100px;
		width: auto;
		transform-origin: top right;
		user-select: none; 
	}

	.logo span{
		position: absolute;
		font-family: 'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;
		font-size: larger;
		color: var(--goldenrod);
		left: 20px;
		top: 38%;
		transform-origin: top right;
		text-shadow: 0px 0px 8px rgba(0, 0, 0, 256);
		user-select: none; 
	}

	.logo span:active {
		color: var(--indigo-dye);
	}

	.player {
		position: relative;
		width: 600px;
		height: 90px;
		background-color: var(--viridian-dark);
		border-radius: 1px;
		z-index: 2;
		box-shadow: 0px 2px 10px rgba(0, 0, 0, 0.3);
	}

	.streamControl {
		position: absolute;
		display: flex;
		width: 200px;
		height: 30px;
		left: 200px;
		top: 98px;
		justify-content: center;
		align-items: center;
		z-index: 10;
		background-color: var(--bright-pink-crayola);
		border-radius: 10px;
		box-shadow: 0px 2px 10px rgba(0, 0, 0, 0.3);
	}

	.streamControl button:hover {
		filter: invert(75%);
		cursor: pointer;
	}

	.mute {
		width: 5%;
		aspect-ratio: 1;
		background: none;
		background-repeat: no-repeat;
		background-position: 50% 50%;
		border-radius: 50%;
		border-width: 0px;
		background-image: url(./assets/volume.svg);
	}

	.pause {
		width: 15%;
		aspect-ratio: 1;
		background: none;
		background-repeat: no-repeat;
		background-position: 50% 50%;
		border-color: var(--goldenrod);
		border-width: 0px;
	}

	.next {
		width: 15%;
		aspect-ratio: 1;
		background: none;
		background-repeat: no-repeat;
		background-position: 50% 50%;
		border-radius: 50%;
		border-color: var(--goldenrod);
		border-width: 0px;
		background-image: url(./assets/skip_next.svg);
	}

	.previous {
		width: 15%;
		aspect-ratio: 1;
		background: none;
		background-repeat: no-repeat;
		background-position: 50% 50%;
		border-radius: 50%;
		border-color: var(--goldenrod);
		border-width: 0px;
		background-image: url(./assets/skip_previous.svg);
	}

	[aria-label="pause"] {
		background-image: url(./assets/pause.svg);
	}

	[aria-label="play"] {
		background-image: url(./assets/play.svg);
	}

	.albumArt {
		position: absolute;
		height: 88%;
		width: 15%;
		margin: 1px 0px 2px 2px;
		background-color: var(--viridian-darker);
		border-radius: 1px;
	}
	
	.albumArt span{
		font-family: 'Franklin Gothic Medium', 'Arial Narrow', Arial, sans-serif;
		font-size: 50pt;
		display: table;
		margin: 0 auto;	
		margin-top: 5%;
		color: var(--goldenrod);
		user-select: none; 
	}
	.seekbar {
		position: absolute;
		bottom: 0px;
		height: 10%;
		width: 100%;
		background-color: var(--viridian-darker);
	}

	.seekbar span{
		min-width: 0;
		max-width: 100%;
		height: 100%;
		background-color: var(--bright-pink-crayola);
		box-shadow: 0 0 2px rgba(255, 0, 191, 0.4);
	}

	.playhead {
		position: absolute;
		right: -4px;
		height: 100%;
		width: 5px;
		background-color: var(--goldenrod);
		border-radius: 2px;
		border-width: 0px;
	}

	.details {
		-webkit-user-select: none; /* Safari */
  		-ms-user-select: none; /* IE 10 and IE 11 */
  		user-select: none; /* Standard syntax */
		position: absolute;
		left: 90px;
		bottom: 10%;
		top: 0;
		justify-content: start;
		border-radius: 5px;
	}
	
	.title {
		font-family: 'Segoe UI Bold', Tahoma, Geneva, Verdana, sans-serif;
		font-size: medium;
		margin: 5px 10px 0 5px;
		color: var(--indigo-dye);
	}
	
	.artist {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		color: var(--indigo-dye);
		font-size: medium;
		margin: 0px 10px 2px 5px;
	}

	.album {
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		color: var(--indigo-dye);
		position: absolute;
		bottom: 0;
		font-size: small;
		margin: 0px 10px 2px 5px;
		font-style: oblique;
	}

	.time {
		font-family: 'Segoe UI Bold', Tahoma, Geneva, Verdana, sans-serif;
		font-size: small;
		z-index: 100px;
		position: relative;
		margin: 60px 0px 50px 89%;
		color: var(--indigo-dye);
		user-select: none; 
	}

	.volume-control {
		display: flex;
		align-items: center;
		margin-top: 10px;
	}

	.volume-slider {
		appearance: none;
		width: 15%;
		height: 10px;
		background-color: var(--viridian-dark);
		border-radius: 5px;
		cursor: pointer;
		outline: none;
	}

	.volume-slider::-webkit-slider-thumb {
		appearance: none;
		width: 14px;
		height: 18px;
		border-radius: 5px;
		background: var(--goldenrod);
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
		border-width: 0;
		cursor: pointer;
	}

	.volume-slider::-moz-range-thumb {
		width: 14px;
		height: 18px;
		border-radius: 5px;
		background: var(--goldenrod);
		box-shadow: 0 0 10px rgba(0, 0, 0, 0.4);
		border-width: 0;
		cursor: pointer;
	}

	.volume-slider::-ms-thumb {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		background: var(--goldenrod);
		cursor: pointer;
	}

	.volume-slider::-webkit-slider-runnable-track {
		background-color: var(--bright-pink-crayola);
		height: 100%;
		border-radius: 5px;
		box-shadow: 0 0 5px rgba(255, 0, 191, 1);
	}

	.volume-slider::-moz-range-progress {
		background-color: var(--bright-pink-crayola);
		height: 100%;
		border-radius: 5px;
		box-shadow: 0 0 5px rgba(255, 0, 191, 1);
	}

	.volume-slider::-moz-range-track {
		background-color: var(--viridian-dark);
		height: 100%;
		border-radius: 5px;
	}

	.volume-slider::-ms-fill-lower {
		background-color: var(--bright-pink-crayola);
		height: 100%;
		border-radius: 5px;
		box-shadow: 0 0 5px rgba(255, 0, 191, 1);
	}

	.volume-slider::-ms-fill-upper {
		background-color: var(--viridian-dark);
		height: 100%;
		border-radius: 5px;
	}
</style>