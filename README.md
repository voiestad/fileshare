# Fileshare
A file sharing website written in Go.

## DISCLAIMER!
This program is currently in the early stages of development and is not ready for production usage.

## Developer
**Vebjørn Øiestad**

## Running Application
```
go run .
```

## Docker
```
docker build -t voiestad/fileshare .
docker run -p 8080:8080 \
  -v $(pwd)/files:/app/files \
  voiestad/fileshare
```
