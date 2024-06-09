import random
import string

from .kubernetes import (
    apply_k8s_namespaced_object,
    apply_k8s_object,
    apply_k8s_namespace,
    apply_k8s_secret,
    apply_k8s_service,
    apply_k8s_stateful_set,
)

from .postgres import (
    create_pg_database,
    create_pg_database_role,
)


def new_password(len: int = 16):
    characters = string.ascii_letters + string.digits + string.punctuation
    return ''.join(random.choice(characters) for i in range(len))
