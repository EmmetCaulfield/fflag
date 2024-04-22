PKG:=

.PHONY: default
default: README.md
	go build

.PHONY: test
test: ./pkg/*
	go test ./pkg/*
	go test

README.md: flag.go
	@echo '# `fflag`\n' > $@
	sed -n '0,/^$$/{s|^// \?||;p}' $< >> $@

.PHONY: clean
clean:
	rm -f README.md *~
