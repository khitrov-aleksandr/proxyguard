FROM golang AS build
ENV DEBIAN_FRONTEND noninteractive
ARG APP_USER_ID
WORKDIR /app
COPY . .
RUN go mod download
RUN useradd -m -s /bin/bash -u$APP_USER_ID app
RUN chown -R app\: .
USER app
RUN CGO_ENABLED=0 GOOS=linux go build -o /tmp/run *.go

FROM alpine AS recording
ARG APP_USER_ID
RUN apk update
RUN apk upgrade
RUN apk add --no-cache tzdata

RUN adduser -u$APP_USER_ID -D app
USER app

COPY --from=build /tmp/run /app/run
COPY --chown=app:app docker/docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT [ "/docker-entrypoint.sh" ]
