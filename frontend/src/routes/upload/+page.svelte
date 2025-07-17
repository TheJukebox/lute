<script>
    import AudioPlayer from "../../AudioPlayer.svelte";
    import Upload from "../../Upload.svelte";
</script>

<div>
	<AudioPlayer></AudioPlayer>
</div>

<div class="upload-box">
    <Upload></Upload>
</div>

<style>
    .upload-box {
        padding-top: 200px;
    }
</style>
