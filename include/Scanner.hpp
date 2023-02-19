/**
 * Scanner.hpp
 * Author: Ahmad Janjua
 * Purpose: Reads in the line of code and tokenizes the string
*/
#ifndef SCANNER_H
#define SCANNER_H

#include "TokenType.hpp"
#include "Almond.hpp"
#include "Token.hpp"
#include <string>
#include <map>
#include <list>
#include <variant>
using namespace std;

class Scanner {
public:
    // Fields

    uint64_t start{0}, current{0}, line{1}; // Helper fields in tracking the start, end and lines
    const string source; // The single line of code
    list<Token> tokens; // List of all the tokens that are extracted from source.
    static map<string, TokenType> keywords; // A map of keywords to the TokenType Enums

    // Methods

    // CONSTRUCTOR
    // PROMISES: creates a scanner object.
    Scanner(string& source);
    // PROMISES: Read the source code or the line and parse every ticket and add it to the token list.
    list<Token> scanTokens();
    // REQUIRES: Requires the type of token and the actual information inside of it, if any
    // PROMISES: Creates a new token object and adds it to the list of tokens
    void addToken(TokenType type, variant<monostate, bool, double, string> literal=monostate{});
    // PROMISES: Parses through the string and finds distinct token types and their info
    void scanToken();
    // PROMISES: Checks to see if token is an Identifier (variable name) or a keyword
    void identifier();
    // PROMISES: Parses through a number token and adds it to the tokens list.
    void number();
    // PROMISES: Parses through a string token and adds it to the tokens list.
    void parseString();
    // PROMISES: returns the character that the current index is at. If its at the end it will return a null character.
    char peek();
    // PROMISES: returns the next character after the current index and returns null character if its beyond the end.
    char peekNext();
    // PROMISES: returns the current character and increments the current index.
    char advance();
    // PROMISES: checks if the current index is beyond the scope of the input string or line.
    bool isAtEnd();
    // PROMISES: checks if the current character matches the expected character then increments the current index if theres a match.
    bool match(char expected);
};

#endif // SCANNER_H
