FROM golang:1.11-stretch AS build
# Put source in GOPATH
RUN mkdir --parents /go/src/github.com/stuartleeks/validate-resource-limits
WORKDIR /go/src/github.com/stuartleeks/validate-resource-limits
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o validate-resource-limits


FROM alpine
WORKDIR /app
EXPOSE 8080
COPY --from=build /go/src/github.com/stuartleeks/validate-resource-limits /app
ENTRYPOINT [ "./validate-resource-limits" ]

