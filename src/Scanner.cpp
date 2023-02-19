/**
 * Scanner.cpp
 * Author: Ahmad Janjua
*/
#include "../include/Scanner.hpp"
// Fields
map<string, TokenType> Scanner::keywords = {
    {"Class", CLASS}, {"Fun", FUN}, {"Var", VAR},
    {"True", TRUE}, {"False", FALSE},
    {"this", THIS}, {"super", SUPER}, {"return", RETURN},
    {"for", FOR}, {"while", WHILE},
    {"if", IF}, {"else", ELSE},
    {"AND", AND}, {"OR", OR},
    {"Print", PRINT},
    {"NIL", NIL}
};

// Methods
// ctor
Scanner::Scanner(string& source): source(source){}

list<Token> Scanner::scanTokens() {
    // keep checking for tokens till the end of the string is reached
    while (!isAtEnd()) {
        start = current;
        scanToken();
    }
    // add an END token to show that it is the last token and all others have been found.
    tokens.emplace_back(Token(TokenType::END, "", monostate(), line));
    return tokens;
}

void Scanner::addToken(TokenType type, variant<monostate, bool, double, string> literal) {
    // Parse the text from start to end of the token
    string text = source.substr(start, current - start);
    // add it to the list with appropriate type and information
    tokens.emplace_back(Token(type, text, literal, line));
}

void Scanner::scanToken() { 
    // have a switch statement to deal with all scenarios of inputs
    char c = advance();
    switch (c) {
        // Single characters
        case '(':
            addToken(L_PAREN); break;
        case ')':
            addToken(R_PAREN); break;
        case '{':
            addToken(L_BRACE); break;
        case '}':
            addToken(R_BRACE); break;
        case ',':
            addToken(COMMA); break;
        case '.':
            addToken(DOT); break;
        case '-':
            addToken(MINUS); break;
        case '+':
            addToken(PLUS); break;
        case ';':
            addToken(SEMICOLON); break;
        case '*':
            addToken(STAR); break; 
        case '!':
            addToken(match('=') ? BANG_EQUAL : BANG); break;
        case '=':
            addToken(match('=') ? EQUAL_EQUAL : EQUAL); break;
        case '<':
            addToken(match('=') ? LESS_EQUAL : LESS); break;
        case '>':
            addToken(match('=') ? GREATER_EQUAL : GREATER); break;
        // Check for comments //
        case '/':
            if (match('/')) {
                while (peek() != '\n' && !isAtEnd())
                    advance();
                } else {
                    addToken(SLASH);
                }
            break;
        // Check for whitespace
        case ' ':
        case '\r':
        case '\t':
            break;
        case '\n':
            line++; break;
        // Check string
        case '"': 
            parseString(); break;
        default:
            // Check if its a number
            if (isdigit(c)) {
                number();
            // otherwise its meant to be an identifier
            } else if (isalpha(c)) {
                identifier();
            // If nothing matches give an error
            } else {
                Almond::error(line, "Unexpected character.");
            break;
 }
            
    }
}

void Scanner::identifier() {
    // Default type is identifier
    TokenType type = IDENTIFIER;
    // loop past all numbers and characters while updating current index
    while (isalnum(peek()))
        advance();
    // parse the entire text extracted from start to finish
    string text = source.substr(start, current - start);
    // Check if the text matches a keyword and change the type to that
    if(keywords.find(text) != keywords.end()) {
        type = keywords[text];
        // check if the types were of boolean
        switch (type) {
        case TRUE:
            addToken(BOOLEAN, true); return;
        case FALSE:
            addToken(BOOLEAN, false); return;
        default:
            break;
        }
    }
    // Add the identifier or the keyword
    addToken(type);
}

void Scanner::number() {
    // loop through past all other numbers
    while (isdigit(peek()))
        advance();
    // check if its a fraction (check for decimal point)
    if (peek() == '.' && isdigit(peekNext())) {
        advance();
        // check for any numbers after the decimal point updating the current index again
        while (isdigit(peek())) 
            advance();
    }
    // add number to token list
    addToken(NUMBER, stod(source.substr(start, current - start)));
}

void Scanner::parseString() {
    // loop past all characters until string ends or there is a closing quote
    while (peek() != '"' && !isAtEnd()) {
        if (peek() == '\n') 
            line++;
        advance();
    }
    // If the end is reached and not the ending quote, report an error
    if (isAtEnd()) {
        Almond::error(line, "Unterminated string.");
        return;
    } else // otherwise move past end quote
    advance();
    // Trim the surrounding quotes and save the value
    string value = source.substr(start + 1, (current - start - 2) );
    addToken(STRING, value);
}

char Scanner::peek() {
    // check if end is reached
    if (isAtEnd())
        return '\0';
    // otherwise return the character at the current index
    else
        return source[current];
}

char Scanner::peekNext() {
    // check beyond the current index and return if there is a character there
    if (current + 1 >= source.length()) 
        return '\0';
    else
        return source[current + 1];
} 

char Scanner::advance() { 
    // return the current index while incrementing it
    return source[current++];
}

bool Scanner::isAtEnd() { return current >= source.length(); }

bool Scanner::match(char expected) {
    // check if the input character matches the current index
    if (isAtEnd() || (source[current] != expected) )
        return false;
    current++;
    return true;
}
