FROM alpine

WORKDIR /app

RUN apk update
RUN apk add ca-certificates
RUN update-ca-certificates

ARG TARGETOS
ARG TARGETARCH
COPY bld/sciensano-$TARGETOS-$TARGETARCH /app/sciensano

COPY build/entrypoint.sh /app
RUN chmod 755 /app/entrypoint.sh

EXPOSE 8080
ENTRYPOINT [ "/app/entrypoint.sh" ]
CMD []
