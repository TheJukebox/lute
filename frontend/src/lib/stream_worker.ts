self.onmessage = async (event: MessageEvent<{data: ArrayBuffer; sampleRate: number}>) => {
    const { data, sampleRate } = event.data; 
    const context: AudioContext = new AudioContext({"sampleRate": sampleRate})
    self.postMessage({ success: true, audioBuffer: await context.decodeAudioData(data)});
}