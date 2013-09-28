#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
generate_base_128_test_data.py

generate data for unit tests using the actual python code from
twisted/spread/banana.py
"""
import io
import random
import sys

_big_int = 1024 * 1024 * 1024 # sys.maxsize is too big for a go int literal

def _int2b128(integer, stream):
    if integer == 0:
        stream(chr(0))
        return
    assert integer > 0, "can only encode positive integers"
    while integer:
        stream(chr(integer & 0x7f))
        integer = integer >> 7


def _b1282int(st):
    e = 1
    i = 0
    for char in st:
        n = ord(char)
        i += (n * e)
        e <<= 7
    return i

def _marshall_int(n):
    accumulator = io.StringIO()
    _int2b128(n, accumulator.write)
    return accumulator.getvalue()

def main():
    """
    main entry point
    """
    test_ints = [0, 1, 2, 3, 4, 5, 6, 7, 88, 89, 100, 128, 256, _big_int]
    for _ in range(10):
        test_ints.append(random.randint(0, _big_int))
    for test_int in test_ints:
        marshalled_int = _marshall_int(test_int)
        print("            testItem{{{0}, []byte{{{1}}}}},".format(
            test_int, ", ".join([str(ord(x)) for x in marshalled_int])))

if __name__ == "__main__":
    sys.exit(main())
