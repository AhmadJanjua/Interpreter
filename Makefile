# Makefile for windows
# Compile the program
all:
	g++ -Wall src\main.cpp src\Almond.cpp src\Scanner.cpp src\Token.cpp  -o Almond
# Compile the debugging version
debug:
	g++ -g src\main.cpp src\Almond.cpp src\Scanner.cpp src\Token.cpp  -o debug