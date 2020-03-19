FROM golang:1.12.0-alpine as build
ENV NAME oamctl
# Build oamctl
RUN apk add git alpine-sdk upx \
 && git clone https://github.com/oam-dev/oamctl.git ${NAME} \
 && cd ${NAME} \
 && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -ldflags='-w -s -extldflags "-static"' -o /go/bin/${NAME}\
 && upx --brute /go/bin/${NAME}

FROM busybox:latest
COPY --from=build /go/bin/${NAME} /${NAME}
ENTRYPOINT [ "/oamctl" ]
