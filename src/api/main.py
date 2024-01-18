from typing import Union
from pathlib import Path
import struct
import asyncio

from audio.stream import audio_stream
import asyncio

from fastapi import FastAPI, WebSocket
from fastapi.responses import StreamingResponse
from fastapi.staticfiles import StaticFiles

app = FastAPI()
app.mount("/static", StaticFiles(directory="lute/api/static"), name="static")

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
