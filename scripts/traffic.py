#!/usr/bin/env python3
"""Gera tráfego contínuo nos módulos de exemplo para popular os dashboards do Grafana.

Cada ciclo exercita o que os dashboards do colibri-sdk-go mostram:

- HTTP: criação e consulta de cursos, alunos e matrículas no school-module, contas e
  faturas no finantial-module, mais erros intencionais (404 de curso inexistente e 400 de
  id inválido) para os painéis de status code;
- SQL e Redis: as próprias consultas dos módulos;
- messaging: a criação de cursos, alunos e matrículas publica eventos que o
  finantial-module consome (messaging.published, consumed, process.duration, consume.lag);
- storage: upload de um documento por aluno (storage.operation, transferred).

Uso:
    python3 scripts/traffic.py                # roda até Ctrl+C
    python3 scripts/traffic.py --duration 300 # roda por 5 minutos

Só usa a biblioteca padrão do Python 3.
"""

import argparse
import json
import random
import signal
import sys
import time
import urllib.error
import urllib.request
import uuid
from collections import Counter

BOUNDARY = "colibri-traffic-boundary"


class Client:
    def __init__(self, timeout):
        self.timeout = timeout
        self.stats = Counter()

    def request(self, method, url, route, body=None, raw=None, content_type="application/json"):
        """Envia a requisição e devolve (status, corpo). Erros de rede viram status 0."""
        data = raw if raw is not None else (json.dumps(body).encode() if body is not None else None)
        req = urllib.request.Request(url, data=data, method=method, headers={"Content-Type": content_type})
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                status, payload = resp.status, resp.read()
        except urllib.error.HTTPError as err:
            status, payload = err.code, err.read()
        except (urllib.error.URLError, TimeoutError, ConnectionError):
            status, payload = 0, b""

        self.stats[(method, route, status)] += 1
        return status, payload


def json_field(payload, *path):
    try:
        value = json.loads(payload)
        for key in path:
            value = value[key]
        return value
    except (ValueError, KeyError, IndexError, TypeError):
        return None


def cycle(client, school, finantial):
    suffix = uuid.uuid4().hex[:8]

    status, payload = client.request("POST", f"{school}/public/v1/courses", "/public/v1/courses",
                                     {"name": f"Course {suffix}", "value": random.randint(100, 900)})
    course_id = json_field(payload, "id") if status < 300 else None
    client.request("GET", f"{school}/public/v1/courses", "/public/v1/courses")

    status, payload = client.request("POST", f"{school}/public/v1/students", "/public/v1/students",
                                     {"name": f"Student {suffix}", "email": f"{suffix}@example.com",
                                      "birthday": "2000-01-01"})
    student_id = json_field(payload, "id") if status < 300 else None
    if student_id is None:
        _, payload = client.request("GET", f"{school}/public/v1/students?page=1&pageSize=1&name={suffix}",
                                    "/public/v1/students")
        student_id = json_field(payload, "items", 0, "id")
    client.request("GET", f"{school}/public/v1/students?name=Student&page=1&pageSize=10", "/public/v1/students")

    # erros intencionais, para os painéis de status code
    client.request("GET", f"{school}/public/v1/courses/{uuid.uuid4()}", "/public/v1/courses/:id")
    client.request("GET", f"{school}/public/v1/students/not-a-uuid", "/public/v1/students/:id")

    if course_id and student_id:
        client.request("POST", f"{school}/public/v1/enrollments", "/public/v1/enrollments",
                       {"courseId": course_id, "studentId": student_id, "installments": random.randint(1, 6)})

        content = "x" * random.randint(100, 50_000)
        body = (f"--{BOUNDARY}\r\n"
                f'Content-Disposition: form-data; name="file"; filename="document-{suffix}.txt"\r\n'
                "Content-Type: text/plain\r\n\r\n"
                f"{content}\r\n--{BOUNDARY}--\r\n").encode()
        client.request("POST", f"{school}/public/v1/students/{student_id}/upload-document",
                       "/public/v1/students/:id/upload-document", raw=body,
                       content_type=f"multipart/form-data; boundary={BOUNDARY}")

    client.request("GET", f"{school}/public/v1/enrollments?page=1&pageSize=10", "/public/v1/enrollments")
    client.request("GET", f"{finantial}/public/accounts", "/public/accounts")
    client.request("GET", f"{finantial}/public/invoices", "/public/invoices")


def print_summary(client, elapsed):
    total = sum(client.stats.values())
    print(f"\n{total} requisições em {elapsed:.0f}s ({total / max(elapsed, 1):.1f}/s)")
    for (method, route, status), count in sorted(client.stats.items()):
        label = "falha de conexão" if status == 0 else status
        print(f"  {method:6} {route:45} {label}: {count}")


def main():
    parser = argparse.ArgumentParser(description=__doc__.split("\n\n")[0])
    parser.add_argument("--duration", type=float, default=0, help="segundos de execução; 0 roda até Ctrl+C")
    parser.add_argument("--interval", type=float, default=0.2, help="pausa entre ciclos, em segundos")
    parser.add_argument("--school", default="http://localhost:8080", help="URL base do school-module")
    parser.add_argument("--finantial", default="http://localhost:8081", help="URL base do finantial-module")
    parser.add_argument("--timeout", type=float, default=10, help="timeout de cada requisição, em segundos")
    args = parser.parse_args()

    # SIGTERM (docker stop, pkill) encerra como o Ctrl+C, imprimindo o resumo
    signal.signal(signal.SIGTERM, lambda *_: sys.exit(0))

    client = Client(args.timeout)
    start = time.monotonic()
    deadline = start + args.duration if args.duration > 0 else None
    print(f"Gerando tráfego em {args.school} e {args.finantial}"
          + (f" por {args.duration:.0f}s" if deadline else " até Ctrl+C"), flush=True)

    try:
        while deadline is None or time.monotonic() < deadline:
            cycle(client, args.school, args.finantial)
            time.sleep(args.interval)
    except (KeyboardInterrupt, SystemExit):
        pass
    finally:
        print_summary(client, time.monotonic() - start)


if __name__ == "__main__":
    main()
