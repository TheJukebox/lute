from typing import Union
from pathlib import Path
import struct
import asyncio
import logging

import streaming.transcode as tcode
import asyncio

from fastapi import FastAPI
from fastapi import Request
from fastapi import WebSocket
from fastapi.responses import StreamingResponse
from fastapi.staticfiles import StaticFiles

app = FastAPI()
app.mount("/static", StaticFiles(directory="src/api/static"), name="static")

logger = logging.getLogger(__name__)

@app.get("/stream/{path}")
async def stream(path: str, request: Request):
    client_host = request.client.host
    logger.info(f"{client_host} is requesting a stream...")
	
    # Transcode to DASH
    path = Path(path)
    output_path = tcode.transcode_to_dash(path)
    
    # get the MPD file and initial segment
    mpd_path = Path(f"{output_path}/{path.stem}.mpd")
    with open(mpd_path, "rb") as f:
        logger.info(f"Sending MPD '{mpd_path}' to client {client_host}")
        StreamingResponse(content=f.read().decode("utf-8"), media_type="application/dash+xml")
    logger.info(f"Finished streaming MPD '{mpd_path}' to client {client_host}")
    
    segments = [segment for segment in output_path.iterdir() if segment != mpd_path]
    segments.sort()
    for segment in segments:
        with open(segment, "rb") as f:
            logger.info(f"Sending segment '{segment}' to client {client_host}")
            StreamingResponse(content=f.read(), media_type="audio/mpeg")
        logger.info(f"Finished streaming segment '{segment}' to client {client_host}")
        
@app.websocket("/audio/{path}")
async def audio(websocket: WebSocket, path: str):
	"""A websocket endpoint that streams audio on request to clients.

	Args:
		websocket (WebSocket): The websocket.
		path (str): A path to an audio file on the server.
	"""
 
	await websocket.accept()
	for data in audio_stream(Path(path)):
		print(f"Sending {len(data)} bytes of data...")
		await asyncio.sleep(0.1)
		await websocket.send_bytes(data)
	await websocket.close()
