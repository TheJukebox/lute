import { fetchTracks } from '$lib/tracklist.svelte';

export interface Track {
    name: string;
    uriName: string;
    contentType: string;
    artist: string;
    date: string;
    release: string;
    album: string;
    number: number;
    disk: number;
}

export interface PresignedUrl {
    url: string;
    fields: Record<string, string>;
}

export async function uploadTrack(
    track: Track,
    file: File,
) {
    const response: Response = await fetch(
        "http://localhost:7001/upload",
        {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(track), 
        },
    );
    if (!response.ok) {
        console.error(`Failed to fetch presigned URL for ${file.name}`);
        return;
    }
    const data: any = await response.json();
    const formData = new FormData();
    for (const [k, v] of Object.entries(data.fields)) {
        formData.append(k, v);
    }
    formData.append("file", file);
    const uploadResponse: Response = await fetch(
        "http://localhost:9000/lute-audio",
        {
            method: "POST",
            body: formData,
        },
    );
    if (uploadResponse.ok) {
        console.log(`Successfully uploaded ${file.name}.`);
        fetchTracks();
    } else {
        console.error(`Failed to upload ${file.name}`);
    }
}
