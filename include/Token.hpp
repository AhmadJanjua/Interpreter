/**
 * Token.hpp
 * Author: Ahmad Janjua
 * Purpose: Encapsulates a single unit of the language.
*/
#ifndef TOKEN_H
#define TOKEN_H


#include "TokenType.hpp"
#include <variant>
#include <ostream>
using namespace std;

class Token {
public:
    // Fields

    const TokenType type; // The type of token
    const string lexeme; // A sequence of characters forming the Token
    const variant<monostate, bool, double, string> literal; // The literal is one of those types or is empty
    const int line; // The line token is found at

    // Methods

    // PROMISES: Converts the Token to a string for printing
    string token_string() const;
    // REQUIRES: Requires info to form Token and also where it is.
    // PROMISES: Forms a Token object.
    Token(TokenType type, string lexeme, variant<monostate, bool, double, string> literal, int line);
    // REQUIRES: A TokenType enum (a number)
    // PROMISES: Returns the string name value of the enum
    static string TokenTypeString(TokenType type);
    // REQUIRES: An ostream and token object. Binary overloading.
    // PROMISES: Prints information about Token to console.
    friend ostream& operator<<(ostream& os, const Token& token);
};


#endif // TOKEN_H