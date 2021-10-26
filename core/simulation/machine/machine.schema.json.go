package machine

// Compiled version of the schema for use as a default schema - needs to be
// updated whenever the master copy gets updated
const SCHEMA = `
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://machinist.flapflap.io/machine.schema.json",
  "title": "Machine",
  "description": "A graph datastructure representing a state machine",
  "type": "object",

  "properties": {
    "Type": {
      "description": "What type of state machine this is. Must be one of: 'DFA', 'NFA', 'PDA', or 'TM'",
      "type": "string",
      "pattern": "(DFA|NFA|PDA|TM)"
    },

    "Alphabet": {
      "description": "The symbols that are accepted by the machine. This is a string where every character is a valid symbol accepted by the machine. If this field is omitted, then the alphabet will be inferred from the Transitions field.",
      "type": "string"
    },

    "Start": {
      "description": "The 'Id' field for the starting state of the machine",
      "type": "string",
      "pattern": "q([1-9]\\d*|0)"
    },

    "States": {
      "description": "The collection of states that are part of the machine",
      "type": "array",
      "minItems": 0,
      "uniqueItems": true,
      "items": {
        "type": "object",
        "properties": {
          "Id": {
            "description": "The id (unique) of the state e.g. 'q0', 'q1'. No leading zeros.",
            "type": "string",
            "pattern": "q([1-9]\\d*|0)"
          },
          "Ending": {
            "description": "Whether or not this state is an ending state. If absent, this value should be considered 'false'",
            "type": "boolean"
          }
        },
        "required": ["Id"]
      }
    },

    "Transitions": {
      "description": " The collection of transitions that are part of the machine",
      "minItems": 0,
      "uniqueItems": true,
      "items": {
        "type": "object",
        "properties": {
          "Start": {
            "description": "The 'Id' field for the starting state of the transition",
            "type": "string",
            "pattern": "q([1-9]\\d*|0)"
          },
          "End": {
            "description": "The 'Id' field for the ending state of the transition",
            "type": "string",
            "pattern": "q([1-9]\\d*|0)"
          },
          "Symbol": {
            "description": "The symbol(s) that is consumed from the input tape in order to traverse this transition",
            "type": "string"
          }
        },
        "required": ["Start", "End", "Symbol"]
      }
    }
  },

  "required": ["Start", "States", "Transitions", "Type"]
}
`
