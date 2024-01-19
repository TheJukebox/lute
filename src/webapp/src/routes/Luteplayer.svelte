<script lang="ts">
    export let current_audio: string;
</script>

<!-- Markup -->
<audio id="luteplayer" controls>
    <source src={current_audio}>
    Your browser does not support the audio player.
</audio>

<style>

</style>