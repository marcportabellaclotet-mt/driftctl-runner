FROM alpine:3.6 as alpine
ARG DRIFTCTL_VERSION=v0.37.0
ARG DRIFTCTL_SHA256=23ee969cbfb960b348fa42e1b6c9b171883bef1c3dfff8f5542cf696f30f071f
RUN apk add -U --no-cache ca-certificates
RUN addgroup -g 1001 -S driftctl-runner && adduser -u 1001 -S driftctl-runner  -G driftctl-runner
RUN pass=$(echo date +%s | sha256sum | base64 | head -c 32; echo | mkpasswd) && \
    echo "driftctl-runner:${pass}" | chpasswd
RUN apk update
RUN apk add curl
RUN pwd
RUN curl -L -o driftctl -s https://github.com/snyk/driftctl/releases/download/${DRIFTCTL_VERSION}/driftctl_linux_amd64
RUN echo "${DRIFTCTL_SHA256}  driftctl" | sha256sum -c -
RUN chmod +x driftctl

FROM scratch
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine /etc/passwd /etc/passwd
COPY --from=alpine /etc/group /etc/group
COPY --from=alpine /etc/shadow /etc/shadow
COPY --from=alpine /driftctl /usr/local/bin/driftctl

COPY --chown=driftctl-runner build/driftctl-runner_linux_amd64 /app

USER driftctl-runner
CMD ["/app"]
