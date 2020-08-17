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
http http://localhost:8080/api/v1/alpr/private/recognize?image_url=https://example.com
```

```json
{

}
```

## License

Project released under the terms of the MIT [license](./LICENSE).
