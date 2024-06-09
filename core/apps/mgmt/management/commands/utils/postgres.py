from kubernetes.stream import stream
from k8s import core


def create_pg_database(namespace: str, dbname: str) -> None:
    psql_connect = 'psql -U postgres'
    exec_command = [
        'bash',
        '-c',
        f'{psql_connect} -tc "SELECT 1 FROM pg_database WHERE datname = \'{dbname}\'" | grep -q 1 || '
        f'{psql_connect} -c "CREATE DATABASE {dbname} ENCODING utf8;"',
    ]
    resp = stream(core.connect_get_namespaced_pod_exec,
                  'postgres-0',
                  namespace,
                  container='postgres',
                  command=exec_command,
                  stderr=True, stdin=False,
                  stdout=True, tty=False)
    print("Response: " + resp)


def create_pg_database_role(namespace:str, dbname: str, role: str, password: str) -> None:

    '''
    $ docker exec -it postgres psql -U postgres -c "CREATE DATABASE harbor ENCODING utf8;"
    $ docker exec -it postgres psql -U postgres -c "CREATE ROLE harbor WITH LOGIN PASSWORD '********';"
    $ docker exec -it postgres psql -U postgres -d harbor -c "GRANT ALL ON SCHEMA public TO harbor;"
    '''

    psql_connect = 'psql -U postgres'
    exec_command = [
        'bash',
        '-c',
        f'{psql_connect} -c "CREATE ROLE {role} WITH LOGIN PASSWORD \'{password}\';"',
    ]
    resp = stream(core.connect_get_namespaced_pod_exec,
                  'postgres-0',
                  namespace,
                  container='postgres',
                  command=exec_command,
                  stderr=True, stdin=False,
                  stdout=True, tty=False)
    print("Response: " + resp)

    exec_command = [
        'bash',
        '-c',
        f'{psql_connect} -d {dbname} -c "GRANT ALL ON SCHEMA public TO {role};"',
    ]
    resp = stream(core.connect_get_namespaced_pod_exec,
                  'postgres-0',
                  namespace,
                  container='postgres',
                  command=exec_command,
                  stderr=True, stdin=False,
                  stdout=True, tty=False)
    print("Response: " + resp)
