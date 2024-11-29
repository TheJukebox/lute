<script lang='ts'>
	let { src, title, artist } = $props();

	let time: number = $state(0);
	let duration: number = $state(0);
	let paused: boolean = $state(true);

	let mouseDown: boolean = false;

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

		const percentage: number = clamp(0, 100, Math.floor((relativePos / width) * 100));
		let seekFill: HTMLElement = seekbar.querySelector(".seekbar span") as HTMLElement;
		seekFill.style.width = `${percentage}%`;
	}

	function seek(event: MouseEvent) {
		if (mouseDown) {
			updateFill(event);
		}
	}
</script>
<svelte:window onmouseup={() => mouseDown = false} onmousemove={(event => seek(event))}></svelte:window>

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
	>
		<span class='seekbar' id="seekFill"></span>
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
		width: 0%;
		max-width: 100%;
		height: 100%;
		background-color: #15ff00;
	}

	.details {
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