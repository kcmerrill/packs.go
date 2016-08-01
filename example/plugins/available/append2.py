#!/usr/bin/python
import sys
import json
import md5

def main():
    if len(sys.argv) == 1:
        for line in sys.stdin:
            line = json.loads(line.strip("\n") )
            print json.dumps(line + "-kcwashere2")
    else:
        register = [{
                "action": "filter",
                "trigger": "filter_append",
                "priority":1
            }]
        print json.dumps(register)


main()
