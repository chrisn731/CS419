all: pow-cr pow-ch

pow-cr: pow-create.go
	go build $^

pow-ch: pow-check.go
	go build $^

clean:
	ls -1 | grep -vE "(p4-samples|Makefile|*.go)" | xargs rm -f
