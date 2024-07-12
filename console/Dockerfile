# syntax=docker/dockerfile:1.7.0

ARG BUILDPLATFORM="linux/arm64"

FROM --platform=$BUILDPLATFORM python:3.12.4-slim AS base
LABEL co.ultary.image.authors="ghilbut@gmail.com"

WORKDIR /app

ENV LANG=C.UTF-8 \
    LC_ALL=C.UTF-8 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONFAULTHANDLER=1 \
    PYTHONUNBUFFERED=1


FROM base AS deps

RUN pip install pipenv

COPY Pipfile .
COPY Pipfile.lock .
RUN pipenv install --system


FROM base AS builder

COPY --from=deps /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY . ./

RUN apt-get update -y  \
 && python ./manage.py collectstatic --clear --no-input


FROM base AS release
LABEL repository="https://github.com/ultary/monokube"

CMD ["hypercorn", "--bind", "0.0.0.0:8000", "monokube.asgi:application"]

ENV DJANGO_SETTINGS_MODULE=monokube.settings \
    PYTHONDONTWRITEBYTECODE=0
EXPOSE 8000
VOLUME /app/static

COPY --from=deps /usr/local/lib/python3.12/site-packages /usr/local/lib/python3.12/site-packages
COPY --from=deps /usr/local/bin/hypercorn /usr/local/bin/hypercorn
COPY --from=builder /tmp/static    /app/static
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/monokube  /app/monokube
