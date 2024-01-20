from pathlib import Path
import filecmp
import os
from os.path import getsize

from streaming.transcode import transcode_to_mp3

import pytest


@pytest.mark.unit
def test_transcode_to_mp3():
    # Could be a fixture later
    cwd = Path(os.getcwd())
    os.chdir("src/streaming/tests")

    test_path = Path("src/streaming/tests")
    sample_path = Path("sample_test_audio.mp3")
    wav_path = Path("test_audio.wav")

    result_path = transcode_to_mp3(wav_path)
    assert result_path == Path("test_audio/test_audio.mp3")
    assert getsize(sample_path) == getsize(result_path)
    assert filecmp.cmp(sample_path, result_path)
    os.remove(result_path)
    os.removedirs(result_path.parent)
    os.chdir(cwd)
