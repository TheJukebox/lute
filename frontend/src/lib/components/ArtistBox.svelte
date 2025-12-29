<script lang=ts>
    import { SquareUserRound, Shuffle } from "@lucide/svelte";
    export let artist;
</script>

<a href={`/tracks?artist=${artist.id}`} class="text-lime-600 hover:text-green-600">
    <div class="p-2 flex flex-col md:w-80 md:h-80 w-50 h-50 items-center justify-center hover:scale-[1.05] hover:cursor-pointer transition">
        <div class="flex flex-col md:w-65 md:h-65 w-40 h-40 bg-lime-200 rounded-xl shadow-sm aspect-square items-center justify-center">
            <SquareUserRound class="text-lime-100" size="256" />
        </div>
        <div class="p-2 font-bold text-sm text-center text-clip text-ellipsis overflow-hidden whitespace-no-wrap">
            {artist.name}
        </div>
    </div>
</a>
