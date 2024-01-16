from pathlib import Path
from typing import Union
from typing import List
import logging

import ffmpeg

logger = logging.getLogger(__name__)

def encode_audio(path: Path) -> Union[List[Path], None]:
    if not path.exists:
        logger.error(f"{path} is not a valid path to an audio file.")
        return None

    try:
        logger.info(f"Converting '{path}' for DASH.")
        output_dir = Path(f"{path.stem}-converted")
        output_dir.mkdir()
        logger.info(f"Created output directory '{output_dir}'")
        (
        ffmpeg.input(path)
            .output(
                f"{output_dir}/{path.stem}.mpd", 
                f="dash",
                init_seg_name=f"{path.stem}-$Number%05d$.mp3",
                media_seg_name=f"{path.stem}-$Number%05d$.mp3",
            )
            .run()
        )
        return output_dir
    except ffmpeg.Error as e:
        logger.error(f"{e.stderr}")
        return None