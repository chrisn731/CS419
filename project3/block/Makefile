all: encrypt decrypt

encrypt: sbencrypt.go
	go build $^

decrypt: sbdecrypt.go
	go build $^

clean:
	ls -1 | grep -vE "(Makefile|*.go)" | xargs rm -f

