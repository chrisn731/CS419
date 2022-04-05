#!/bin/sh

# Set the directories of where everything is
SAMPLES_DIR="./samples"
V_DIR="./vigenere"
S_DIR="./stream"
B_DIR="./block"

# Password to use in testing
PASSWORD="dirtybubble"

verify_file_exists() {
	if test ! -e $1
	then
		echo $1 "does not exist! Make sure it exists."
		exit 1
	fi
}

verify_dir_exists() {
	if test ! -d $1
	then
		echo $1 "does not exist! Ensure script variables are correctly set" &&
		exit 1
	fi
}

run_test() {
	local password=$1
	local encrypter=$2
	local decrypter=$3
	local encrypt_type=$4
	verify_file_exists $encrypter
	verify_file_exists $decrypter
	for f in `ls -1 $SAMPLES_DIR`; do
		local file_path=$SAMPLES_DIR/$f
		$encrypter $password $file_path cipher;
		$decrypter $password cipher plain;
		diff -q plain $file_path >/dev/null;
		local ret_val=$?
		if test $ret_val -eq 0
		then
			echo "Success: " $file_path "using" $encrypt_type
		else
			echo "Fail: " $file_path "using" $encrypt_type
		fi
	done
}

test_vigenere() {
	echo "##### Testing Vigenere #####"
	echo $PASSWORD > keyfile
	local encrypt=$V_DIR/vencrypt
	local decrypt=$V_DIR/vdecrypt
	local encrypt_type="vigenere encryption"
	run_test keyfile $encrypt $decrypt $encrypt_type
}

test_stream() {
	echo "##### Testing Stream #####"
	local encrypt=$S_DIR/scrypt
	local decrypt=$S_DIR/scrypt
	local encrypt_type="stream cipher encryption"
	run_test $PASSWORD $encrypt $decrypt $encrypt_type
}

test_block() {
	echo "##### Testing Block #####"
	local encrypt=$B_DIR/sbencrypt
	local decrypt=$B_DIR/sbdecrypt
	local encrypt_type="sb-encryption"
	run_test $PASSWORD $encrypt $decrypt $encrypt_type
}
print_usage() {
	cat << EOF
$0 [OPTIONS]
 Test Options:
  NONE - Runs all tests
  v - Run vigenere test
  s - Run stream test
  v - Run block test

 Utility Options:
  make - 'Make' all encryption/decryption executables
  clean - Clean the directory of all temporary files made by the tester
EOF
	exit 0
}

verify_dir_exists $V_DIR
verify_dir_exists $S_DIR
verify_dir_exists $B_DIR
verify_dir_exists $SAMPLES_DIR

case $1 in
	"clean")
		rm -f keyfile plain cipher
		echo "Directory cleaned."
		exit 0
		;;

	"make")
		echo "Making..."
		before_move=$PWD
		cd $V_DIR && make && cd $before_move
		cd $S_DIR && make && cd $before_move
		cd $B_DIR && make && cd $before_move
		echo "Done."
		exit 0
		;;

	"help")
		print_usage
		;;

	"v")
		test_vigenere
		;;
	"b")
		test_block
		;;
	"s")
		test_stream
		;;
	*)
		echo "~~~~~ Testing All Encryptions ~~~~~"
		test_vigenere
		test_stream
		test_block
		;;
esac
echo "##### Tests Completed #####"
