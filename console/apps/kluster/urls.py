from django.urls import include, path
from . import views


app_name = "kluster"
urlpatterns = [
    path('', views.index, name="home"),
]
