FROM alpine:latest as certs
RUN apk add -U --no-cache ca-certificates

FROM scratch
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY agbridge /usr/bin/agbridge
ENTRYPOINT [ "/usr/bin/agbridge" ]