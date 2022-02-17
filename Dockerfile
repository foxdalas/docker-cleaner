FROM       alpine:3.15
MAINTAINER Maxim Pogozhiy <foxdalas@gmail.com>

COPY docker-cleaner /bin/docker-cleaner

ENTRYPOINT ["/bin/docker-cleaner"]
EXPOSE     9203
