FROM golang:1.20 AS build

WORKDIR /src
COPY main.go go.mod .
RUN CGO_ENABLED=0 go build -o /bin/app

FROM scratch
COPY --from=build /bin/app /bin/app
ENTRYPOINT ["/bin/app"]
