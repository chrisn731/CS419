package main

import (
	"fmt"
	"os"
)

func fillCipherGrid(grid [][256]uint8) {
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[i]) - i; j++ {
			grid[i][j] = uint8(i + j)
		}
		for j := len(grid[i]) - i; j < len(grid[i]); j++ {
			grid[i][j] = uint8(j - (len(grid[i]) - i))
		}
	}
}

func readFile(fname string) []byte {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dat
}

func doEncrypt(key, plaintext []byte, cipher [256][256]uint8) (ciphertext []byte) {
	for i := 0; i < len(plaintext); i++ {
		ciphertext = append(ciphertext,
				byte(cipher[plaintext[i] % 255][key[i] % 255]))
	}
	return
}

func main() {
	var cipherGrid [256][256]uint8

	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("Bad usage")
		os.Exit(0)
	}

	fillCipherGrid(cipherGrid[:])
	key := readFile(args[0])
	plaintext := readFile(args[1])
	for len(key) < len(plaintext) {
		key = append(key, key...)
	}
	output := args[2]
	//fmt.Printf("Performing encrypt using key (%s) and plaintext (%s)\n",
			//string(key), string(plaintext))
	err := os.WriteFile(output, doEncrypt(key, plaintext, cipherGrid), 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
