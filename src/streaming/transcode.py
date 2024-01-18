from pathlib import Path
from typing import Union
from typing import List
import logging

import ffmpeg

logger = logging.getLogger(__name__)

def transcode_to_mp3(path: Path) -> Union[Path, None]:
    if not path.exists:
        logger.error(f"{path} is not a valid path to an audio file.")
        return None
    
    try:
        logger.info(f"Converting '{path}' to MP3.")
        output_path = Path(f"{path.stem}.mp3")
        output_path.mkdir()
        logger.info(f"Created output directory '{output_path}'")
        (
            ffmpeg.input(path)
                .output(
                    f"{output_path}/{path.stem}.mp3",
                    acodec="libmp3lame",
                )
                .run()
        )
        return output_path
    except ffmpeg.Error as e:
        logger.error(f"{e.stderr}")
    return None

def transcode_to_dash(path: Path) -> Union[List[Path], None]:
    if not path.exists:
        logger.error(f"{path} is not a valid path to an audio file.")
        return None

    try:
        logger.info(f"Converting '{path}' for DASH.")
        output_path = Path(f"{path.stem}-converted")
        if output_path.exists:
            return output_path
        output_path.mkdir()
        logger.info(f"Created output directory '{output_path}'")
        (
            ffmpeg.input(path)
                .output(
                    f"{output_path}/{path.stem}.mpd", 
                    f="dash",
                    init_seg_name=f"{path.stem}-$Number%05d$.mp3",
                    media_seg_name=f"{path.stem}-$Number%05d$.mp3",
            )
            .run()
        )
        return output_path
    except ffmpeg.Error as e:
        logger.error(f"{e.stderr}")
        return None