from django.http import HttpResponse
from django.shortcuts import render

from api.kluster import client


def index(request):

    resp = client.ping()
    print('Pong: ', resp.pong)

    context = {}
    tempfile_name = 'index.html'
    return render(request, tempfile_name, context)


def ready(request):
    return HttpResponse('OK')


def healthz(request):
    return HttpResponse('OK')
