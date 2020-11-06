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


# are we reading an A or a B file?
READING_Artist = False
READING_Music = False

# Hadoop may break each input file into several small chunks for processing
# and the streaming mode only shows us one row (line of text) at a time.
#
# If we want to know what file the input data is coming from, this is
# stored in the environment variable `mapreduce_map_input_file`:
if 'artist' in os.environ['mapreduce_map_input_file']:
	READING_Artist = True
elif 'track' in os.environ['mapreduce_map_input_file']:
	READING_Music = True
else:
	raise RuntimeError('Could not determine input file!')


# input comes from STDIN (stream data that goes to the program)
for line in sys.stdin:

    # Remove leading and trailing whitespace
	line = line.strip()

    # Split line into array of entry data
	entry = line.split(",")
    # If this is an entry in matrix A...
	if READING_Artist:
        # Generate the necessary key-value pairs
		artist_id,tag = entry
		print('{}\t{}'.format(artist_id,(0,tag)))
		#print(artist_id+'$',0,'$'+tag)
    # Otherwise, if this is an entry in matrix B...
	else:
        # Generate the necessary key-value pairs
		track_id,title,album,year,duration,artist_id = entry
		year = int(year)
		if year>1990:
			print('{}\t{}'.format(artist_id,(1,track_id)))
		#print(artist_id+'$',0,'$'+tag)
	    		#print(artist_id+'$',1,'$'+track_id)

