
encrypt: vencrypt.go
	go build $^

decrypt: vdecrypt.go
	go build $^

clean:
	ls -1 | grep -vE "(Makefile|*.go)" | xargs rm -f

