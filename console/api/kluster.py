import grpc

from django.conf import settings
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2

from api.grpc.v1 import system_pb2_grpc


class Client(object):

    def __init__(self):
        host = settings.MK_GRPC_HOST
        port = settings.MK_GRPC_PORT
        channel = grpc.insecure_channel(f'{host}: {port}')
        #channel = grpc.insecure_channel('unix:///tmp/kluster.sock')
        self.stub = system_pb2_grpc.SystemStub(channel)

    def ping(self):
        empty = google_dot_protobuf_dot_empty__pb2.Empty()
        pong = self.stub.Ping(empty)
        return pong


client = Client()
