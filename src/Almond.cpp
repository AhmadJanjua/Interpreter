/**
 * Almond.cpp
 * Author: Ahmad Janjua
*/
#include "../include/Almond.hpp"

// Fields
bool Almond::hadError = false;
bool Almond::hadRuntimeError = false;

// Methods
void Almond::error(int line, string message) {
    report(line, "", message);
}

void Almond::error(Token token, string message) {
    if (token.type == TokenType::END) {
        report(token.line, " at end", message);
    } else {
        report(token.line, " at '" + token.lexeme + "'", message);
    }
}

void Almond::report(int line, string where, string message) {
    cout <<  "[line " << line << "] Error" << where << ": " << message << endl;
    hadError = true;
}

void Almond::run(string& source) { 
    Scanner scanner(source);

    // Track all tokens; string seperated by whitespace
    list<Token> tokens = scanner.scanTokens();

    for( auto token : tokens) {
        cout << token << endl;
    }
}

void Almond::runFile(const string& path) {
    // open a read only input stream
    ifstream file(path);

    // output string buffer stream
    ostringstream ostring;

    // Read from file into string buffer stream
    ostring << file.rdbuf();  

    // Convert the string buffer stream into a string
    string source = ostring.str();

    // Pass the string into the run function
    run(source);

    if(Almond::hadError) exit(65);
    
    return;
}

void Almond::runPrompt(){
    string line;

    // Read until nothing is received
    while(true) {
        // indicate to type here
        cout << "> ";

        // Read line terminated by \n
        getline(cin, line);

        // Check if its empty to exit
        if(line.compare("") == 0)
            break;
        // Otherwise run the line
        run(line);
        Almond::hadError = false;

    }
    return;
}