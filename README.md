# Lute

Lute is an open-source, self-hosted audio streaming platform, that allows users to listen to and interact with the same audio stream together.

## Installation

To get started with Lute, create a Python virtual environment and install the requirements:

```bash
# Navigate to the root of the project
$ cd /path/to/lute
$ python -m venv lute-venv
$ source lute-venv/bin/activate
(lute-venv) $ pip install -r requirements.txt
```

## Running the project

To run the uvicorn server and access the SwaggerUI for FastAPI:

```py
# Navigate to the root of the project
$ cd /path/to/lute
# Activate your lute venv
$ source lute-venv/bin/activate
(lute-venv) $ uvicorn --app-dir src api.main:app --reload
```

You should now be able to see the SwaggerUI at `http://127.0.0.1:8000/docs`