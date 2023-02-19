/**
 * Almond.cpp
 * Author: Ahmad Janjua
 * Purpose: Main simply reads in any commandline arguements that are passed to it and acts accordingly
*/
#include "../include/Almond.hpp"
#include <iostream>
using namespace std;

int main(int argc, char const *argv[])
{
    // If more than just the file name is passed, exit
    if(argc > 2) {
        cout << "Almond Usage: too many commandline arguements!" << endl;
        exit(64);
    }
    // If the filename is passed, attempt to run a file
    else if( argc == 2)
        Almond::runFile(argv[1]);
    // Run line by line otherwise
    else
        Almond::runPrompt();
    return 0;
}