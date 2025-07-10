.PHONY: clean handd-run handd-tmpl handd-build handd-watch

clean:
	rm -rf ./bin

handd-run:
	go run ./cmd/handd/

handd-tmpl:
	go run ./cmd/handd/ -tmpl $(name)

handd-build:
	go build -o ./bin/handd ./cmd/handd/

handd-watch:
	air -c ./cmd/handd/.air.toml

.DEFAULT_GOAL := handd-watch
