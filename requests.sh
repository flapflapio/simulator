#!/bin/bash

main() {
    curl -i -X POST http://localhost:8080/simulate?tape=aaba --data '
    {
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
    '
}

[[ ${BASH_SOURCE[0]} == $0 ]] && main "$@"
