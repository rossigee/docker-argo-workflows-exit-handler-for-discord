FROM golang:1.20 AS build

WORKDIR /src
COPY main.go go.mod .
RUN CGO_ENABLED=0 go build -o /bin/handler

FROM scratch
COPY --from=build /bin/handler /bin/handler
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ['/bin/handler']
