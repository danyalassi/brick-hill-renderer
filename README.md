# FauxGL Renderer For Local

Fnproject was a stupid idea so here we go

## how to test

1. install [golang](https://go.dev/)

2. initalize go module
```bash
go mod init github.com/hawl1/brick-hill-renderer # you may want to change the link if you are gonna make changes
```

3. get packages
```bash
go get github.com/gofrs/uuid
go get github.com/hawl1/brickgl
go get github.com/spf13/pflag
go get github.com/nfnt/resize # i dont give a damn even if its not being developed
```

4. compile it down
```bash
go build
```

5. run the server
```bash
./brick-hill-renderer # this may depend because go puts the repo name from the link (i.e. github.com/hawl1/brick-hill-renderer)
```

6. test the server
```bash
curl -X POST -H "Content-Type: application/json" -d '{"avatarJSON":"","size":512}' http://0.0.0.0:8080/render
```

Changing settings and stuff
```bash
go run main.go --host 0.0.0.0 --port 7231 --accesskey change-it-or-die
```
- `--host` argument allows you to bind the ip that will get used
- `--port` argument allows you to bind the port that will get used
- `--accesskey` argument allows you to change the accesskey that will be used to access the renderer

That's all!
