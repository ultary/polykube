import grpc

from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2

from api.grpc.v1 import system_pb2
from api.grpc.v1 import system_pb2_grpc


class Client(object):

    def __init__(self):
        host = '127.0.0.1'
        port = 50051
        channel = grpc.insecure_channel(f'{host}: {port}')
        self.stub = system_pb2_grpc.SystemServiceStub(channel)

    def ping(self):
        empty = google_dot_protobuf_dot_empty__pb2.Empty()
        pong = self.stub.Ping(empty)
        return pong


client = Client()
