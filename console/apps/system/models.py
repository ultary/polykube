from django.db import models

# Create your models here.


class Application(models.Model):

    class Apps(models.TextChoices):
        CILIUM = "cilium"
        OTEL_AGENT = "otel-agent"
        OTEL_COLLECTOR = "otel-collector"
        OTEL_GATEWAY = "otel-gateway"

    class Status(models.TextChoices):
        CREATE = "create"
        UPDATE = "update"
        DELETE = "delete"
        DONE = "done"

    name = models.CharField(max_length=64, choices=Apps)
    status = models.CharField(max_length=8, choices=Status, default=Status.CREATE)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    deleted_at = models.DateTimeField(default=None, null=True)

    class Meta:
        app_label = 'kluster'
        db_table = u'kube_system_applications'
        unique_together = ['name']
