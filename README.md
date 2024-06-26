# FauxGL Renderer For Local

Fnproject was a stupid idea so here we go

## how to test

1. install [golang](https://go.dev/)

2. install packages
```sh
go mod tidy
```

3. run server
```sh
go run main.go
```

4. test the server
```sh
curl -X POST -H "Content-Type: application/json" -d '{"avatarJSON":"","size":512}' http://0.0.0.0:8080/render
```

That's all!
By the way, you can compile it and you'll never need to compile again! (you can compile on freebsd and prob get it working on mydevil ?)
I won't be explaining that one though
