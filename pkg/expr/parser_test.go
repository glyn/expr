/*
 * Copyright 2020 Go YAML Path Authors
 *
 * SPDX-License-Identifier: Apache-2.0
 */

package expr

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewNode(t *testing.T) {
	cases := []struct {
		name     string
		lexemes  []lexeme
		expected *node
		focus    bool // if true, run only tests with focus set to true
	}{
		{
			name:     "no lexemes",
			lexemes:  []lexeme{},
			expected: nil,
		},
		{
			name:    "integer",
			lexemes: []lexeme{"1"},
			expected: &node{
				lexeme:   "1",
				children: []*node{},
			},
		},
		{
			name:    "add",
			lexemes: []lexeme{"1", "+", "2"},
			expected: &node{
				lexeme: "+",
				children: []*node{
					{
						lexeme:   "1",
						children: []*node{},
					},
					{
						lexeme:   "2",
						children: []*node{},
					},
				},
			},
		},
		{
			name:    "subtract",
			lexemes: []lexeme{"1", "-", "2"},
			expected: &node{
				lexeme: "-",
				children: []*node{
					{
						lexeme:   "1",
						children: []*node{},
					},
					{
						lexeme:   "2",
						children: []*node{},
					},
				},
			},
		},
		{
			name:    "add and subtract",
			lexemes: []lexeme{"1", "+", "2", "-", "3"},
			expected: &node{
				lexeme: "-",
				children: []*node{
					{
						lexeme: "+",
						children: []*node{
							{
								lexeme:   "1",
								children: []*node{},
							},
							{
								lexeme:   "2",
								children: []*node{},
							},
						},
					},
					{
						lexeme:   "3",
						children: []*node{},
					},
				},
			},
		},
		// {
		// 	name:    "multiplication",
		// 	lexemes: []lexeme{"1", "*", "2"},
		// 	expected: &node{
		// 		lexeme: "*",
		// 		children: []*node{
		// 			{
		// 				lexeme:   "1",
		// 				children: []*node{},
		// 			},
		// 			{
		// 				lexeme:   "2",
		// 				children: []*node{},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:    "1*2+3",
		// 	lexemes: []lexeme{"1", "*", "2", "+", "3"},
		// 	expected: &node{
		// 		lexeme: "+",
		// 		children: []*node{
		// 			{
		// 				lexeme: "*",
		// 				children: []*node{
		// 					{
		// 						lexeme:   "1",
		// 						children: []*node{},
		// 					},
		// 					{
		// 						lexeme:   "2",
		// 						children: []*node{},
		// 					},
		// 				},
		// 			},
		// 			{
		// 				lexeme:   "3",
		// 				children: []*node{},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:    "1+2*3",
		// 	lexemes: []lexeme{"1", "+", "2", "*", "3"},
		// 	expected: &node{
		// 		lexeme: "+",
		// 		children: []*node{
		// 			{
		// 				lexeme:   "1",
		// 				children: []*node{},
		// 			},
		// 			{
		// 				lexeme: "*",
		// 				children: []*node{
		// 					{
		// 						lexeme:   "2",
		// 						children: []*node{},
		// 					},
		// 					{
		// 						lexeme:   "3",
		// 						children: []*node{},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:    "1+2+3",
		// 	lexemes: []lexeme{"1", "+", "2", "+", "3"},
		// 	expected: &node{
		// 		lexeme: "+",
		// 		children: []*node{
		// 			{
		// 				lexeme: "+",
		// 				children: []*node{
		// 					{
		// 						lexeme:   "1",
		// 						children: []*node{},
		// 					},
		// 					{
		// 						lexeme:   "2",
		// 						children: []*node{},
		// 					},
		// 				},
		// 			},
		// 			{
		// 				lexeme:   "3",
		// 				children: []*node{},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:    "1*2*3",
		// 	lexemes: []lexeme{"1", "*", "2", "*", "3"},
		// 	expected: &node{
		// 		lexeme: "*",
		// 		children: []*node{
		// 			{
		// 				lexeme:   "1",
		// 				children: []*node{},
		// 			},
		// 			{
		// 				lexeme: "*",
		// 				children: []*node{
		// 					{
		// 						lexeme:   "2",
		// 						children: []*node{},
		// 					},
		// 					{
		// 						lexeme:   "3",
		// 						children: []*node{},
		// 					},
		// 				},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	name:    "simple bracket",
		// 	lexemes: []lexeme{"(", "1", ")"},
		// 	expected: &node{
		// 		lexeme: "()",
		// 		children: []*node{
		// 			{
		// 				lexeme:   "1",
		// 				children: []*node{},
		// 			},
		// 		},
		// 	},
		// },
		// {
		// 	focus:   true,
		// 	name:    "(1+2)*3",
		// 	lexemes: []lexeme{"(", "1", "+", "2", ")", "*", "3"},
		// 	expected: &node{
		// 		lexeme: "*",
		// 		children: []*node{
		// 			{
		// 				lexeme: "()",
		// 				children: []*node{
		// 					{
		// 						lexeme: "+",
		// 						children: []*node{
		// 							{
		// 								lexeme:   "1",
		// 								children: []*node{},
		// 							},
		// 							{
		// 								lexeme:   "2",
		// 								children: []*node{},
		// 							},
		// 						},
		// 					},
		// 				},
		// 			},
		// 			{
		// 				lexeme:   "3",
		// 				children: []*node{},
		// 			},
		// 		},
		// 	},
		// },
	}

	focussed := false
	for _, tc := range cases {
		if tc.focus {
			focussed = true
			break
		}
	}

	for _, tc := range cases {
		if focussed && !tc.focus {
			continue
		}
		t.Run(tc.name, func(t *testing.T) {
			actual := newNode(tc.lexemes)
			if focussed {
				// sometimes easier to read this than a diff
				fmt.Println("Expected:")
				fmt.Println(tc.expected.String())
				fmt.Println("Actual:")
				fmt.Println(actual.String())
			}
			require.Equal(t, tc.expected, actual)
		})
	}

	if focussed {
		t.Fatalf("testcase(s) still focussed")
	}
}
