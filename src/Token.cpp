/**
 * Token.hpp
 * Author: Ahmad Janjua
 * Purpose: ***
*/
#include "../include/Token.hpp"

// Methods
// ctor
Token::Token(TokenType type, string lexeme, variant<monostate, bool, double, string> literal, int line)
: type(type),lexeme(lexeme), literal(literal), line(line) {}

string Token::token_string() const{
    string temp_literal = "";
    // check what type is stored in the literal and extract the information to string
    // boolean
    if (holds_alternative<bool>(literal))
        temp_literal = to_string(get<bool>(literal));
    // double
    else if (holds_alternative<double>(literal))
        temp_literal = to_string(get<double>(literal));
    // string
    else if (holds_alternative<std::string>(literal))
        temp_literal = get<string>(literal);

    return TokenTypeString(type) + " " + lexeme + " " + temp_literal;
}

string Token::TokenTypeString(TokenType type) {
    // returns the string of the type passed in
    switch (type) {
        case TokenType::GREATER_EQUAL: return "Greater Equals: ";
        case TokenType::EQUAL_EQUAL: return "Equal Equal: ";
        case TokenType::LESS_EQUAL: return "Less Equals: ";
        case TokenType::IDENTIFIER: return "Identifier: ";
        case TokenType::BANG_EQUAL: return "Bang Equals: ";
        case TokenType::SEMICOLON: return "Semicolon: ";
        case TokenType::R_PAREN: return "Right Parenthesis: ";
        case TokenType::L_PAREN: return "Left  Parenthesis: ";
        case TokenType::R_BRACE: return "Right Brace: ";
        case TokenType::L_BRACE: return "Left  Brace: ";
        case TokenType::GREATER: return "Greater: ";
        case TokenType::BOOLEAN: return "Boolean: ";
        case TokenType::RETURN: return "Return: ";
        case TokenType::NUMBER: return "Number: ";
        case TokenType::STRING: return "String: ";
        case TokenType::SLASH: return "Slash: ";
        case TokenType::COMMA: return "Comma: ";
        case TokenType::MINUS: return "Minus: ";
        case TokenType::EQUAL: return "Equal: ";
        case TokenType::CLASS: return "Class: ";
        case TokenType::FALSE: return "False: ";
        case TokenType::PRINT: return "Print: ";
        case TokenType::SUPER: return "Super: ";
        case TokenType::WHILE: return "While: ";
        case TokenType::ELSE: return "Else: ";
        case TokenType::LESS: return "Less: ";
        case TokenType::PLUS: return "Plus: ";
        case TokenType::STAR: return "Star: ";
        case TokenType::BANG: return "Bang: ";
        case TokenType::THIS: return "This: ";
        case TokenType::TRUE: return "True: ";
        case TokenType::VAR: return "Var: ";
        case TokenType::FUN: return "Fun: ";
        case TokenType::NIL: return "Nil: ";
        case TokenType::END: return "End: ";
        case TokenType::AND: return "And: ";
        case TokenType::FOR: return "For: ";
        case TokenType::DOT: return "Dot: ";
        case TokenType::IF: return "If: ";
        case TokenType::OR: return "Or: ";
        default: return "UNKNOWN: ";
    }
}

// overload << operators so that token can be printed easily
ostream& operator<<(ostream& os, const Token& token) {
    os << token.token_string();
    return os;
}