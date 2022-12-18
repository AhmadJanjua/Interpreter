#include "iostream"
#include "fstream"
#include "sstream"
#include "string"
#include <list>

using namespace std;

// Requires: A path to a source file to read
// Promises: Will open and pass the contents of the file to be parsed
void runFile(string path);

// Promises: Will read inputs from terminal line by line
void runPrompt();

// Requires: a string of the contents from the source file
// Promises: parses through and extracts the tokens
void run(string& source);

// Main simply reads in any commandline arguements that are passed
int main(int argc, char const *argv[])
{
    // If more than just the file name is passed, exit
    if(argc > 2) {
        cout << "Usage: Too many arguements!" << endl;
        exit(64);
    }

    // If the filename is passed, run the file
    else if( argc == 2)
        runFile(argv[1]);
    
    // If nothing is passed, run line by line
    else
        runPrompt();
    
    return 0;
}

void runFile(string path) {
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

    return;
}

void runPrompt(){
    string line;

    // Read until nothing is received
    while(true) {
        // indicate to type here
        cout << ">> ";

        // Read line terminated by \n
        getline(cin, line);

        // Check if its empty to exit
        if(!line.compare("exit")) {
            // Confirm that the person wants to quit
            cout << "Would you like to exit? (Y/N)" << endl;
            getline(cin, line);
            // if they select y, close program
            if( !line.compare("y") || !line.compare("Y"))
                break;
        }
        // Otherwise run the line
        else{
            run(line);
        }
    }
    return;
}

void run(string& source) {
    string token;

    // Track all tokens; string seperated by whitespace
    list<string> tokens;

    // input string steam
    istringstream istring(source);

    // Extract a single token
    while (istring >> token)
    {
        // add token to list
        tokens.push_back(token);
        
        // Print token
        cout << token << endl;
    }
}