
from typing import Union
from pathlib import Path
from asyncio import sleep

import ffmpeg
# import wave

def convert_to_mp3(path: Path, output: Path) -> Union[Path, None]:
	"""Uses ffmpeg to convert an input file to an MP3 with the libmp3lame codec.
	
		Args:
			path (Path): The path of an input file.
			output (Path): The desired output destination of the converted file.
		Returns:
			output (Path): The desired output destination of the converted file.
			None: Returns None if ffmpeg fails to convert the file.
 	"""
 
	try:
		(
			ffmpeg.input(path)
			.output(str(output), acodec="libmp3lame")
			.run()
		)
		return output
	except ffmpeg.Error as e:
		print("Didn't convert:",e)
		return None

def audio_stream(path: Path):
	"""A generator function that yields chunks of mp3 data.

	Args:
		path (Path): The path of an audio file on the server.

	Yields:
		Bytes: A string of at most 4096 bytes.
	"""

	# We need files in mp3 format to work with web for the moment.
	mp3 = convert_to_mp3(path, Path("test.mp3"))
	if mp3:
		with open(mp3, "rb") as f:
			data = f.read(4096)
			while data:
				yield data
				data = f.read(4096)

# def audio_stream(path: Path):
# 	with wave.open(str(path), "rb") as f:
# 		chunk = 4096
# 		data = f.readframes(chunk)
# 		while data:
# 			yield data
# 			data = f.readframes(chunk)
