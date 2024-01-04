FROM scratch
LABEL org.opencontainers.image.authors="ehoekema@gmail.com"
ADD activities activities

EXPOSE 8080
ENTRYPOINT ["/activities"]
