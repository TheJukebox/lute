from typing import Union
from pathlib import Path
import struct
import asyncio
import logging
import json
from uuid import uuid4

import streaming.transcode as transcode
import streaming.operations as operations

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


@app.get("/mpd/{uid}")
async def mpd(uid: str, request: Request):
    """Endpoint for the client to receive an MPD manifest for some file in the library.

    Args:
        uid (str): The UID of the requested library file.
        request (Request): The request object created by FastAPI.
    """
    client_host = request.client.host
    logger.info(f"{client_host} has requested the MPD file for {uid}")

    library_path = Path("library.json")
    if not library_path.exists():
        return {"error": "There is no library, so no MPD. Oops!"}

    library = {}
    # Open the library file
    with open(library_path, "r") as f:
        library = json.loads(f.read())

    # get the MPD file and stream it to the client
    content = library[uid]
    if not content:
        return {"error": "No library entry with uid '{uid}'"}

    mpd_path = Path(library[uid]["mpd_path"])
    logger.info(f"Sending MPD '{mpd_path}' to client {client_host}")
    return FileResponse(path=mpd_path, media_type="application/dash+xml")


@app.post("/upload/")
async def upload_file(file: UploadFile, request: Request):
    client_host = request.client.host
    logger.info(f"{client_host} is uploading {file.filename}")

    # Asynchronously read the file from the client
    path = Path(file.filename)
    with open(file.filename, "wb") as output:
        contents = await file.read()
        output.write(contents)

    # Transcode to MP3 and then to DASH
    logger.info(f"Transcoding {path} to MP3...")
    mp3_output_path = transcode.transcode_to_mp3(path)
    logger.info(f"Created {mp3_output_path}")
    logger.info(f"Transcoding {mp3_output_path} to DASH")
    dash_output_path = transcode.transcode_to_dash(mp3_output_path)
    logger.info(f"Created {dash_output_path}")

    # Quick and dirty library file, abstracting a database
    library_path = Path("library.json")
    if library_path.exists():
        with open(library_path, "r") as f:
            library = json.loads(f.read())
    else:
        library = {}

    uid = str(uuid4())
    library[uid] = {
        "title": path.stem,
        "path": str(dash_output_path),
        "mpd_path": str(Path(f"{dash_output_path}/{path.stem}.mpd")),
    }

    library_json = json.dumps(library)
    with open("library.json", "w") as library_file:
        library_file.write(library_json)

    # Return some info to the client
    return {
        "uid": uid,
        "filename": file.filename,
        "file_size": file.size,
        "title": str(library[uid]["title"]),
        "mpd_path": str(library[uid]["mpd_path"]),
        "output_path": str(library[uid]["path"]),
    }


@app.get("/stream/{uid}")
def stream(uid: str, request: Request):
    client_host = request.client.host
    logger.info(f"{client_host} is requesting a stream for {uid}")

    library_path = Path("library.json")
    if not library_path.exists():
        logger.error("There is no library file")
        yield {"error": "There is no library file"}

    logger.info("Opening library file...")
    library = {}
    with open(library_path, "r") as library_file:
        library = json.loads(library_file.read())

    logger.info("Sending segments to client")
    path = Path(library[uid]["path"])
    mpd_path = Path(library[uid]["mpd_path"])
    segments = [segment for segment in path.iterdir() if segment != mpd_path]
    segments.sort()
    for segment in segments:
        yield FileResponse(path=segment, media_type="audio/mpeg")
