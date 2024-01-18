from typing import Union
from pathlib import Path
import struct
import asyncio
<<<<<<< HEAD

from audio.stream import audio_stream
import asyncio

from fastapi import FastAPI, WebSocket
from fastapi.responses import StreamingResponse
from fastapi.staticfiles import StaticFiles

app = FastAPI()
app.mount("/static", StaticFiles(directory="lute/api/static"), name="static")

=======
import logging

import streaming.transcode as tcode
import asyncio

from fastapi import FastAPI
from fastapi import Request
from fastapi import WebSocket
from fastapi import UploadFile
from fastapi.responses import StreamingResponse
from fastapi.responses import FileResponse
from fastapi.staticfiles import StaticFiles

app = FastAPI()
app.mount("/static", StaticFiles(directory="src/api/static"), name="static")

logger = logging.getLogger(__name__)

def iterfile(path: Path):
    """A generator for file-like objects.

    Args:
        path (Path): A path for a file-like object.

    Yields:
        bytes: Yields bytes from the file-like object.
    """
    
    logger.info(f"Opening {path} for streaming.")
    with open(path, "rb") as f:
        yield from f

@app.get("/mpd/{path}")
async def mpd(path: str, request: Request):
    """Endpoint for the client to receive an MPD manifest for a given file.

    Args:
        path (str): The path to the file the user wants to stream.
        request (Request): The request object created by FastAPI.
    """
    
    client_host = request.client.host
    logger.info(f"{client_host} has requested the MPD file for {path}")
     # get the MPD file and stream it to the client
    mpd_path = Path(f"{output_path}/{path.stem}.mpd")
    with open(mpd_path, "rb") as f:
        logger.info(f"Sending MPD '{mpd_path}' to client {client_host}")
        FileResponse(content=f.read().decode("utf-8"), media_type="application/dash+xml")
    
@app.post("/upload/")
async def upload_file(file: UploadFile, request: Request):
    client_host = request.client.host
    logger.info(f"{client_host} is uploading {file.filename}")
    
    with open(file.filename, "wb") as output:
        contents = await file.read()
        output.write(contents)
    
    mp3_output_path = tcode.transcode_to_mp3(Path(file.filename))
    dash_output_path = tcode.transcode_to_dash(mp3_output_path)
    
    return {"filename": file.filename, "file_size": file.size}

@app.get("/stream/{path}")
async def stream(path: str, request: Request):   
    client_host = request.client.host
    logger.info(f"{client_host} is requesting a stream...")
	
    # Transcode to DASH
    path = Path(path)
    output_path = tcode.transcode_to_dash(path)
    mpd_path = Path("{output_path}/{path.stem}.mpd")
    
    segments = [segment for segment in output_path.iterdir() if segment != mpd_path]
    segments.sort()
    for segment in segments:
        with open(segment, "rb") as f:
            logger.info(f"Sending segment '{segment}' to client {client_host}")
            StreamingResponse(content=iterfile(segment), media_type="audio/mpeg")
        logger.info(f"Finished streaming segment '{segment}' to client {client_host}")
        
>>>>>>> f12d44c6ebae3c47afbdca8334e91f9b7e335661
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
