<script lang="ts">
    import { onMount } from 'svelte';
    import { 
        uploadTrack,
    } from '$lib/upload';
    import type {
        Track
    } from '$lib/upload';
    import { parseBlob, parseStream } from 'music-metadata';
    import type {
        IAudioMetadata
    } from 'music-metadata';
    import { fetchTracks} from '$lib/tracklist.svelte.ts';

    let uploading: boolean = false;
    let uploadName: string = "";
    let uploadingCount: number = 0;
    let uploadCurrent: number = 0;

    onMount(() => {
        window.addEventListener("dragover", (e) => {
            const fileItems = [...e.dataTransfer.items].filter(
                (item) => item.kind === "file",
            );
            if (fileItems.length > 0) {
                e.preventDefault();
            }
        });
    });

    const onDrop = async (e: DragEvent) => {
        if ([...e.dataTransfer.items].some((item) => item.kind === "file")) {
            e.preventDefault();
            uploading = true;
            console.log(e);
            uploadingCount = e.dataTransfer.files.length;
            uploadCurrent = 1;
            for (const file of e.dataTransfer.files) {

                const metadata: IAudioMetadata = await parseBlob(file);
                const name: string = metadata.common.title || file.name;
                const album: string = metadata.common.album || "Unknown";
                const artist: string = metadata.common.artist || "Unknown";
                const date: string = metadata.common.date || ""; 
                const release: string = metadata.common.originaldate || ""; 
                const number: number = metadata.common.track.no || 1;
                const disk: number = metadata.common.disk.no || 1;
                uploadName = name;

                const track: Track = {
                    name: name,
                    uriName: encodeURIComponent(`${artist}/${album}/${name}`),
                    contentType: file.type,
                    artist: artist,
                    album: album,
                    date: date,
                    release: release,
                    trackNumber: number,
                    diskNumber: disk,
                };
                console.debug(`Uploading track: ${JSON.stringify(track)}`);
                uploadTrack(track, file);
                uploadCurrent++;
            }
            uploading = false;
            uploadCurrent = 0;
            fetchTracks();
        }
    };
</script>

<div 
    class={`relative p-4 min-h-30 max-w-100 rounded-xl shadow hover:bg-lime-100 flex justify-center items-center ${uploading ? "bg-lime-100 scale-[0.95]" : "bg-lime-50"} transition`}
>
    <label class={`flex flex-col flex-1 justify-center items-center absolute inset-0 hover:cursor-pointer ${uploading ? "text-lime-800" : "text-lime-600"}`} ondrop={onDrop}>
    <div class="p-2 text-center">{uploading ? `Uploading ${uploadName} (${uploadCurrent} of ${uploadingCount})...` : "Drop audio files or click here to upload."}</div>
        <input class="hidden" type="file" accept="audio/*" />
    </label>
</div>
