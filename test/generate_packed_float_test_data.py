#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
generate_packed_float_test_data.py

generate data for unit tests using the actual python code from
twisted/spread/banana.py
"""
import random
import struct
import sys

_packed_float_format = "!d"

def main():
    """
    main entry point
    """
    test_floats = [-1000000.000001, -0,1, 0.0, 1.0, 2.0, ]
    for _ in range(10):
        test_floats.append(random.random())
    for test_float in test_floats:
        marshalled_float = struct.pack(_packed_float_format, test_float)
        print("            testItem{{{0}, []byte{{{1}}}}},".format(
            test_float, ", ".join([str(x) for x in marshalled_float])))

if __name__ == "__main__":
    sys.exit(main())
