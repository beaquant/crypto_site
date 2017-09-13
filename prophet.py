#!/usr/bin/python3

import sys
import datetime

for line in sys.stdin:
    parts = line.strip().replace("\"", "").split(",")
    print(parts[1] + " " + parts[2])
    print(datetime.datetime.fromtimestamp(float(parts[1])).strftime('%Y-%m-%dT%H:%M:%SZ') + " " + parts[2])
