all: dev

build: 
	go build -o ./bin/api ./cmd/web/

prod:
	go build -o ./bin/api ./cmd/web/
	
dev:
	air --build.cmd "go build -o ./bin/api ./cmd/web/" --build.bin './bin/api -addr=":1337"'
