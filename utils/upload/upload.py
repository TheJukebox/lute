#!/usr/bin/env python
from pathlib import Path
import sys
import click

import grpc
import grpc._channel

import upload_pb2 as uploadpb
import upload_pb2_grpc as upload_grpc

SERVER_ADDRESS = "localhost:50051"

def upload_file(path: Path):
    channel = grpc.insecure_channel(
        SERVER_ADDRESS
    )
    stub = upload_grpc.UploadStub(channel)
    
    try:
        with open(path, "rb") as f:
            data = f.read()
            size = len(data)
    except FileNotFoundError:
       click.echo(click.style(f"File not found: {path}", fg="red")) 
       exit(1)
        
    request = uploadpb.UploadRequest(
        file_name=path.name,
        checksum="test",
        size=size,
    )

    click.echo(f"Requesting upload of '{path}'...")
    try:
        response = stub.StartUpload(request)
        if not response.file_id:
            click.echo(click.style(f"Request failed!", fg="red"))
            exit(1) 

        file_id = response.file_id
        click.echo(click.style(f"Got file ID: {file_id}", fg="green"))
    except grpc._channel._InactiveRpcError as e:
        click.echo(click.style(f"Request was rejected by server with message:", fg="red"))
        click.echo(click.style(f"\t{e.details()}", fg="yellow"))
        exit(1)
    
    chunk_size = 1024 * 1024 * 3
    click.echo(f"Uploading with a chunk size of {chunk_size / 1024 / 1024}MB")
    
    i = 0
    with click.progressbar(range(size), length=size, label="Uploading") as bar:
        while i < size:
            chunk = uploadpb.Chunk(
                file_id=file_id,
                data=data[i:i + chunk_size],
                final=(size - i) <= chunk_size,
            )
            response = stub.UploadChunk(chunk)
            i += len(chunk.data)
            bar.update(len(chunk.data))
    click.echo(click.style(f"Upload complete!", fg="green"))

@click.command()
@click.argument("path", required=True)
def upload(path: Path):
    """A click command for uploading a file to the Lute backend"""
    
    upload_file(Path(path))

if __name__ == "__main__":
    upload()