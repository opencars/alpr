# ALPR

[godoc]: https://godoc.org/github.com/opencars/alpr
[godoc-img]: https://godoc.org/github.com/opencars/alpr?status.svg
[goreport]: https://goreportcard.com/report/github.com/opencars/alpr
[goreport-img]: https://goreportcard.com/badge/github.com/opencars/alpr
[version]: https://img.shields.io/github/v/tag/opencars/alpr?sort=semver

[![Docs][godoc-img]][godoc]
[![Go Report][goreport-img]][goreport]
[![Version][version]][version]

:camera: Golang web server for vehicle number plates recognition

## Development

Build docker image with needed libraries

```sh
docker build . -f Dockerfile.dev -t alpr-dev
```

Mount current directory and use bash inside docker container for development

```sh
docker run -it -v ${PWD}:/go/src/app -p 8080:8080 alpr-dev /bin/bash
```

## Usage

```sh
http http://localhost:8080/api/v1/alpr/private/recognize?image_url="https://example.com"
```

```json
[
    {
        "candidates": [
            {
                "confidence": 91.67285,
                "plate": "AA9359PC"
            },
            {
                "confidence": 84.00614,
                "plate": "AA9359PX"
            },
            {
                "confidence": 83.37887,
                "plate": "AA9359XC"
            },
            {
                "confidence": 83.22216,
                "plate": "AA935XPC"
            },
            {
                "confidence": 83.00714,
                "plate": "AA93X9PC"
            }
        ],
        "coordinates": [
            {
                "x": 243,
                "y": 496
            },
            {
                "x": 414,
                "y": 506
            },
            {
                "x": 413,
                "y": 542
            },
            {
                "x": 242,
                "y": 533
            }
        ],
        "plate": "AA9359PC"
    }
]
```

## License

Project released under the terms of the MIT [license](./LICENSE).
