all: dev

prod:
	go build -o ./bin/api ./cmd/web/
	./bin/api -addr=":8080"

dev:
	air --build.cmd "go build -o ./bin/api ./cmd/web/" --build.bin './bin/api -addr=":1337"'
