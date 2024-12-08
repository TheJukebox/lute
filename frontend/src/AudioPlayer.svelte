<script lang='ts'>
	import stream from '$lib/gen/stream_grpc_web_pb';
	import '$lib/audio_processing'
	import { fetchStream, playFromBuffer, togglePlayback } from '$lib/audio_processing';
	let { src, title, artist } = $props();

	let time: number = $state(0);
	let duration: number = $state(0);
	let paused: boolean = $state(true);

	let mouseDown: boolean = false;

	function startStream() {
		paused = !paused
		fetchStream('http://127.0.0.1:8080', '../../output.aac', 'test-session');
		togglePlayback();
	}
	
	function toggle() {
		paused = !paused
		togglePlayback();
	}

	function format(time: number): string {
		if (isNaN(time)) return '...';

		const minutes: number = Math.floor(time / 60);
		const seconds: number = Math.floor(time % 60);

		return `${minutes}:${seconds < 10 ? `0${seconds}` : seconds}`;
	}

	function clamp(min: number, max: number, x: number) {
		return Math.min(Math.max(x, min), max);
	}

	function updateFill(event: MouseEvent) {
		let seekbar: HTMLElement = document.querySelector(".seekbar") as HTMLElement;
		let bounds: DOMRect = seekbar.getBoundingClientRect();	
		let relativePos: number = event.clientX - bounds.left;

		let style = window.getComputedStyle(seekbar);
		let width = parseFloat(style.width);

		let percentage: number = clamp(0, 100, Math.floor((relativePos / width) * 100));
		let seekFill: HTMLElement = seekbar.querySelector(".seekbar span") as HTMLElement;
		seekFill.style.width = `${percentage}%`;
		seekFill.ariaValueNow = `${percentage}`;
	}

	function seek(event: MouseEvent) {
		if (mouseDown) {
			updateFill(event);
		}
	}
</script>
<svelte:window 
	onmouseup={() => mouseDown = false} 
	onmousemove={(event => seek(event))}
></svelte:window>

<div class='player' class:paused>
	<audio>
	</audio>
	<div class='albumArt'>
	</div>
	<div class='details'>
		<p class='title'>Weezer Does Weezer</p>
		<p class='artist'>Weezer</p>
		<p class='album'>Weezer: The Weezer Years</p>
	</div>
	<div class='seekbar' id="seekbar"
		onmousedown={() => mouseDown = true }
		onmouseup={() => mouseDown = false }
		onmousemove={(event) => seek(event)}
		role="slider"
		aria-valuenow=0
		aria-valuemin=0
		aria-valuemax=100
		tabindex=0
	>
		<span class='seekbar' id="seekFill" style="width: 0%">
			<div class='playhead' id="playhead"></div>
		</span>
	</div>
	<div class='controls'>
		<button 
			class='previous'
			aria-label='previous'
			onclick={startStream}
		>prev</button>
		<button 
			class='pause'
			onclick={toggle}
			aria-label={paused ? 'play' : 'pause'}
		></button>
		<button 
			class='next'
			aria-label='next'
		>next</button>
	</div>
</div>

<style>
	.player {
		position: relative;
		border: 2px dashed #000;
		width: 600px;
		height: 100px;
		background-color: #f9f9f9;
	}

	.controls {
		position: absolute;
		display: flex;
		width: 200px;
		height: 30px;
		background-color:#a0a0a0;
		border: 2px dashed #000;	
		left: 200px;
		top: 105px;
		justify-content: center;
		align-items: center;
	}

	.pause {
		width: 15%;
		aspect-ratio: 1;
		background: none;
		background-repeat: no-repeat;
		background-position: 50% 50%;
		border-radius: 50%;
	}

	[aria-label="pause"] {
		background-image: url(./assets/pause.svg);
	}

	[aria-label="play"] {
		background-image: url(./assets/play.svg);
	}

	.albumArt {
		position: absolute;
		height: 90%;
		width: 90px;
		border: 0px solid #15ff00;
		background-color: #ff0000;
	}
	
	.seekbar {
		position: absolute;
		bottom: 0;
		height: 10%;
		width: 100%;
		border: 0px solid #00bbff;
		background-color: #8b8b8b;
	}

	.seekbar span{
		min-width: 0;
		max-width: 100%;
		height: 100%;
		background-color: #15ff00;
	}

	.playhead {
		position: absolute;
		right: 0px;
		height: 100%;
		width: 5px;
		background-color: #252525;
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
		background-color: #ffffff;
	}
	
	.title {
		font-size: medium;
		margin: 2px 10px 0 2px;
	}
	
	.artist {
		font-size: medium;
		margin: 0 10px 0 2px;
	}

	.album {
		position: absolute;
		bottom: 0;
		font-size: small;
		margin: 0;
		font-style: oblique;
	}
</style>