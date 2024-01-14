const socket = new WebSocket("ws://127.0.0.1:8000/audio/1TORCHLEFT_3.wav");
socket.binaryType = "arraybuffer";

const ctx = new window.AudioContext();
let chunks = new ArrayBuffer();
let isPlaying = false;

socket.addEventListener('open', (event) => {
    console.log("opened");
});

socket.addEventListener('close', (event) => {
    console.log("closed");
});

socket.addEventListener('error', (event) => {
    console.log("There was an error:",event.data);
});

socket.addEventListener('message', async (event) => {
    if (event.data instanceof ArrayBuffer) {
        chunks = addChunks(chunks, event.data);
        console.log(chunks);
        if (chunks.byteLength > 409600 && !isPlaying) {
            playback();
            isPlaying = true;
         } else {
             console.log("Collecting more chunks...");
         }
    } else {
        console.log(event.data);
    }
});

function playback() {
    const curChunks = chunks.slice(0);
    ctx.decodeAudioData(
        curChunks,
        (buffer) => {
            src = ctx.createBufferSource();
            src.buffer = buffer;
            src.connect(ctx.destination);
            src.start();
        },
        (error) => {
            console.error("Error decoding audio data",error);
        }
    );
}

function addChunks(i, j) {
    const length = i.byteLength + j.byteLength;
    const k = new ArrayBuffer(length);

    const iv = new Uint8Array(i);
    const jv = new Uint8Array(j);
    const kv = new Uint8Array(k);

    kv.set(iv, 0);
    kv.set(jv, i.byteLength);

    return k;
}
