# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: upload.proto
# Protobuf Python Version: 5.28.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    28,
    1,
    '',
    'upload.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x0cupload.proto\x12\x06upload\"B\n\rUploadRequest\x12\x11\n\tfile_name\x18\x01 \x01(\t\x12\x10\n\x08\x63hecksum\x18\x02 \x01(\t\x12\x0c\n\x04size\x18\x03 \x01(\x05\"!\n\x0eUploadResponse\x12\x0f\n\x07\x66ile_id\x18\x01 \x01(\t\"5\n\x05\x43hunk\x12\x0f\n\x07\x66ile_id\x18\x01 \x01(\t\x12\x0c\n\x04\x64\x61ta\x18\x02 \x01(\x0c\x12\r\n\x05\x66inal\x18\x03 \x01(\x08\"1\n\rChunkResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\x12\x0f\n\x07message\x18\x02 \x01(\t2\x7f\n\x06Upload\x12>\n\x0bStartUpload\x12\x15.upload.UploadRequest\x1a\x16.upload.UploadResponse\"\x00\x12\x35\n\x0bUploadChunk\x12\r.upload.Chunk\x1a\x15.upload.ChunkResponse\"\x00\x42\rZ\x0b/gen/uploadb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'upload_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\013/gen/upload'
  _globals['_UPLOADREQUEST']._serialized_start=24
  _globals['_UPLOADREQUEST']._serialized_end=90
  _globals['_UPLOADRESPONSE']._serialized_start=92
  _globals['_UPLOADRESPONSE']._serialized_end=125
  _globals['_CHUNK']._serialized_start=127
  _globals['_CHUNK']._serialized_end=180
  _globals['_CHUNKRESPONSE']._serialized_start=182
  _globals['_CHUNKRESPONSE']._serialized_end=231
  _globals['_UPLOAD']._serialized_start=233
  _globals['_UPLOAD']._serialized_end=360
# @@protoc_insertion_point(module_scope)
