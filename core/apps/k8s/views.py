from django.shortcuts import render

# Create your views here.

def index(request):
    context = None
    tempfile_name = 'k8s/index.html'
    return render(request, tempfile_name, context)
