<script lang="ts">
    import { onMount } from 'svelte';

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
                const name: string = file.name.split(".")[0];
                uploadName = name;
                const uriName: string = encodeURIComponent(name);
                const contentType: string = file.type;
                const response: Response = await fetch(
                    "http://172.31.204.147:7001/upload",
                    {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        body: JSON.stringify({
                            Name: name,
                            UriName: uriName,
                            ContentType: contentType,
                        }),
                    },
                );
                const data: Object = await response.json();
                const formData = new FormData();
                for (const [k, v] of Object.entries(data.fields)) {
                    formData.append(k, v);
                }
                formData.append("file", file);
                const uploadResponse: Response = await fetch(
                    data.url, 
                    {
                        method: "POST",
                        body: formData,
                    },
                );
                console.debug(uploadResponse);
                uploadCurrent++;
            }
            uploading = false;
            uploadCurrent = 0;
        }
    };
</script>

<div 
    class={`relative p-4 min-h-30 min-w-90 max-h-30 max-w-90 rounded-xl shadow hover:bg-lime-100 flex justify-center items-center ${uploading ? "bg-lime-100 scale-[0.95]" : "bg-lime-50"} transition`}
>
    <label class={`flex flex-col flex-1 justify-center items-center absolute inset-0 hover:cursor-pointer ${uploading ? "text-lime-800" : "text-lime-600"}`} ondrop={onDrop}>
    <div class="p-2 text-center">{uploading ? `Uploading ${uploadName} (${uploadCurrent} of ${uploadingCount})...` : "Drop audio files or click here to upload."}</div>
        <input class="hidden" type="file" accept="audio/*" />
    </label>
</div>
