from http import HTTPStatus
from typing import Callable

from kubernetes.client import (
    V1Namespace,
    V1Secret,
    V1Service,
    V1StatefulSet,
)
from kubernetes.client.exceptions import ApiException
from k8s import apps, core


ALREADY_EXISTS = HTTPStatus.CONFLICT


def apply_k8s_namespace(m: dict[str, any]) -> None:
    apply_k8s_object(core.create_namespace,
                     core.patch_namespace,
                     V1Namespace,
                     m)


def apply_k8s_secret(namespace: str, m: dict[str, any]) -> None:
    patch = None if m['immutable'] == True else core.patch_namespaced_secret
    apply_k8s_namespaced_object(namespace,
                                core.create_namespaced_secret,
                                patch,
                                V1Secret,
                                m)


def apply_k8s_stateful_set(namespace: str, m: dict[str, any]) -> None:
    apply_k8s_namespaced_object(namespace,
                                apps.create_namespaced_stateful_set,
                                apps.patch_namespaced_stateful_set,
                                V1StatefulSet,
                                m)


def apply_k8s_service(namespace: str, m: dict[str, any]) -> None:
    apply_k8s_namespaced_object(namespace,
                                core.create_namespaced_service,
                                core.patch_namespaced_service,
                                V1Service,
                                m)


def apply_k8s_namespaced_object(namespace: str, create: Callable, patch: Callable | None, body_class, m: dict[str, any]) -> None:
    name = m['metadata']['name']
    m['api_version'] = m.pop('apiVersion')
    body = body_class(**m)
    try:
        create(namespace=namespace, body=body)
    except ApiException as e:
        if e.status != ALREADY_EXISTS:
            raise
        if patch is not None:
            patch(name=name, namespace=namespace, body=body)


def apply_k8s_object(create: Callable, patch: Callable | None, body_class, m: dict[str, any]) -> None:
    m['api_version'] = m.pop('apiVersion')
    body = body_class(**m)
    try:
        create(body)
    except ApiException as e:
        if e.status != ALREADY_EXISTS:
            raise
        if patch is not None:
            name = m['metadata']['name']
            patch(name, body)
