See [Let's Build a Compiler](https://compilers.iecc.com/crenshaw/) by Jack Crenshaw.

This code parses expressions of the form:
```
<expression> ::= <term> [<addop> <term>]*
<term> ::= <factor> [ <mulop> <factor> ]*
<addop> ::= "+" | "-"
<mulop> ::= "*" | "/"
<factor> ::= <digit> | "(" <expression> ")"
```
