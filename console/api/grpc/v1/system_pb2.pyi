from google.protobuf import empty_pb2 as _empty_pb2
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional

DESCRIPTOR: _descriptor.FileDescriptor

class Pong(_message.Message):
    __slots__ = ("pong",)
    PONG_FIELD_NUMBER: _ClassVar[int]
    pong: str
    def __init__(self, pong: _Optional[str] = ...) -> None: ...

class EnableOpenTelemetryCollectorRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class EnableOpenTelemetryCollectorResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class DisableOpenTelemetryCollectorRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class DisableOpenTelemetryCollectorResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class UpdateOpenTelemetryCollectorRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class UpdateOpenTelemetryCollectorResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...
