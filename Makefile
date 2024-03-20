PKG:=

.PHONY: default
default:
	go build

.PHONY: test
test: ./pkg/*
	go test ./pkg/*
	go test

.PHONY: clean
clean:
	rm -f *~
