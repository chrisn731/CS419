<h1 align="center">CS419 - Computer Security | Spring 2022</h1>

## Assignment 1 - Access Control
The goal of this project is to simulate access control mechanisms with things like:
<ul>
	<li> Users </li>
	<li> Domains </li>
	<li> Types </li>
	<li> Access Permissions </li>
</ul>
	
## Assignment 2 - Rootkit
### Part 1: Hiding Files
The first part of this project involves hiding files from a user when calling 'ls'.
This is done by simulating the functionality of readdir, but when given a filename that
we should hide we are not to return it to the user. File names to hide are put under a 
environment variable called "HIDDEN".

### Part 2: Changing Time
Similarly to part 1, this part involves overwriting the libc "time" function. On first run
of the program we are to lie to any program about the current time, and any subsequent call 
should return the true time.

## Assignment 3 - Encryption
### Part 1: Binary Vigenere Cipher
The first part of this project involves implementing a very simple cipher based on the vigenere 
cipher: 
<br>
`ciphertext[i] = (plaintext[i] + key[i % len(key)]) % len(alphabet)`

### Part 2: Stream cipher
Part 2 of this project requires implementing a keystream and using it in order to encrypt plaintext.

### Part 3: Block Encryption with Cipher Block Chaining and Padding
The final part builds off of part 2. It uses a combination of stream cipher and encrypting in blocks
in order to create a stronger cipher.

### Notes
For a more detailed explanation of the project check 'project3.pdf' under 'project3'.

## Assignment 4 - Proof-of-Work
### Part 1: Proof-of-Work Generator
The first part of this project requires building a proof-of-work generator. For any given set of text, the
program should be able to compute a hash that builds a string with a ceratin amount of leading zeros in 
front of the original hash.

### Part 2: Proof-of-Work Validator
The second part of this project requires building a validator for any proof-of-work generated. Whether it be
the generated material from part 1, or any given test cases.
