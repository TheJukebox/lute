from pathlib import Path
import sys
import click

import grpc

import upload_pb2 as upload
import upload_pb2_grpc as upload_grpc

SERVER_ADDRESS = "localhost:50051"

def upload_file(path: Path):
    channel = grpc.insecure_channel(
        SERVER_ADDRESS
    )
    stub = upload_grpc.UploadStub(channel)
    
    with open(path, "rb") as f:
       data = f.read()
       size = len(data)
        
    print(f"Creating new request to upload '{path}'...")
    request = upload.UploadRequest(
        file_name=path.name,
        checksum="test",
        size=size,
    )
    response = stub.StartUpload(request)
    if response.file_id:
        print("Success!")
        file_id = response.file_id
    
    chunk_size = 1024 * 1024 * 3
    i = 0
    with click.progressbar(range(size), length=size, label="Uploading") as bar:
        while i < size:
            chunk = upload.Chunk(
                file_id=file_id,
                data=data[i:i + chunk_size],
                final=(size - i) <= chunk_size,
            )
            response = stub.UploadChunk(chunk)
            i += len(chunk.data)
            bar.update(len(chunk.data))

upload_file(Path("../../Breezeblocks.mp3"))
