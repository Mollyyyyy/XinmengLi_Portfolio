#!/usr/bin/env python
#Reduce function for computing matrix multiply A*B

import sys

# Create data structures to hold the current row/column values 

###
# (if needed; your code goes here)
###

current_A = []
current_T = []
key = None
current_key = None
	
#def join(artist_id,curr_A,curr_T):
#	if artist_id == 'AR5DFHO1187B9A3CC4':
#		print(curr_A)
#		print(curr_T)
#	if len(curr_A)==0 or len(curr_T)==0:
#		return
#	for t in current_T:
#		for a in current_A:
	#		print((artist_id,t,a))

# input comes from STDIN (stream data that goes to the program)
for line in sys.stdin:

    # Remove leading and trailing whitespace
	line = line.strip()

    # Get key/value and split by tab
	key, value = line.split('\t', 1)
#	print("read key",key,"value",value) 
	value = value[1:-1].strip().split(', ')
    # Parse key/value input (your code goes here)
#	line = line.split('$')
#	key = line[0]
#	value = [line[1],line[2]]
#	if key == 'AR5DFHO1187B9A3CC4':
#		print(key,value)
    # If we are still on the same key...
	if key == current_key:

        # Process key/value pair
		if int(value[0]) == 0:
			current_A.append(value[1])
		else:
			current_T.append(value[1])
    # Otherwise, if this is a new key...
	else:
        #i If this is a new key and nmt(current_word, current_count))t the first key we've seen
		if current_key:
			#if current_key == 'AR5DFHO1187B9A3CC4':
			#	print(current_A,current_T)
			if len(current_A) >0 and len(current_T)>0:
				for t in current_T:
					for a in current_A:
						print('(',current_key,',',t,',',a,')')
		
        # Process input for new key
		current_key = key
		current_A = []
		current_T = []
#		print(value)
		if int(value[0]) == 0:
			current_A.append(value[1])
		else:
			current_T.append(value[1])   
	#if current_key == 'AR5DFHO1187B9A3CC4':
	#	print(current_A)
	#	print(current_T)

#Compute/output result for the last key
if current_key:
	if len(current_A)>0 and len(current_T)>0:
		for t in current_T:
			for a in current_A:
				print('(',current_key,',',t,',',a,')')
