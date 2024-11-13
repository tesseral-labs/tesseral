FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build

RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

WORKDIR /work

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "-w -s" -o frontend ./cmd/frontend

FROM scratch

COPY --from=build /usr/share/zoneinfo /user/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /work/frontend /frontend

EXPOSE 8083
ENTRYPOINT ["/frontend"]