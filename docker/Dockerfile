FROM golang:alpine as builder

RUN apk update && apk upgrade && apk add --no-cache git openssh

COPY /ImgCrop /go/src/ImgCrop
WORKDIR /go/src/ImgCrop/cmd

RUN go get github.com/disintegration/imaging
RUN go get github.com/spf13/viper
RUN go get go.uber.org/zap
RUN go get github.com/gorilla/mux

RUN go build -v .
FROM alpine
RUN adduser -S -h /app app_user
COPY --from=builder /go/src/ImgCrop /app/
WORKDIR /app/cmd

USER root
RUN chmod -R 777 /app
RUN chown -R app_user /app

USER app_user
CMD ["./main"]
