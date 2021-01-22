FROM golang:1.15-stretch AS build-env
RUN mkdir -p /go/src/github.com/lreimer/k8s-experts-day/mutating-admission
WORKDIR /go/src/github.com/lreimer/k8s-experts-day/mutating-admission
COPY  . .
RUN useradd -u 10001 webhook
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o mutating-admission

FROM scratch
COPY --from=build-env /go/src/github.com/lreimer/k8s-experts-day/mutating-admission/mutating-admission .
COPY --from=build-env /etc/passwd /etc/passwd
USER webhook
ENTRYPOINT ["/mutating-admission"]