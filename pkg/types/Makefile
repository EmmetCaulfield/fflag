.PHONY: default
default:
	go build

more_test.go: gentests.sh
	bash $< > $@

.PHONY:test
test: types_test.go more_test.go
	go test

.PHONY: clean
clean:
	rm -f *~ more_test.go
