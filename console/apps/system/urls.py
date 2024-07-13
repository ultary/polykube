from django.urls import include, path
from .htmx import OpenTelemetryView


app_name = "apps"
urlpatterns = [
    path('otlp/', OpenTelemetryView.as_view(), name="otlp"),
]
