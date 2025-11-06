# fileshare
A file sharing website written in Go

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
