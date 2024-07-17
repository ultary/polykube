from django.db import models


class ResourceStatus(models.Model):
    api_group = models.CharField(blank=True, default='', max_length=32)
    api_version = models.CharField(max_length=16)
    kind = models.CharField(max_length=128)
    name = models.CharField(max_length=253)
    namespace = models.CharField(max_length=63, blank=True, default='')
    requested = models.JSONField(max_length=2**20) # 1MiB
    status = models.JSONField(max_length=2**21)    # 2MiB
    resource_version = models.PositiveBigIntegerField()
    uid = models.UUIDField()
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        app_label = 'kluster'
        db_table = u'kluster_resources_status'
        unique_together = ['api_group', 'kind', 'name', 'namespace', 'uid']


class LatestRsourceKindVersion(models.Model):
    resource_version = models.PositiveBigIntegerField(primary_key=True)
    updated_at = models.DateTimeField(auto_now=True)

    class Meta:
        app_label = 'kluster'
        db_table = u'kluster_latest_event_resource_version'
        indexes = [
            models.Index(fields=['updated_at']),
        ]
