package main

import (
	"fmt"
	"os"
)

const (
	multiplier = 1103515245
	increment = 12345
)

func unused(x interface{}) { }

func readFile(fname string) []byte {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dat
}

func linearCongruentGen(seed uint64) byte {
	return byte((multiplier * seed + increment) % 256)
}

func sdbm(str []byte) (hash uint64) {
	for i := 0; i < len(str); i++ {
		c := uint64(str[i])
		hash = c + (hash << 6) + (hash << 16) - hash
	}
	return
}

func doStreamCipher(password, text []byte, outfile string) {
	file, err := os.Create(outfile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer file.Close()

	seed := sdbm([]byte(password))
	for i := 0; i < len(text); i++ {
		stream := linearCongruentGen(seed)
		result := text[i] ^ stream
		file.Write([]byte{result})
		seed = uint64(stream)
	}
}

func main() {
	args := os.Args[1:]

	if len(args) != 3 {
		fmt.Println("Usage: scrypt password ...")
		os.Exit(1)
	}
	password := args[0]
	in := readFile(args[1])
	outFile := args[2]
	doStreamCipher([]byte(password), in, outFile)
}

