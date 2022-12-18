# Makefile for windows

# Compile the program
all:
	g++ -Wall .\Interpreter\source\main.cpp -o Almond

# delete the executable from the directory
clean:
	 rm .\Almond.exe