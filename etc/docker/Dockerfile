# Build
FROM golang:1.19-alpine AS build

WORKDIR /app
COPY . .

RUN go build -ldflags="-s -w" -o hue-picker main.go

# Deploy
FROM alpine:3.17 AS deploy

LABEL maintainer="thomas@crissyfield.de"
LABEL description="Philips Hue Color Picker"

RUN apk --no-cache add tini \
                       tzdata

COPY --from=build /app/hue-picker /hue-picker

EXPOSE 80

ENTRYPOINT [ "/sbin/tini", "--", "/hue-picker" ]
CMD [ "serve" ]
