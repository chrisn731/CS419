package main

import (
	"os"
	"fmt"
)

const (
	// Constants for the keystream generator
	multiplier = 1103515245
	increment = 12345

	// The blocksize of the cipher
	blockSize = 16
)

var (
	_seed uint64 = 0
)

func readFile(fname string) []byte {
	dat, err := os.ReadFile(fname)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return dat
}

// Seed function
func sdbm(str []byte) (hash uint64) {
	for i := 0; i < len(str); i++ {
		c := uint64(str[i])
		hash = c + (hash << 6) + (hash << 16) - hash
	}
	return
}

func linearCongruentGen(seed uint64) byte {
	return byte((multiplier * seed + increment) % 256)
}

func genBlockKeyStream() [blockSize]byte {
	var gen [blockSize]byte

	for i := 0; i < blockSize; i++ {
		ret := linearCongruentGen(_seed)
		_seed = uint64(ret)
		gen[i] = ret
	}
	return gen
}

func shuffleBytes(plaintext, keystream []byte) {
	for i := blockSize - 1; i >= 0; i-- {
		first := keystream[i] & 0x0f
		second := (keystream[i] >> 4) & 0x0f
		temp := plaintext[first]
		plaintext[first] = plaintext[second]
		plaintext[second] = temp
	}

}

func doDecryption(ciphertext []byte) []byte {
	var plaintext, prevCipher []byte

	iv := genBlockKeyStream()
	isFirst := true
	for len(ciphertext) != 0 {
		var tempBlock [blockSize]byte

		// XOR ciphertext with generated keystream
		keystream := genBlockKeyStream()
		for i := 0; i < blockSize; i++ {
			tempBlock[i] = ciphertext[i] ^ keystream[i]
		}

		// Shuffle the bytes
		shuffleBytes(tempBlock[:], keystream[:])

		// XOR our (almost) plaintext with previous block of cipher text.
		// If this is the first iteration, use the initialization vector
		if isFirst {
			prevCipher = iv[:]
			isFirst = false
		}
		for i := 0; i < blockSize; i++ {
			tempBlock[i] ^= prevCipher[i]
		}
		prevCipher = ciphertext[:blockSize]
		ciphertext = ciphertext[blockSize:]
		plaintext = append(plaintext, tempBlock[:]...)
	}
	// Strip off the padding
	paddingSize := int(plaintext[len(plaintext) - 1])
	plaintext = plaintext[:len(plaintext) - paddingSize]
	return plaintext
}

func main() {
	program := os.Args[0]
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Printf("Usage: %s password ciphertext plaintext\n", program)
		os.Exit(1)
	}

	// Initalize the keystream generator seed
	password := []byte(args[0])
	_seed = sdbm(password)

	ciphertext := readFile(args[1])
	outFile := args[2]

	ret := doDecryption(ciphertext)
	err := os.WriteFile(outFile, ret, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
