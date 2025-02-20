# Usage
In root directory run:
```
   go run main.go
   go run main.go <filename>
```
or
Build
```
  go build main.go
  ./main.exe
  ./main.exe <filename>
```

# Example Code:
```
print "variable declaration and data types";
var a;
print a;

var b = 0;
print b;

var c = "string";
print c;

var d = true;
var e = false;
print e;

var f = d | e;
print f;

var g = true & false;
print g;

var h = 10;
var i = b >= h;
print i;

print "";
print "if statements";
if (d | e) {
    print "d or e is true";
} else {
    print "d and e are false";
}

print "";
print "Scope Demonstration";
var a = "global a";
var b = "global b";
var c = "global c";
{
    var a = "outer a";
    var b = "outer b";
    {
        var a = "inner a";
        print a;
        print b;
        print c;
    }
    print a;
    print b;
    print c;
}
print a;
print b;
print c;

print "";
print "for loop fibbonacci";
var a = 0;
var temp;
for (var b = 1; a < 8; b = temp + b) {
    print a;
    temp = a;
    a = b;
}

print "";
print "while loop powers of 2";
var a = 1;
var b = 0;

while (a < 32) {
    print a;
    b = a;
    a = a + b;
}

print"";
print "Function Demonstration";
fn hello(first, last) {
    print "Hello Function: " + first + " " + last + "!";
}

hello("hello", "world");

print"";
print"Recursive fibbonacci Function";
fn fib(n) {
    if (n <= 1) return n;
    return fib(n - 2) + fib(n - 1);
}
for (var i = 0; i < 6; i = i + 1) {
    print fib(i);
}

print "";
print "Native Function Test";
print clock();

print "";
print "Native Sleep Test";
sleepMS(500);
print clock();
```

# Resulting output
```
variable declaration and data types
NULL
0
string
FALSE
TRUE
FALSE
FALSE

if statements
d or e is true

Scope Demonstration
inner a
outer b
global c
outer a
outer b
global c
global a
global b
global c

for loop fibbonacci
0
1
1
2
3
5

while loop powers of 2
1
2
4
8
16

Function Demonstration
Hello Function: hello world!

Recursive fibbonacci Function
0
1
1
2
3
5

Native Function Test
0.0031184

Native Sleep Test
0.5039167
```
