GOFILES := \
	$(wildcard *.go)

all: .web_static .genstatic_gen pingtrack

clean:
	-rm .genstatic_gen .web_static pingtrack

.web_static:
	cd ../web && gatsby build && cd ../cmd && touch .web_static

.genstatic_gen: $(wildcard ./cmd/*)
	go generate $(GOFILES) && touch .genstatic_gen

pingtrack: $(GOFILES) .genstatic_gen
	go build -o pingtrack $(GOFILES)