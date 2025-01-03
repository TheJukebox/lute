# Lute Developer Utilities
This directory contains Python utilities to assist in development tasks.

## Installing Utility Dependencies
You should use a [Python virtual environment](https://docs.python.org/3/library/venv.html) to keep your Python path clean:

```bash
$ python -m venv lute-utils
$ source lute-utils/bin/activate
(lute-utils) $ pip install -r requirements.txt
```

## Running Utilities

### pb_build
The `pb_build` utility can be used to compile [protobuffer](https://protobuf.dev/) files for Lute.

```bash
# Builds all protobuffers found in the current directory tree.
$ utils/pb_build.py --all .

# Build a specific protobuffer
$ utils/pb_build.py api/proto/stream.proto
```

### upload
The `upload` utility can be used to upload files to the server via RPC.

```bash
# Uploads a file to the server
$ utils/upload/upload.py path/to/file.mp3
```

## (WIP) Building New Utilities
Python utilities for Lute are built using [click](https://click.palletsprojects.com/en/stable/), to provide
consistent and featured CLIs out of the gate. New utilities should follow the same pattern where possible. If the
utility requires compiling proto files, they should be packed into a flat directory (see the [upload](./upload) utility).

When you build a new utility, ensure that its requirements are added to the [requirements](requirements.txt) file:

```bash
pip freeze > utils/requirements.txt
```

Utilities should be friendly for use in CI/CD, using appropriate exit codes and returning enough input to ease debugging.