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