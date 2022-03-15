package main

import (
	"fmt"
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

func doDecrypt(key, ciphertext []byte) (plaintext []byte) {
	for i := 0; i < len(ciphertext); i++ {
		plaintext = append(plaintext, ciphertext[i] - key[i])
	}
	return
}

func main() {
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Bad usage")
		os.Exit(0)
	}

	key := readFile(args[0])
	ciphertext := readFile(args[1])
	for len(key) < len(ciphertext) {
		key = append(key, key...)
	}
	output := args[2]
	err := os.WriteFile(output, doDecrypt(key, ciphertext), 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
