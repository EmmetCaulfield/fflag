PKG:=

.PHONY: default
default:
	go build

.PHONY: test
test: ./pkg/*
	go test ./pkg/*
	go test

README.md: flag.go
	@echo '# `fflag`' > $@
	go doc | sed -n '2,/^const/p' | head -n-1 >> $@

.PHONY: clean
clean:
	rm -f *~
