from django.http import HttpResponse
from django.shortcuts import render
from django.views import View

from api.grpc.client import client


class OpenTelemetryView(View):

    def post(self, request, *args, **kwargs):
        print('enable OpenTelemetry')
        client.enable_opentelemetry()
        return HttpResponse('OK')

    def put(self, request, *args, **kwargs):
        print('update OpenTelemetry')
        client.disable_opentelemetry()
        return HttpResponse('OK')

    def delete(self, request, *args, **kwargs):
        print('disable OpenTelemetry')
        client.update_opentelemetry()
        return HttpResponse('OK')
