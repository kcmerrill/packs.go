#!/usr/bin/python
import sys
import json
import md5

def main():
    if len(sys.argv) == 1:
        for line in sys.stdin:
            print json.dumps(
                    md5.new(
                        json.loads(
                            line.strip("\n")
                            )
                        ).hexdigest()
                    )
    else:
        register = [{
                "action": "filter",
                "trigger": "filter_password",
                "priority":1
            }]
        print json.dumps(register)


main()
