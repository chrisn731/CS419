package main

import (
	"fmt"
	"io"
	"os"
)

// Read in the entire stream of bytes from a file
func readFile(fname string) []byte {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dat
}

func doDecrypt(key []byte, ctxt, ptxt string) error {
	f, err := os.Open(ctxt)
	if err != nil {
		return err
	}
	defer f.Close()

	pf, err := os.OpenFile(ptxt, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer pf.Close()

	buf := make([]byte, 4096)
	var nr uint64 // Need to keep track of which byte we are on
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			// M[i] = (C[i] - K[i % len(K)]) % len(alphabet)
			keylen := uint64(len(key))
			b := (int(buf[i]) - int(key[nr % keylen])) % 256
			buf[i] = byte(b)
			nr++
		}
		n, err = pf.Write(buf[:n])
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
		fmt.Printf("Usage: %s keyfile ciphertext plaintext\n", program)
		os.Exit(1)
	}

	key := readFile(args[0])
	err := doDecrypt(key, args[1], args[2])
	if err != nil {
		panic(err)
	}
}
