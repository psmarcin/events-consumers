# dependencies
FROM golang:1.19-alpine as build-env

WORKDIR /src
ADD go.* /src/

RUN go mod download

ADD . /src

RUN go build -o ./app

# base
FROM golang:1.19-alpine
COPY --from=build-env /src/app /app
COPY --from=build-env /src/views /views/
ENV TEMPLATE_PATH="/"

CMD ["/app"]
