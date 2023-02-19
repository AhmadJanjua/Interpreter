/**
 * Almond.hpp
 * Author: Ahmad Janjua
 * Purpose: Takes in raw source code ad converts it into tokens
*/
#ifndef ALMOND_H
#define ALMOND_H

#include "../include/Scanner.hpp"
#include "TokenType.hpp"
#include "Token.hpp"
#include <iostream>
#include <fstream>
#include <sstream>
#include <string>
#include <list>
using namespace std;

class Almond {
public:
    // Fields

    static bool hadError; // Checks if there was an error in the program
    static bool hadRuntimeError; // Checks if there was an error in the runtime environment

    // Methods

    // REQUIRES: A location and message of where the error occured.
    // PROMISES: Returns a error message on the console.
    static void error(int line, string message);
    // REQUIRES: A token and a message
    // PROMISES: Returns a error message on the console.
    static void error(Token token, string message);
    // REQUIRES: A line and location of error as well as the error message
    // PROMISES: Returns a error message on the console.
    static void report(int line, string where, string message);
    // REQUIRES: a string of the contents from the source file
    // PROMISES: parses through and extracts the tokens
    static void run(string& source);
    // REQUIRES: A path to a source file to read
    // PROMISES: Will open and pass the contents of the file to be parsed
    static void runFile(const string& path);
    // PROMISES: Will read inputs from terminal line by line
    static void runPrompt();
};

#endif // ALMOND_H
