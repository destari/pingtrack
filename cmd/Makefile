GOFILES := \
	$(wildcard *.go)

all: .web_static .genstatic_gen pinger

clean:
	-rm .genstatic_gen .web_static pinger

.web_static:
	cd ../web && gatsby build && cd ../cmd && touch .web_static

.genstatic_gen: $(wildcard ./cmd/*)
	go generate $(GOFILES) && touch .genstatic_gen

pinger: $(GOFILES) .genstatic_gen
	go build $(GOFILES)