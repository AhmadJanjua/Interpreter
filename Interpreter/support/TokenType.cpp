// Enum for all the language token types

enum TokenType {
    // Language Literals x3
    IDENTIFIER, STRING, NUMBER,

    // Symbol based tokens x 11
    L_PAREN, R_PAREN, L_BRACE, R_BRACE,
    COMMA, DOT, MINUS, PLUS, SEMICOLON, SLASH, STAR,
    
    // Comparison based tokens x8
    NOT, NOT_EQUAL, EQUAL, EQUAL_EQUAL,
    GREATER, GREATER_EQUAL, LESS, LESS_EQUAL,

    // Keywords x17
    AND, OR, TRUE, FALSE, IF, ELSE, FOR, WHILE, NA,
    CLASS, VAR, FUNC, THIS, SUPER, EOF, PRINT, RETURN
};