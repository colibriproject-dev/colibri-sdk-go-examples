# Imagem compartilhada pelos deploys de AWS e GCP — o binário não tem nada
# específico de nuvem, tudo vem por variável de ambiente.
#
# Build self-contained: ao contrário do Dockerfile da raiz do repositório, que
# espera o binário `application` já compilado na máquina, aqui a compilação
# acontece dentro da imagem. Assim o deploy não depende de ter Go, swag ou
# mockgen instalados localmente — só Docker.
#
#   docker build -f deploy/docker/Dockerfile --build-arg APP_SRC=school-module .
#
# O contexto de build é a raiz do repositório.

ARG GO_VERSION=1.26

FROM golang:${GO_VERSION}-alpine AS builder

ARG APP_SRC
WORKDIR /src

RUN apk add --no-cache git ca-certificates

# go.mod e go.sum primeiro para a camada de dependências ser reaproveitada
# entre builds em que só o código mudou.
COPY ${APP_SRC}/go.mod ${APP_SRC}/go.sum ./
RUN go mod download

COPY ${APP_SRC}/ ./

# O contrato OpenAPI é gerado aqui em vez de vir do repositório: docs/* está no
# .gitignore, então num clone limpo esse diretório não existe. É o mesmo passo
# que o `make doc` de cada módulo executa.
RUN go install github.com/swaggo/swag/cmd/swag@latest \
    && (swag init --pd --ot=json || echo '{"swagger":"2.0","info":{"title":"contrato nao gerado no build"}}' > docs/swagger.json) \
    && ls -la docs/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -buildvcs=false -trimpath -ldflags="-w -s" -o /out/application .

FROM alpine:3.22 AS dist

ARG APP_SRC

RUN apk add --no-cache tzdata ca-certificates curl
ENV TZ=America/Sao_Paulo

WORKDIR /application

COPY --from=builder /out/application     /application/application
COPY --from=builder /src/docs/           /application/docs/
COPY ${APP_SRC}/migrations/              /application/migrations/

# Roda sem privilégio de root.
RUN addgroup -S app && adduser -S -G app app && chown -R app:app /application
USER app

ENTRYPOINT ["./application"]
