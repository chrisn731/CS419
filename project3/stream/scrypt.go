package main

import (
	"fmt"
	"os"
)

const (
	multiplier = 1103515245
	increment  = 12345
)

var (
	seed uint64 = 0
)

func readFile(fname string) []byte {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dat
}

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

func doStreamCipher(password, text []byte, outfile string) {
	var fullResult []byte
	for _, b := range text {
		fullResult = append(fullResult, b ^ linearCongruentGen())
	}
	err := os.WriteFile(outfile, fullResult, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func seedGenerator(password []byte) {
	seed = sdbm(password)
}

func main() {
	program := os.Args[0]
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Printf("Usage: %s password plaintext ciphertext\n" +
				"\tor\n" +
				"Usage: %s password ciphertext plaintext\n",
				program, program)
		os.Exit(1)
	}

	password := []byte(args[0])
	seedGenerator(password)
	in := readFile(args[1])
	outFile := args[2]
	doStreamCipher([]byte(password), in, outFile)
}

