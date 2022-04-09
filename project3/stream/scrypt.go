package main

import (
	"fmt"
	"io"
	"os"
)

const (
	multiplier = 1103515245
	increment  = 12345
)

var (
	seed uint64 = 0
)

func linearCongruentGen() byte {
	ret := byte((multiplier * seed + increment) % 256)
	seed = uint64(ret)
	return ret
}

func sdbm(str []byte) (hash uint64) {
	for _, b := range str {
		c := uint64(b)
		hash = c + (hash << 6) + (hash << 16) - hash
	}
	return
}

func seedGenerator(password []byte) {
	seed = sdbm(password)
}

func doStreamCipher(i, o string) error {
	in, err := os.Open(i)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(o, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	buf := make([]byte, 4096)
	for {
		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			buf[i] ^= linearCongruentGen()
		}
		_, err = out.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	program := os.Args[0]
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Printf("Usage: %s password plaintext ciphertext\n" +
				"\tor\n" +
				"       %s password ciphertext plaintext\n",
				program, program)
		os.Exit(1)
	}

	password := []byte(args[0])
	seedGenerator(password)
	err := doStreamCipher(args[1], args[2])
	if err != nil {
		panic(err)
	}
}
