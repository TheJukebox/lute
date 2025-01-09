type Frame = {
    frame: Uint8Array;
    seq: number;
}

type FrameMessage = {
    type: string | undefined;
    frame: Frame | undefined;
}


let frameQueue: Array<Frame> = [];

self.onmessage = async (event: MessageEvent<FrameMessage>) => {
    const { type, frame } = event.data;
    let msg: FrameMessage = {type: undefined, frame: undefined};
    switch (type) {
        case 'queue':
            if (frame) {
                frameQueue.push(frame);
                frameQueue.sort((a, b) => a.seq - b.seq);
                msg.type = 'queue_ok';
            }
            break;
        case 'dequeue':
            let next: Frame | undefined = frameQueue.shift();
            if (next) {
                msg.frame = next;
                msg.type = 'dequeue_ok';
            } else {
                msg.type = 'dequeue_fail';
            }
            break;
    }
    postMessage(msg);
};


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
