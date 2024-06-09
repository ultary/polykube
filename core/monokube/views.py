from django.http import HttpResponse
from django.shortcuts import render


def index(request):
    context = {}
    tempfile_name = 'index.html'
    return render(request, tempfile_name, context)


def ready(request):
    return HttpResponse('OK')


def healthz(request):
    return HttpResponse('OK')
