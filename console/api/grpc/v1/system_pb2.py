# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: api/grpc/v1/system.proto
# Protobuf Python Version: 5.29.1
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    1,
    '',
    'api/grpc/v1/system.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\x18\x61pi/grpc/v1/system.proto\x12\tdokevy.v1\x1a\x1bgoogle/protobuf/empty.proto\"\x14\n\x04Pong\x12\x0c\n\x04pong\x18\x01 \x01(\t\"%\n#EnableOpenTelemetryCollectorRequest\"&\n$EnableOpenTelemetryCollectorResponse\"&\n$DisableOpenTelemetryCollectorRequest\"\'\n%DisableOpenTelemetryCollectorResponse\"%\n#UpdateOpenTelemetryCollectorRequest\"&\n$UpdateOpenTelemetryCollectorResponse2\xca\x03\n\x06System\x12\x31\n\x04Ping\x12\x16.google.protobuf.Empty\x1a\x0f.dokevy.v1.Pong\"\x00\x12\x81\x01\n\x1c\x45nableOpenTelemetryCollector\x12..dokevy.v1.EnableOpenTelemetryCollectorRequest\x1a/.dokevy.v1.EnableOpenTelemetryCollectorResponse\"\x00\x12\x84\x01\n\x1d\x44isableOpenTelemetryCollector\x12/.dokevy.v1.DisableOpenTelemetryCollectorRequest\x1a\x30.dokevy.v1.DisableOpenTelemetryCollectorResponse\"\x00\x12\x81\x01\n\x1cUpdateOpenTelemetryCollector\x12..dokevy.v1.UpdateOpenTelemetryCollectorRequest\x1a/.dokevy.v1.UpdateOpenTelemetryCollectorResponse\"\x00\x42\x66\n\x19\x63o.ultary.kluster.grpc.v1B\x0bSystemProtoP\x01Z.github.com/ultary/polykube/kluster/api/grpc/v1\xa2\x02\tKlusterV1b\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'api.grpc.v1.system_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'\n\031co.ultary.kluster.grpc.v1B\013SystemProtoP\001Z.github.com/ultary/polykube/kluster/api/grpc/v1\242\002\tKlusterV1'
  _globals['_PONG']._serialized_start=68
  _globals['_PONG']._serialized_end=88
  _globals['_ENABLEOPENTELEMETRYCOLLECTORREQUEST']._serialized_start=90
  _globals['_ENABLEOPENTELEMETRYCOLLECTORREQUEST']._serialized_end=127
  _globals['_ENABLEOPENTELEMETRYCOLLECTORRESPONSE']._serialized_start=129
  _globals['_ENABLEOPENTELEMETRYCOLLECTORRESPONSE']._serialized_end=167
  _globals['_DISABLEOPENTELEMETRYCOLLECTORREQUEST']._serialized_start=169
  _globals['_DISABLEOPENTELEMETRYCOLLECTORREQUEST']._serialized_end=207
  _globals['_DISABLEOPENTELEMETRYCOLLECTORRESPONSE']._serialized_start=209
  _globals['_DISABLEOPENTELEMETRYCOLLECTORRESPONSE']._serialized_end=248
  _globals['_UPDATEOPENTELEMETRYCOLLECTORREQUEST']._serialized_start=250
  _globals['_UPDATEOPENTELEMETRYCOLLECTORREQUEST']._serialized_end=287
  _globals['_UPDATEOPENTELEMETRYCOLLECTORRESPONSE']._serialized_start=289
  _globals['_UPDATEOPENTELEMETRYCOLLECTORRESPONSE']._serialized_end=327
  _globals['_SYSTEM']._serialized_start=330
  _globals['_SYSTEM']._serialized_end=788
# @@protoc_insertion_point(module_scope)
