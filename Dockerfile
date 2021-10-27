FROM golang:1.17 as build

WORKDIR /shinobi_exporter

# Get deps (cached)
COPY ./go.mod /shinobi_exporter
COPY ./go.sum /shinobi_exporter
COPY ./Makefile /shinobi_exporter
RUN make dependencies

# Compile
COPY ./ /shinobi_exporter
RUN make test
RUN make lint
RUN make build
RUN ./shinobi_exporter --help

#############################################
# FINAL IMAGE
#############################################
FROM gcr.io/distroless/static
ENV LOG_JSON=1
COPY --from=build /shinobi_exporter /
USER 1000:1000
ENTRYPOINT ["/shinobi_exporter"]