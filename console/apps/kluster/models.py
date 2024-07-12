from django.db import models


class Manifests(models.Model):
    api_group = models.CharField(blank=True, default='', max_length=32)
    api_version = models.CharField(max_length=16)
    kind = models.CharField(max_length=128)
    name = models.CharField(max_length=253)
    namespace = models.CharField(max_length=63, blank=True, default='')
    raw = models.TextField(max_length=2**20) # 1MiB
    created_at = models.DateTimeField(auto_now_add=True)
    requested_at = models.DateTimeField(auto_now=True)
    committed_at = models.DateTimeField(null=True)

    class Meta:
        app_label = 'kluster'
        db_table = u'kluster_manifests'
        unique_together = ('api_group', 'api_version', 'kind', 'name', 'namespace')


class LatestResourceKindChanged(models.Model):
    api_version = models.CharField(max_length=253)
    kind = models.CharField(max_length=63)
    name = models.CharField(max_length=253)
    namespace = models.CharField(max_length=63)

    class Meta:
        app_label = 'kluster'
        db_table = u'kluster_last_resource_changed'
        unique_together = ('api_version', 'kind', 'name', 'namespace')
