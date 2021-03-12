FROM alpine:latest

ENV EXECUTABLE=sport4all
ENV CONFIG=docker_hosted_config.yml

ARG BUILDER=sport-builder
ARG TAG=latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app

COPY --from=sport-builder:latest /application/${EXECUTABLE} /app
COPY --from=sport-builder:latest /application/cmd/${CONFIG} /app

CMD /app/${EXECUTABLE} -c ${CONFIG}