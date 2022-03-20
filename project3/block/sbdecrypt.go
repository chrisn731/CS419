package main

import (
	"fmt"
	"os"
)

const (
	// Constants for the keystream generator
	multiplier = 1103515245
	increment  = 12345

	// The blocksize of the cipher
	blockSize  = 16
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

// Seed function
func sdbm(str []byte) (hash uint64) {
	for i := 0; i < len(str); i++ {
		c := uint64(str[i])
		hash = c + (hash << 6) + (hash << 16) - hash
	}
	return
}

func linearCongruentGen() byte {
	ret := byte((multiplier * seed + increment) % 256)
	seed = uint64(ret)
	return ret
}

func genBlockKeyStream() [blockSize]byte {
	var gen [blockSize]byte

	for i := 0; i < blockSize; i++ {
		gen[i] = linearCongruentGen()
	}
	return gen
}

func shuffleBytes(plaintext, keystream []byte) {
	for i := blockSize - 1; i >= 0; i-- {
		first := keystream[i] & 0x0f
		second := (keystream[i] >> 4) & 0x0f
		plaintext[first], plaintext[second] = plaintext[second], plaintext[first]
	}
}

func seedGenerator(password []byte) {
	seed = sdbm(password)
}

func doDecryption(ciphertext []byte) []byte {
	var plaintext, prevCipher []byte

	iv := genBlockKeyStream()
	// On our first iteration use the initialization vector
	prevCipher = iv[:]
	for len(ciphertext) != 0 {
		var tempBlock [blockSize]byte

		// XOR ciphertext with generated keystream
		keystream := genBlockKeyStream()
		for i := 0; i < blockSize; i++ {
			tempBlock[i] = ciphertext[i] ^ keystream[i]
		}

		// Shuffle the bytes
		shuffleBytes(tempBlock[:], keystream[:])

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
	seedGenerator(password)

	ciphertext := readFile(args[1])
	outFile := args[2]

	ret := doDecryption(ciphertext)
	err := os.WriteFile(outFile, ret, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
