FROM quay.io/opencars/openalpr:alpine AS openalpr
FROM golang:1.19-alpine AS build

RUN apk update && apk upgrade && apk add --no-cache ca-certificates tesseract-ocr-dev tiff-dev zlib-dev libpng-dev jpeg-dev make gcc build-base

ENV LD_LIBRARY_PATH="/usr/local/lib"
ENV PKG_CONFIG_PATH="/usr/local/share/pkgconfig:/usr/local/lib/pkgconfig:/usr/local/lib64/pkgconfig"
ENV LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib64:/usr/lib64:/usr/lib:/lib64:/lib"

COPY --from=openalpr /usr/local/include /usr/local/include
COPY --from=openalpr /usr/local/lib /usr/local/lib
COPY --from=openalpr /usr/local/lib64 /usr/local/lib64
COPY --from=openalpr /usr/local/share/openalpr/ /usr/local/share/openalpr/

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN BLDDIR=/go/bin make build

FROM alpine:3.12

ENV LD_LIBRARY_PATH="/usr/local/lib"
ENV PKG_CONFIG_PATH="/usr/local/share/pkgconfig:/usr/local/lib/pkgconfig:/usr/local/lib64/pkgconfig"
ENV LD_LIBRARY_PATH="/usr/local/lib:/usr/local/lib64:/usr/lib64:/usr/lib:/lib64:/lib"

RUN apk update && apk upgrade && apk add --no-cache tesseract-ocr tiff libpng jpeg

WORKDIR /app

COPY --from=openalpr /usr/local/lib /usr/local/lib
COPY --from=openalpr /usr/local/lib64 /usr/local/lib64
COPY --from=openalpr /usr/local/share/openalpr/ /usr/local/share/openalpr/

COPY ./config/ua.conf /usr/local/share/openalpr/runtime_data/config/ua.conf
COPY ./config/ua.patterns /usr/local/share/openalpr/runtime_data/postprocess/ua.patterns
RUN cp /usr/local/share/openalpr/runtime_data/region/eu.xml /usr/local/share/openalpr/runtime_data/region/ua.xml

COPY --from=build /go/bin/ ./
COPY ./config/config.yaml ./config/config.yaml

EXPOSE 8080

CMD ["./server"]
