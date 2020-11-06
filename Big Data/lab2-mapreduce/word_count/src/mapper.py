#!/usr/bin/env python
from __future__ import print_function

# example using hadoop-streaming, adapted from
# http://www.michael-noll.com/tutorials/writing-an-hadoop-mapreduce-program-in-python/

import sys
import string

# input comes from STDIN (standard input)
for line in sys.stdin:

    # remove leading and trailing whitespace
    line = line.strip()

    # split the line into words
    words = line.split()

    # increase counters
    for word in words:
        # write the results to STDOUT (standard output);
        # what we output here will be the input for the
        # Reduce step, i.e. the input for reducer.py
        #
        # tab-delimited; the trivial word count is 1
        # the tab character, '\t' separates the key from the value
        w = word.strip(string.punctuation)
        if w:
            print('{}\t{}'.format(w, 1))
