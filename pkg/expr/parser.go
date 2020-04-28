/*
 * Copyright 2020 Go YAML Path Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package expr

import (
	"fmt"
	"strings"
)

type lexeme string

const (
	add          lexeme = "+"
	subtract     lexeme = "-"
	multiply     lexeme = "*"
	openBracket  lexeme = "("
	closeBracket lexeme = ")"
	eof          lexeme = "."
)

/*
   node represents a node of an expression parse tree. Each node is labelled with a lexeme.

   Terminal nodes have an integer lexeme.

   Non-terminal nodes represent an addition, multiplication, or bracketed expression.

   The following examples illustrate the approach.

   The basic expressions `1+2` is represented as the following parse tree:

               +
              / \
             1   2

   The expression `1 * 2 + 3` is represented as the parse tree:

                               +
                             /    \
                            *      3
                           / \
                          1   2

   The filter expression `(1+2)*3` is represented as the parse tree:

                           *
                        /     \
					  ( )      3
					   |
					   +
				     /   \
				    1     2
*/
type node struct {
	lexeme   lexeme
	children []*node
}

func newNode(lexemes []lexeme) *node {
	return newParser(lexemes).parse()
}

func (n *node) String() string {
	return "---\n" + n.indentedString(0) + "\n---\n"
}

func (n *node) indentedString(indent int) string {
	i := strings.Repeat("    ", indent)
	s := n.lexeme
	c := ""
	for _, child := range n.children {
		c = c + "\n" + child.indentedString(indent+1)
	}
	return fmt.Sprintf("%s%s%s", i, s, c)
}

// pstateFn represents the state of the parser as a function that returns the next state.
// A nil pstateFn indicates lexing is complete.
type pstateFn func(*parser) pstateFn

// parser holds the state of the filter expression parser.
// based on https://compilers.iecc.com/crenshaw/
type parser struct {
	input     []lexeme // the lexemes being scanned
	pos       int      // current position in the input
	stack     []*node  // parser stack
	tree      *node    // parse tree, equivalent of D0
	savedTree *node    // parse tree, equivalent of D1
}

// lex creates a new scanner for the input string.
func newParser(input []lexeme) *parser {
	l := &parser{
		input: input,
		stack: make([]*node, 0),
	}
	return l
}

// push pushes a node on the stack
func (p *parser) push(node *node) {
	p.stack = append(p.stack, node)
}

// pop pops a node from the stack. If the stack is empty, panics.
func (p *parser) pop() *node {
	if len(p.stack) == 0 {
		panic("parser stack underflow")
	}
	index := len(p.stack) - 1
	element := p.stack[index]
	p.stack = p.stack[:index]
	return element
}

// empty returns true if and onl if the stack of nodes is empty.
func (p *parser) emptyStack() bool {
	return len(p.stack) == 0
}

// nextLexeme returns the next item from the input.
func (p *parser) nextLexeme() lexeme {
	if p.pos >= len(p.input) {
		return eof
	}
	next := p.input[p.pos]
	p.pos++
	return next
}

// peek returns the next item from the input without consuming the item.
func (p *parser) peek() lexeme {
	if p.pos >= len(p.input) {
		return eof
	}
	return p.input[p.pos]
}

func (p *parser) getNum() *node {
	n := p.peek()
	if n == "0" || n == "1" || n == "2" || n == "3" || n == "4" || n == "5" ||
		n == "6" || n == "7" || n == "8" || n == "9" {
		p.nextLexeme()
		return &node{
			lexeme:   n,
			children: []*node{},
		}
	}
	panic("digit expected")
}

func (p *parser) match(m lexeme) {
	if p.peek() == m {
		p.nextLexeme()
		return
	}
	panic(fmt.Sprintf("%s expected but found %s", m, p.peek()))
}

func (p *parser) term() {
	p.tree = p.getNum()
}

func (p *parser) add() {
	p.match(add)
	p.term()
	p.tree = &node{
		lexeme: add,
		children: []*node{
			p.savedTree,
			p.tree,
		},
	}
}

func (p *parser) subtract() {
	p.match(subtract)
	p.term()
	p.tree = &node{
		lexeme: subtract,
		children: []*node{
			p.savedTree,
			p.tree,
		},
	}
}

func (p *parser) expression() {
	p.term()
	for p.peek() == "+" || p.peek() == "-" {
		p.savedTree = p.tree
		switch p.peek() {
		case add:
			p.add()
		case subtract:
			p.subtract()
		default:
			panic(fmt.Sprintf("+ or - expected, found %s", p.peek()))
		}
	}
}

func (p *parser) parse() *node {
	if p.peek() == eof {
		return nil
	}
	p.expression()
	return p.tree
}
