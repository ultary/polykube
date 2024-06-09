from django.db import models


class Cluster(models.Model):
    name = models.CharField(max_length=255)
    network = models.CharField(max_length=255)
    config = models.JSONField(default=dict)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)
    deleted_at = models.DateTimeField(null=True)

    class Meta:
        db_table = u'k8s_cluster'
        unique_together = (('name', 'network'),)
