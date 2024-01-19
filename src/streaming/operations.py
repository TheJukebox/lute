from pathlib import Path


def iterfile(path: Path):
    """A generator for file-like objects, which yields bytes as requested.

    Args:
        path (Path): A path for a file-like object.

    Yields:
        bytes: Yields bytes from the file-like object.
    """

    with open(path, "rb") as f:
        yield from f
