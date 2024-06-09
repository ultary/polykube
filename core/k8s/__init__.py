import yaml

from django.conf import settings
from kubernetes import config, client, utils
from kubernetes.client import (
    AppsV1Api,
    CoreV1Api,
)
from kubernetes.client.api_client import ApiClient


if settings.MK_IN_CLUSTER:
    config.load_incluster_config()
else:
    f = settings.MK_KUBECONFIG_FILE
    c = settings.MK_KUBECONFIG_CONTEXT
    config.load_kube_config(config_file=f, context=c)

client = ApiClient()
apps = AppsV1Api(client)
core = CoreV1Api(client)


def create_from_yaml(manifest: str):
    obj = yaml.safe_load(manifest)
    return utils.create_from_yaml(client.api_client.ApiClient(), yaml_objects=[obj])
