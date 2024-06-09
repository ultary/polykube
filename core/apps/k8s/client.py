from django.conf import settings
from kubernetes import client, config


if settings.MK_IN_CLUSTER:
    config.load_incluster_config()
else:
    f = settings.MK_KUBECONFIG_FILE
    c = settings.MK_KUBECONFIG_CONTEXT
    config.load_kube_config(config_file=f, context=c)

apps = client.AppsV1Api()
core = client.CoreV1Api()
