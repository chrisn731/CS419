DIRS = vigenere stream block

all:
	for dir in $(DIRS); do 		\
		make -C $$dir;		\
	done

clean:
	for dir in $(DIRS); do 		\
		make -C $$dir clean;	\
	done
	./test.sh clean

test: all
	./test.sh

submission:
	@ zip cmn134-p3.zip block/*.go stream/*.go vigenere/*.go

.PHONY: all clean test submission
