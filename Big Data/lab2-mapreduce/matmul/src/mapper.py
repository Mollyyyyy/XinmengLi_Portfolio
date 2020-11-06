#!/usr/bin/env python
# map function for matrix multiply

# The input matrices should be contained in a folder in HDFS of the form
#   matrix/A.txt
#   matrix/B.txt
#
# to compute the product A * B

# Input files are assumed to have lines of the form "i,j,x", where
#   i is the row index
#   j is the column index
#   and x is the value M[i, j]
# for matrix M
#
# Indices are assumed to start at 0.

# It is assumed that the matrix dimensions are such that the product A*B exists.

#Input arguments:
# n should be set to the number of rows in A
# m should be set to the number of columns in B.

import os
import sys

# number of rows in A
n = int(sys.argv[1])

# number of columns in B
m = int(sys.argv[2])


# are we reading an A or a B file?
READING_A = False
READING_B = False

# Hadoop may break each input file into several small chunks for processing
# and the streaming mode only shows us one row (line of text) at a time.
#
# If we want to know what file the input data is coming from, this is
# stored in the environment variable `mapreduce_map_input_file`:
if 'A' in os.environ['mapreduce_map_input_file']:
	READING_A = True
elif 'B' in os.environ['mapreduce_map_input_file']:
	READING_B = True
else:
	raise RuntimeError('Could not determine input file!')


# input comes from STDIN (stream data that goes to the program)
for line in sys.stdin:

    # Remove leading and trailing whitespace
	line = line.strip()

    # Split line into array of entry data
	entry = line.split(",")

    # Set row, column, and value for this entry
	row = int(entry[0])
	col = int(entry[1])
	value = float(entry[2])

    # If this is an entry in matrix A...
	if READING_A:
        # Generate the necessary key-value pairs
	
        ###
        # (your code goes here)
        ###
		for k in range(int(m)):
			print('{}\t{}'.format((row,k),(0,col,value)))
    # Otherwise, if this is an entry in matrix B...
	else:
        # Generate the necessary key-value pairs

        ###
        # (your code goes here)
        ###
		for i in range(int(n)):
	    		print('{}\t{}'.format((i,col),(1,row,value)))

