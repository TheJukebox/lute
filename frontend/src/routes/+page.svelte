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
	<ul class='menu'>
		<a href='/library'>
			<li>
				<img src='./assets/library.svg' alt='Your library'/>
				<span>Your Library</span>
			</li>
		</a>
	</ul>
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

	.menu {
		color: var(--indigo-dye);
		font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
		font-size: x-large;
		font-weight: bolder;
		margin-left: 25%;
		display: flex;
		align-items: center;
		gap: 20px;
	}

	.menu a {
		text-decoration: none;
		color: var(--viridian);
	}

	.menu a:active {
		color: var(--indigo-dye);
	}

	.menu li {
		list-style: none;
		line-height: 24px;
	}

	.menu img{
		filter: invert(100%);	
	}
</style>