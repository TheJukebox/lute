<script lang="ts">
    import { onMount } from 'svelte';

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
            console.log(e);
            const name: string = e.dataTransfer.files.item(0).name.split(".")[0];
            const uriName: string = encodeURIComponent(name);
            const contentType: string = e.dataTransfer.files.item(0).type;
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
            formData.append("file", e.dataTransfer.files.item(0));
            const uploadResponse: Response = await fetch(
                data.url, 
                {
                    method: "POST",
                    body: formData,
                },
            );
            console.log(uploadResponse);
        }
    };
</script>

<div 
    class="relative p-4 min-h-30 min-w-90 max-h-30 max-w-90 rounded-xl shadow bg-lime-50 hover:bg-lime-100 flex justify-center items-center"
>
    <label class="flex flex-col flex-1 justify-center items-center absolute inset-0 hover:cursor-pointer target:text-orange-500" ondrop={onDrop}>
        <div class="p-2 text-center">Drop audio files or click here to upload.</div>
        <input class="hidden" type="file" accept="audio/*" />
    </label>
</div>
