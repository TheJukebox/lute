<script>
	import AudioPlayer from "../AudioPlayer.svelte";
</script>
<div>
	<AudioPlayer></AudioPlayer>
</div>
<div class="content">
	<div class="big_logo">
		<img src="./assets/logo_lute.svg" alt="lute logo"/>
		<h1>LUTE</h1>
	</div>
	<p>
		You stand in a vast chamber. A Lute rests upon a pedestal. Play it?
	</p>
</div>

<style>
	.big_logo {
		text-align: center;
	}
	.big_logo img {
		position: relative;
		margin-bottom: 0px;
	}
	.big_logo h1 {
		position: relative;
		margin-top: 0px;
		font-family:'Franklin Gothic Medium', 'Arial', Arial, sans-serif;
		font-size: 72pt;
		color: var(--goldenrod);
	}
	.content p {
		font-family: 'Segoe UI Bold', Tahoma, Geneva, Verdana, sans-serif;
		display: block;
		margin-left: 25%;
	}
	.content {
		margin-top: 10%;
	}
</style>