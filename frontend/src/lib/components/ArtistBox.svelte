<script lang=ts>
    import { SquareUserRound, Shuffle } from "@lucide/svelte";
    export let artist;
</script>

<a href={`/tracks?artist=${artist.id}`} class="flex flex-col text-lime-200 hover:text-lime-400 hover:scale-[1.05] transition">
    <div class="flex flex-col w-80 rounded-xs aspect-square bg-gray-800 border-1 border-gray-600">
        <SquareUserRound class="w-full h-full aspect-square"/>
    </div>
    <div class="flex flex-col bg-gray-700">
        <p class="self-center">{artist.name}</p>
    </div>
</a>
