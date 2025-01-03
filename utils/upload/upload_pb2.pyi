from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class UploadRequest(_message.Message):
    __slots__ = ("file_name", "checksum", "size")
    FILE_NAME_FIELD_NUMBER: _ClassVar[int]
    CHECKSUM_FIELD_NUMBER: _ClassVar[int]
    SIZE_FIELD_NUMBER: _ClassVar[int]
    file_name: str
    checksum: str
    size: int
    def __init__(self, file_name: _Optional[str] = ..., checksum: _Optional[str] = ..., size: _Optional[int] = ...) -> None: ...

class UploadResponse(_message.Message):
    __slots__ = ("file_id",)
    FILE_ID_FIELD_NUMBER: _ClassVar[int]
    file_id: str
    def __init__(self, file_id: _Optional[str] = ...) -> None: ...

class Chunk(_message.Message):
    __slots__ = ("file_id", "data", "final")
    FILE_ID_FIELD_NUMBER: _ClassVar[int]
    DATA_FIELD_NUMBER: _ClassVar[int]
    FINAL_FIELD_NUMBER: _ClassVar[int]
    file_id: str
    data: bytes
    final: bool
    def __init__(self, file_id: _Optional[str] = ..., data: _Optional[bytes] = ..., final: bool = ...) -> None: ...

class ChunkResponse(_message.Message):
    __slots__ = ("success", "message")
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    success: bool
    message: str
    def __init__(self, success: bool = ..., message: _Optional[str] = ...) -> None: ...
