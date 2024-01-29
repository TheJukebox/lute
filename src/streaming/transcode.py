from pathlib import Path
from typing import Union
from typing import List
import logging

import ffmpeg

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)
ch = logging.StreamHandler()
ch.setLevel(logging.DEBUG)
formatter = logging.Formatter("%(asctime)s - %(name)s - %(levelname)s - %(message)s")
ch.setFormatter(formatter)
logger.addHandler(ch)


def transcode_to_mp3(path: Path) -> Union[Path, None]:
    if not path.exists():
        logger.error(f"{path} is not a valid path to an audio file.")
        return None

    if Path(path.stem).exists():
        return Path(f"{path.stem}/{path.stem}.mp3")
    try:
        logger.info(f"Converting '{path}' to MP3.")
        output_directory = Path(f"{path.stem}")
        if not output_directory.exists():
            output_directory.mkdir()
        output_path = Path(f"{output_directory}/{path.stem}.mp3")
        logger.info(f"Created output directory '{output_path}'")
        (
            ffmpeg.input(path)
            .output(
                str(output_path),
                acodec="libmp3lame",
            )
            .run(overwrite_output=True)
        )
        return output_path
    except ffmpeg.Error as e:
        logger.error(f"{e.stderr}")
    return None


def transcode_to_hls(path: Path) -> Union[Path, None]:
    if not path.exists():
        logger.error(f"{path} is not a valid path to an audio file.")
        return None

    try:
        logger.info(f"Transcoding '{path}' to HLS format.")
        output_path = Path(f"{path.stem}-hls")
        if not output_path.exists():
            output_path.mkdir()
            logger.info(f"Created directory '{output_path}'.")
        (
            ffmpeg.input(path)
            .output(
                f"{output_path}/{path.stem}.m3u8",
                f="hls",
                map="0:a",
                acodec="aac",
                hls_time="10",
                hls_list_size="0",
                hls_segment_filename=f"{output_path}/{path.stem}-%03d.ts",
            )
            .run(overwrite_output=True)
        )
        return output_path
    except ffmpeg.Error as e:
        logger.error(f"{e.stderr}")
        return None


def transcode_to_dash(path: Path) -> Union[Path, None]:
    if not path.exists():
        logger.error(f"{path} is not a valid path to an audio file.")
        return None

    try:
        logger.info(f"Converting '{path}' for DASH.")
        output_path = Path(f"{path.stem}-converted")
        if not output_path.exists():
            output_path.mkdir()
        logger.info(f"Created output directory '{output_path}'")
        (
            ffmpeg.input(path)
            .output(
                f"{output_path}/{path.stem}.mpd",
                f="dash",
                init_seg_name=f"{path.stem}-init.mp4",
                media_seg_name=f"{path.stem}-$Number%01d$.mp4",
            )
            .run(overwrite_output=True)
        )
        return output_path
    except ffmpeg.Error as e:
        logger.error(f"{e.stderr}")
        return None
