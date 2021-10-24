package dfa

const ODDA = `
{
	"Type": "DFA",
	"Alphabet": "ab",
	"Start": "q0",
	"States": [
	  { "Id": "q0", "Ending": false },
	  { "Id": "q1", "Ending": true }
	],
	"Transitions": [
	  { "Start": "q0", "End": "q1", "Symbol": "a" },
	  { "Start": "q0", "End": "q0", "Symbol": "b" },
	  { "Start": "q1", "End": "q1", "Symbol": "b" },
	  { "Start": "q1", "End": "q0", "Symbol": "a" }
	]
}
`
