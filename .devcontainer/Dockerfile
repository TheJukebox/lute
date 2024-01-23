FROM python:3.12-alpine3.19

# Update mirrors and install ffmpeg and git
RUN apk upgrade -U && \
    apk add ffmpeg && \
    apk add git

# Install python deps
COPY ./requirements.txt requirements.txt
RUN pip install -r requirements.txt

ENTRYPOINT [“/bin/bash”, “-c”]