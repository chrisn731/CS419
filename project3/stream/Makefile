OUT = scrypt.go

all: $(OUT)
	go build $^

clean:
	ls -1 | grep -vE "(Makefile|*.go)" | xargs rm -f

