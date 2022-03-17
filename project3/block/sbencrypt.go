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

func genBlockKeyStream() (keyGen [blockSize]byte) {
	for i := 0; i < blockSize; i++ {
		ret := linearCongruentGen(_seed)
		_seed = uint64(ret)
		keyGen[i] = ret
	}
	return
}

func shuffleBytes(plaintext, keystream []byte) {
	for i := 0; i < blockSize; i++ {
		first := keystream[i] & 0x0f
		second := (keystream[i] >> 4) & 0x0f
		temp := plaintext[first]
		plaintext[first] = plaintext[second]
		plaintext[second] = temp
	}

}

func doEncryption(plaintext []byte) []byte {
	var completeCipher []byte
	var prevCipher [blockSize]byte

	// Add padding where needed
	numBytesToPad := blockSize - (len(plaintext) % blockSize)
	for i := 0; i < numBytesToPad; i++ {
		plaintext = append(plaintext, byte(numBytesToPad))
	}

	iv := genBlockKeyStream()
	isFirst := true
	for len(plaintext) != 0 {
		var tempBlock [blockSize]byte

		// If this is the first iteration, we use our initalization vector
		if isFirst {
			prevCipher = iv
			isFirst = false
		}

		// Apply CBC
		for i := 0; i < blockSize; i++ {
			tempBlock[i] = plaintext[i] ^ prevCipher[i]
		}
		// Read 16 bytes from the keystream
		keystream := genBlockKeyStream()

		// Shuffle the bytes based on keystream data
		shuffleBytes(tempBlock[:], keystream[:])

		// XOR between CBC'd block and the keystream
		for i := 0; i < blockSize; i++ {
			prevCipher[i] = tempBlock[i] ^ keystream[i]
		}
		// Add our ciphered block onto the rest
		completeCipher = append(completeCipher, prevCipher[:]...)

		// Advance the plaintext forward
		plaintext = plaintext[blockSize:]
	}
	return completeCipher
}


func main() {
	program := os.Args[0]
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Printf("Usage: %s password plaintext ciphertext\n", program)
		os.Exit(1)
	}

	// Initalize the keystream generator seed
	password := []byte(args[0])
	_seed = sdbm(password)

	plaintext := readFile(args[1])
	outFile := args[2]

	ret := doEncryption(plaintext)
	err := os.WriteFile(outFile, ret, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
