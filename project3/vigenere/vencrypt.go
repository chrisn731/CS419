package main

import (
	"fmt"
	"io"
	"os"
)

func readFile(fname string) []byte {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dat
}

func doEncrypt(key []byte, ptxt, ctxt string) error {
	f, err := os.Open(ptxt)
	if err != nil {
		return err
	}
	defer f.Close()

	cf, err := os.OpenFile(ctxt, os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer cf.Close()

	buf := make([]byte, 4096)
	var nr uint64
	for {
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i++ {
			// C[i] = (M[i] + K[i % len(K)]) % len(alphabet)
			keylen := uint64(len(key))
			buf[i] = byte((int(buf[i]) + int(key[nr % keylen])) % 256)
			nr++
		}
		n, err = cf.Write(buf[:n])
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
		fmt.Printf("Usage: %s keyfile plaintext ciphertext\n", program)
		os.Exit(1)
	}

	key := readFile(args[0])
	err := doEncrypt(key, args[1], args[2])
	if err != nil {
		panic(err)
	}
}
