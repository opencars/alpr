# ALPR

Golang web server for vehicle number plates recognition

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
http http://localhost:8080/api/v1/alpr/private/recognize?image_url="https://img03.platesmania.com/170327/o/9629055.jpg"
```

```json
[
    {
        "coordinates": [
            {
                "x": 185,
                "y": 421
            },
            {
                "x": 302,
                "y": 430
            },
            {
                "x": 299,
                "y": 459
            },
            {
                "x": 183,
                "y": 448
            }
        ],
        "plate": "AA9359PC"
    }
]
```

## License

Project released under the terms of the MIT [license](./LICENSE).
