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

func decryptBlock(ciphertext, prevCipher []byte) [blockSize]byte {
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
	return tempBlock
}

func doDecryption(i, o string) error {
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
	// Need to keep track of how many bytes left for padding
	sizeLeft := info.Size()

	buf := make([]byte, 4096)
	iv := genBlockKeyStream()
	prevCipher := iv[:]
	for sizeLeft > 0 {
		var plaintext []byte
		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		sizeLeft -= int64(n)
		for i := 0; i < n; i += blockSize {
			start := i
			end := i + blockSize
			ptxt := decryptBlock(buf[start:end], prevCipher)
			plaintext = append(plaintext, ptxt[:]...)
			prevCipher = buf[start:end]
		}
		// Need to hold a copy of the last block of ciphertext for next iteration
		prevCipher = make([]byte, blockSize)
		copy(prevCipher, buf[n - blockSize:n])

		// Remove padding
		if sizeLeft <= 0 {
			numBytesPadded := int(plaintext[n - 1])
			plaintext = plaintext[:len(plaintext) - numBytesPadded]
		}
		_, err = out.Write(plaintext)
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
		fmt.Printf("Usage: %s password ciphertext plaintext\n", program)
		os.Exit(1)
	}

	// Initalize the keystream generator seed
	password := []byte(args[0])
	seedGenerator(password)
	if err := doDecryption(args[1], args[2]); err != nil {
		panic(err)
	}
}
