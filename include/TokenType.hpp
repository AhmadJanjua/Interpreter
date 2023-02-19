/**
 * TokenType.hpp
 * Author: Ahmad Janjua
 * Purpose: A list of keywords and functionalities.
*/
#ifndef TOKENTYPE_H
#define TOKENTYPE_H
#include <string>
using namespace std;

enum TokenType {
    // Single-character tokens
    L_PAREN, R_PAREN,
    L_BRACE, R_BRACE,
    COMMA, DOT, 
    MINUS, PLUS, SLASH, STAR,
    SEMICOLON,
    // One or two character tokens
    BANG, BANG_EQUAL,
    EQUAL, EQUAL_EQUAL,
    GREATER, GREATER_EQUAL,
    LESS, LESS_EQUAL,
    // Literals
    IDENTIFIER, STRING, NUMBER, BOOLEAN,
    // Keywords
    AND, OR,
    IF, ELSE, WHILE, FOR,
    CLASS, VAR, FUN,
    TRUE, FALSE,
    THIS, PRINT, 
    RETURN, SUPER,
    NIL, END
};

#endif // TOKENTYPE_H