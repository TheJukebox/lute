let frameQueue: Array<{frame: Uint8Array, seq: number}> = [];

type Frame = {
    frame: Uint8Array;
    seq: number;
}

let nextSeq: number = 1;

self.onmessage = async (event: MessageEvent<Frame>) => {
    await awaitPrevious(event.data.seq);
    frameQueue.push(event.data);
    if (frameQueue.length > 2) {
        let frames: Uint8Array = new Uint8Array(0);
        let next = frameQueue.shift();
        if (next) frames = concatArrays(frames, next.frame);
        let seq = next?.seq;
        while (frameQueue.length > 0) {
            let next: Uint8Array | undefined = frameQueue.shift()?.frame;
            if (next) frames = concatArrays(frames, next);
        }
        self.postMessage({frames: frames, seq: seq});
    }
};


async function awaitPrevious(seq: number): Promise<void> {
    while (seq !== nextSeq) {
        return new Promise(resolve => setTimeout(resolve, 10));
    }
    nextSeq = seq + 1;
}


/**
 * Concatenates two byte arrays and returns the result.
 * @function
 * @param x {Uint8Array}    The first array.
 * @param y {Uint8Array}    The second array, to be appended.
 * @returns {Uint8Array}    An array with x and y's data combined.
 */
function concatArrays(x: Uint8Array, y: Uint8Array): Uint8Array {
    const combinedData: Uint8Array = new Uint8Array(x.length + y.length);
    combinedData.set(x);
    combinedData.set(y, x.length);
    return combinedData;
}
