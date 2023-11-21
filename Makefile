all: dev

prod:
	air --build.cmd "go build -o ./bin/api ./cmd/web/" --build.bin './bin/api -addr=":80"'

test:
	air --build.cmd "go build -o ./bin/api ./cmd/web/" --build.bin "./bin/api"

dev:
	air --build.cmd "go build -o ./bin/api ./cmd/web/" --build.bin './bin/api -addr=":1337"'
