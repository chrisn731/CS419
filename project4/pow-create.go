package main

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	// The amount of iterations to limit our search to
	searchLimit = 1000000000
)

func readFile(fname string) []byte {
	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		panic(err)
	}
	return dat
}

// Parse the hash string and check how many zeros are leading it
// In this case, leading zeros applies to binary, NOT the actual '0'
func numLeadingZeros(hash []rune) int {
	// By looking at each character in the string, we can map a number
	// of leading zeros to a rune
	leadingZeros := map[rune]int{
		'0': 4, '1': 3, '2': 2, '3': 2,
		'4': 1, '5': 1, '6': 1, '7': 1,
		'8': 0, '9': 0, 'a': 0, 'b': 0,
		'c': 0, 'd': 0, 'e': 0, 'f': 0,
	}
	num := 0
	for _, r := range hash {
		toAdd := leadingZeros[r]
		num += toAdd
		// If we did not get 4, that means the leading zeros will be
		// stopped by a 1.
		if toAdd != 4 {
			break
		}
	}
	return num
}

func applySHA256(hash []rune) []rune {
	h := sha256.New()
	h.Write([]byte(string(hash)))
	ret := hex.EncodeToString(h.Sum(nil))
	return []rune(ret)
}

func generateProofOfWork(_hash string, nbits int) (work, newHash []rune, runs uint64) {
	hash := []rune(_hash)
	newHash = append(newHash, hash...)
	if numLeadingZeros(newHash) == nbits {
		return
	}

	// We want to stay in range of printable characters which is 33 (!) - 126 (~)
	newHash = append(newHash, 33)
	startidx := len(newHash) - 1
	curridx := startidx
	for numLeadingZeros(applySHA256(newHash)) < nbits {
		runs++
		if runs >= searchLimit {
			fmt.Printf("Failed to find proof of work in %d runs" +
					" increase search limit to continue " +
					" for longer.", runs)
			os.Exit(1)
		}
		// If we are out of bounds...
		if newHash[curridx] >= 127 {
			// check the closest variable we can still edit without
			// going out of bounds
			for curridx >= startidx && newHash[curridx] >= 127 {
				curridx--
			}
			// Check out of bounds...
			if curridx < startidx {
				// We went out of bounds, that means that this
				// length proof of work is not sufficient.
				// Extend by one byte and reset.
				newHash = append(newHash, 33)
				for i := startidx; i < len(newHash); i++ {
					newHash[i] = 33
				}
			} else {
				// We did not go out of bounds, that means we
				// can increment the closest variable in our
				// current proof of work and reset ONLY the
				// bytes from the current to the end.
				newHash[curridx]++
				for i := curridx + 1; i < len(newHash); i++ {
					newHash[i] = 33
				}
			}
			// We messed with current, reset it back
			curridx = len(newHash) - 1
		} else {
			// We are not out of the valid byte range, simply increment
			// to the next valid character
			for {
				newHash[curridx]++
				currval := newHash[curridx]
				if currval != '\'' && currval != '"' {
					break
				}
			}
		}
	}
	// Pull out the work we did
	work = newHash[startidx:len(newHash)]
	return
}

func main() {
	program := os.Args[0]
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Usage: %s nbits file\n", program)
		os.Exit(1)
	}

	nbits, err := strconv.Atoi(args[0])
	if err != nil {
		panic(err)
	}
	file := args[1]
	fileBytes := readFile(file)
	h := sha256.New()
	h.Write(fileBytes)
	initialHash := hex.EncodeToString(h.Sum(nil))
	start := time.Now()
	work, newHash, iters := generateProofOfWork(initialHash, nbits)
	duration := time.Since(start)

	shaHash := applySHA256(newHash)
	s := fmt.Sprintf("File: %s\n" +
		"Initial-hash: %s\n" +
		"Proof-of-work: %s\n" +
		"Hash: %s\n" +
		"Leading-zero-bits: %d\n" +
		"Iterations: %d\n" +
		"Compute-time: %f\n",
		file, initialHash, string(work), string(shaHash),
		numLeadingZeros(shaHash), iters, duration.Seconds())
	fmt.Printf("%s", s)
}
