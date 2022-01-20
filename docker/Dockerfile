# STAGE 1: build the app
# ==============================================================================
FROM golang:1.17.1 AS build
ENV DEBIAN_FRONTEND="noninteractive"
WORKDIR /go/src/github.com/flapflapio/simulator/
COPY ./ /go/src/github.com/flapflapio/simulator/
RUN apt-get update                                          \
    && apt-get install -y make bash tzdata ca-certificates  \
    && make build-static                                    \
    && adduser                                              \
        --disabled-password                                 \
        --gecos ""                                          \
        --home "/homeless"                                  \
        --shell "/noshell"                                  \
        --no-create-home                                    \
        --uid 666                                           \
        simulator

# STAGE 2: Running the app :D
# ==============================================================================
FROM scratch AS final

# Some useful files to have
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/localtime /etc/localtime
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

# Our app files
COPY --from=build                                       \
    /go/src/github.com/flapflapio/simulator/app         \
    /go/src/github.com/flapflapio/simulator/config.yml  \
    /

# Healthcheck
HEALTHCHECK             \
    --interval=30s      \
    --timeout=30s       \
    --start-period=5s   \
    --retries=3         \
    CMD [ "/app", "-health" ]

EXPOSE 8080
ENV PORT=8080 \
    TZ='Canada/Eastern'
USER simulator:simulator
ENTRYPOINT ["/app"]
