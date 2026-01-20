VERSION := $(shell jj log --template 'commit_id.short(8)' --no-graph --limit 1)
LDFLAGS := -X 'main.version=$(VERSION)'

.PHONY: build clean

build:
	go build -ldflags "$(LDFLAGS)"

clean:
	rm -f dol
