# Every target here is a command, not a file it produces. `demo` in
# particular names a directory that exists, and without this make would
# report it up to date and record nothing.
.PHONY: build test cover demo

build:
	go build -o spg ./src

test:
	go test ./...

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Records the demo GIFs against the fixture harbour under demo/. Needs vhs and
# ffmpeg, which are not in the flake: they record the demo rather than build or
# test specgetty. Pass tapes to record only some of them:
#   make demo TAPES="demo/tapes/hero.tape"
demo:
	bash scripts/record-demo.sh $(TAPES)
