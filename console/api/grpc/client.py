import grpc

from django.conf import settings
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2

from api.grpc.v1.system_pb2 import (
    EnableOpenTelemetryCollectorRequest,
    DisableOpenTelemetryCollectorRequest,
    UpdateOpenTelemetryCollectorRequest,
)
from api.grpc.v1 import (
    system_pb2_grpc,
)


class Client(object):

    def __init__(self):
        host = settings.MK_GRPC_HOST
        port = settings.MK_GRPC_PORT
        channel = grpc.insecure_channel(f'{host}: {port}')
        #channel = grpc.insecure_channel('unix:///tmp/kluster.sock')
        self.system = system_pb2_grpc.SystemStub(channel)

    def ping(self):
        empty = google_dot_protobuf_dot_empty__pb2.Empty()
        pong = self.system.Ping(empty)
        return pong

    def enable_opentelemetry(self):
        req = EnableOpenTelemetryCollectorRequest()
        res = self.system.EnableOpenTelemetryCollector(req)
        return res

    def disable_opentelemetry(self):
        req = DisableOpenTelemetryCollectorRequest()
        res = self.system.DisableOpenTelemetryCollector(req)
        return res

    def update_opentelemetry(self):
        req = UpdateOpenTelemetryCollectorRequest()
        res = self.system.UpdateOpenTelemetryCollector(req)
        return res


client = Client()
