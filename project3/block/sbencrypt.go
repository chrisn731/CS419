package main

import (
	"fmt"
	"io"
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

func genBlockKeyStream() (keyGen [blockSize]byte) {
	for i := 0; i < blockSize; i++ {
		keyGen[i] = linearCongruentGen()
	}
	return
}

func shuffleBytes(plaintext, keystream []byte) {
	for i := 0; i < blockSize; i++ {
		first := keystream[i] & 0x0f
		second := (keystream[i] >> 4) & 0x0f
		plaintext[first], plaintext[second] = plaintext[second], plaintext[first]
	}
}

func seedGenerator(password []byte) {
	seed = sdbm(password)
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
	// On our first iteration we should use our initalization vector
	prevCipher = iv
	for len(plaintext) != 0 {
		var tempBlock [blockSize]byte

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

func doEncryption1(i, o string) error {
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

	info, err := in.Stat()
	if err != nil {
		return err
	}
	sizeLeft := info.Size()

	buf := make([]byte, 4096)
	iv := genBlockKeyStream()
	prevCipher := iv[:]
	for sizeLeft > 0 {
		var ciphertext []byte
		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		sizeLeft -= int64(n)
		if sizeLeft <= 0 {
			numBytesToPad := blockSize - (n % blockSize)
			for i := 0; i < numBytesToPad; i++ {
				buf[n + i] = byte(numBytesToPad)
			}
			n += numBytesToPad
		}

		for i := 0; i < n; i += blockSize {
			start := i
			end := i + blockSize
			prevCipher = encryptBlock(buf[start:end], prevCipher)
			ciphertext = append(ciphertext, prevCipher...)
		}
		_, err = out.Write(ciphertext)
		if err != nil {
			return err
		}
	}
	return nil
}

func encryptBlock(plaintext, prevCipher []byte) []byte {
	var tempBlock [blockSize]byte

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
	return prevCipher
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
	seedGenerator(password)

	err := doEncryption1(args[1], args[2])
	if err != nil {
		panic(err)
	}
	/*
	plaintext := readFile(args[1])
	outFile := args[2]

	ret := doEncryption(plaintext)
	err := os.WriteFile(outFile, ret, 0644)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	*/
}
