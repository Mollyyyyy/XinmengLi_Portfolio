#!/usr/bin/env python
#Reduce function for computing matrix multiply A*B

import sys

# Create data structures to hold the current row/column values 

#current_A = []
#current_T = []
key = None
current_key = None
max_year = 0
num_term = 0
avg_dur = [0,0.0]
num_track = 0
# input comes from STDIN (stream data that goes to the program)
for line in sys.stdin:

    # Remove leading and trailing whitespace
	line = line.strip()

    # Get key/value and split by tab
	key, value = line.split('\t', 1)
#	print("read key",key,"value",value) 
	value = value[1:-1].strip().split(', ')
    # Parse key/value input (your code goes here)
	value[0] = int(value[0])
	if value[0] == 1:
		value[1] = int(value[1])
		value[2] = float(value[2])
    # If we are still on the same key...
	if key == current_key:

        # Process key/value pair
		if value[0] == 0:
			num_term = num_term+1
			#current_A.append(value[1])
		else:
			num_track = num_track+1
			max_year = max(max_year,value[1])
			avg_dur[0] = avg_dur[0]+1
			avg_dur[1] = avg_dur[1]+value[2]
			#current_T.append((value[1],value[2]))
    # Otherwise, if this is a new key...
	else:
        #i If this is a new key and nmt(current_word, current_count))t the first key we've seen
		if current_key:
			if avg_dur[0]>0:
				avg_dur[1] = avg_dur[1]/avg_dur[0]
			if num_track>0:
				print((current_key,max_year,avg_dur[1],num_term*num_track))
			#if len(current_A) >0 and len(current_T)>0:
				#for t in current_T:
				#	for a in current_A:
				#		print((current_key,t[0],t[1],a))
			#elif len(current_A)==0 and len(current_T) > 0: 
				#for t in current_T:
				#	print((current_key,t[0],t[1],'NaN'))
		
        # Process input for new key
		current_key = key
		#current_A = []
		#current_T = []
#		print(value)
		num_track = 0
		max_year = 0
		num_term = 0
		avg_dur[0] = 0
		avg_dur[1] = 0
		if value[0] == 0:
			num_term = num_term+1
			#current_A.append(value[1])
		#	current_A.append(value[1])
		else:
		#	current_T.append((value[1],value[2]))   
			num_track = num_track+1
			max_year = max(max_year,value[1])
			avg_dur[0] = avg_dur[0]+1
			avg_dur[1] = avg_dur[1]+value[2]
#Compute/output result for the last key
if current_key:
	if avg_dur[0]>0:
		avg_dur[1] = avg_dur[1]/avg_dur[0]
	if num_track>0:
		print((current_key,max_year,avg_dur[1],num_term*num_track))
			#if len(current_A) >0 and len(current_T)>0:
#	if len(current_A)>0 and len(current_T)>0:
#		for t in current_T:
#			for a in current_A:
#				print((current_key,t[0],t[1],a))
#	elif len(current_A) == 0 and len(current_T)>0:
#		for t in current_T:
#			print((current_key,t[0],t[1],'NaN'))

		
