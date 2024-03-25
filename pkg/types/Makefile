.PHONY: default
default: coerce.go types.go
	go build

more_test.go: gentests.sh
	bash $< > $@

coerce.go: gencoerce.sh coerce.go-top
	bash $< > $@

.PHONY:test
test: coerce.go types_test.go more_test.go
	go test

.PHONY: clean
clean:
	rm -f *~ more_test.go coerce.go
