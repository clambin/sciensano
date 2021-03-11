FROM alpine

WORKDIR /app

ARG TARGETOS
ARG TARGETARCH

COPY bld/sciensano-$TARGETOS-$TARGETARCH /app/sciensano

EXPOSE 8080
ENTRYPOINT ["/app/sciensano"]
CMD []
