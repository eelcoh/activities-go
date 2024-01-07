FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
ENV PATH=/bin
LABEL org.opencontainers.image.authors="ehoekema@gmail.com"

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt


ADD activities activities

EXPOSE 8080
ENTRYPOINT ["/activities"]
