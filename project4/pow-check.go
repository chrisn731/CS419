package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"strconv"
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

func checkMissing(lines []string) []*string {
	headers := []string{
		"Initial-hash:",
		"Proof-of-work:",
		"Leading-zero-bits:",
		"Hash:",
	}
	var res []*string
	for _, header := range headers {
		var found_line *string = nil
		for _, line := range lines {
			if strings.Contains(line, header) {
				found_line = &line
				break
			}
		}
		res = append(res, found_line)
	}
	return res

}

func checkInitHash(line *string, initHash string) bool {
	if line == nil {
		fmt.Println("ERROR: missing Initial-hash in header")
		return true
	}
	val := strings.Split(*line, ": ")
	hashInHeader := ""
	if len(val) == 2 {
		hashInHeader = val[1]
	}
	if initHash == hashInHeader {
		fmt.Println("PASSED: initial file hashes match")
		return false
	} else {
		fmt.Printf("ERROR initial hashes don't match\n" +
			"\thash in header: %s\n" +
			"\tfile hash: %s\n",
			hashInHeader, initHash)
		return true
	}
}

func checkProofOfWork(line *string) bool {
	if line == nil {
		fmt.Println("ERROR: missing Proof-of-work in header")
		return true
	}
	return false
}

func checkLeadingZeroBits(line *string, hash string) bool {
	if line == nil {
		fmt.Println("ERROR: missing Leading-zero-bits in header")
		return true
	}
	val := strings.Split(*line, ": ")
	headerVal, _ := strconv.Atoi(val[1])
	numZeros := numLeadingZeros([]rune(hash))
	if numZeros == headerVal {
		fmt.Println("PASSED: leading bits is correct")
		return false
	} else {
		fmt.Printf("ERROR: Leading-zero-bits value: %s, but hash has %d" +
			" leading zero bits\n", val[1], numZeros)
		return true
	}
	return false
}

func checkProofHash(line *string, initHash, work string) bool {
	given := ""
	if line == nil {
		fmt.Println("ERROR: missing Hash in header")
		return true
	} else {
		val := strings.Split(*line, ": ")
		if len(val) == 2 {
			given = val[1]
		}
	}

	h := sha256.New()
	h.Write([]byte(initHash + work))
	hashString := hex.EncodeToString(h.Sum(nil))

	if given == hashString {
		fmt.Println("PASSED: pow hash matches Hash header")
		return false
	} else {
		fmt.Printf("ERROR: pow hash does not match Hash header\n" +
			"\texpected: %s\n" +
			"\theader has: %s\n",
			hashString, given)
		return true
	}

}

func doHeaderChecks(lines []*string, initHash string) bool {
	proofHash := ""
	if lines[3] != nil {
		t := strings.Split(*lines[3], ": ")
		if len(t) == 2 {
			proofHash = t[1]
		}
	}

	work := ""
	if lines[1] != nil {
		t := strings.Split(*lines[1], ": ")
		if len(t) == 2 {
			work = t[1]
		}
	}
	fail1 := checkInitHash(lines[0], initHash)
	fail2 := checkProofOfWork(lines[1])
	fail3 := checkLeadingZeroBits(lines[2], proofHash)
	fail4 := checkProofHash(lines[3], initHash, work)
	return fail1 || fail2 || fail3 || fail4
}

func main() {
	program := os.Args[0]
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Printf("Usage: %s powheader file\n", program)
		os.Exit(1)
	}


	powheader := args[0]
	file := args[1]
	h := sha256.New()
	h.Write(readFile(file))
	initialHash := hex.EncodeToString(h.Sum(nil))

	lines := strings.Split(string(readFile(powheader)), "\n")
	ret := checkMissing(lines)
	fail := doHeaderChecks(ret, initialHash)
	if fail {
		fmt.Println("fail")
	} else {
		fmt.Println("pass")
	}
}
