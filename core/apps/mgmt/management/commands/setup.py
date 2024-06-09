import base64
import yaml

from http import HTTPStatus
from pathlib import Path

from django.core.management.base import BaseCommand, CommandError
from kubernetes.stream import stream

from .utils import (
    apply_k8s_namespace,
    apply_k8s_secret,
    apply_k8s_service,
    apply_k8s_stateful_set,
    create_pg_database,
    create_pg_database_role,
    new_password,
)


ALREADY_EXISTS = HTTPStatus.CONFLICT

NAMESPACE = 'monokube'


class Command(BaseCommand):
    help = 'Initialize management environment'

    def add_arguments(self, parser):
        pass

    def handle(self, *args, **options):

        msg = '==== Initialize management environment ===='
        self.stdout.write(self.style.SUCCESS(msg))

        # ==== Namespace ====
        msg = f'#1 Create Namespace named {NAMESPACE}'
        self.stdout.write(self.style.SUCCESS(msg))
        apply_namespace()

        # ==== PostgreSQL ====
        msg = '#2 Install PostgreSQL'
        self.stdout.write(self.style.SUCCESS(msg))
        apply_postgres()

        # ==== Grafana ====
        msg = '#3 Install Grafana'
        self.stdout.write(self.style.SUCCESS(msg))
        # exec = 'psql -U postgres'
        # dbname = 'grafana'
        # exec_command = [
        #     'bash',
        #     '-c',
        #     f'{exec} -tc "SELECT 1 FROM pg_database WHERE datname = \'{dbname}\'" | grep -q 1 || '
        #     f'{exec} -c "CREATE DATABASE {dbname} ENCODING utf8;"',
        # ]
        # print(exec_command)

        # from k8s import core
        # resp = stream(core.connect_get_namespaced_pod_exec,
        #               'postgres-0',
        #               NAMESPACE,
        #               container='postgres',
        #               command=exec_command,
        #               stderr=True, stdin=False,
        #               stdout=True, tty=False)
        # print("Response: " + resp)
        apply_grafana()

        # ==== Completed ====
        msg = 'Management environment is initialized'
        self.stdout.write(self.style.SUCCESS(msg))


def apply_namespace() -> None:
    y = f'''
        apiVersion: v1
        kind: Namespace
        metadata:
          name: {NAMESPACE}
        '''
    m = yaml.safe_load(y)
    apply_k8s_namespace(m)


def apply_postgres() -> None:
    name = 'postgres'

    # Yaml Manifests
    yamlpath = Path(__file__).resolve().parent / 'manifests/postgres.yaml'
    with open(yamlpath, 'r') as stream:
        objects = yaml.safe_load_all(stream)
        objects = [object for object in objects if object != None]
    manifests: dict((str, str), dict) = {}
    for obj in objects:
        manifests[(obj['kind'], obj['metadata']['name'])] = obj

    # Secret
    pw = new_password()
    pw = base64.b64encode(pw.encode()).decode()
    m = manifests[('Secret', name)]
    m['data']['POSTGRES_PASSWORD'] = pw
    apply_k8s_secret(NAMESPACE, m)

    # StatefulSet
    m = manifests[('StatefulSet', name)]
    apply_k8s_stateful_set(NAMESPACE, m)

    # Service
    m = manifests[('Service', name)]
    apply_k8s_service(NAMESPACE, m)


def apply_grafana() -> None:
    name = 'grafana'
    create_pg_database(NAMESPACE, name)
    #create_pg_database_role(NAMESPACE, name, name)
