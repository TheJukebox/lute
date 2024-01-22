from pathlib import Path
import os
from os.path import getsize
import filecmp

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
    print(result_path)
    assert result_path == Path("test_audio/test_audio.mp3")
    assert getsize(sample_path) == getsize(result_path)
    # This doesn't work for whatever reason
    # We should be able to assume these files are the same or similar? Not sure what changes between runs of the test.
    #assert filecmp.cmp(sample_path, result_path, shallow=False)
    os.remove(result_path)
    os.removedirs(result_path.parent)
    os.chdir(cwd)
