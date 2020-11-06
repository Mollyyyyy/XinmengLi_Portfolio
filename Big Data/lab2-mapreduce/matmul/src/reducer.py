#!/usr/bin/env python
#Reduce function for computing matrix multiply A*B

import sys

# Create data structures to hold the current row/column values 

###
# (if needed; your code goes here)
###

current_value = 0
current_A = []
current_B = []
key = None
current_key = None

#def computeCij(current_A,current_B):
#	current_A = sorted(current_A)
#	current_B = sorted(current_B)
#	Aind = 0
#	Bind = 0
#	Cij = 0
#	while Aind < len(current_A) and Bind < len(current_B):
#		k1 = current_A[Aind][0]
#		k2 = current_B[Bind][0]
#		if k1 == k2:
#			Cij = Cij+current_A[Aind][1]*current_B[Bind][1]
#			Aind = Aind +1
#			Bind = Bind +1
#		elif k1< k2:
#			Aind = Aind+1
#		else:
#			Bind = Bind+1
#	print("calculating Cij",Cij)
#	return Cij
#print("______Reducer Starts___________")
# input comes from STDIN (stream data that goes to the program)
for line in sys.stdin:

    # Remove leading and trailing whitespace
	line = line.strip()

    # Get key/value and split by tab
	key, value = line.split('\t', 1)
#	print("read key",key,"value",value) 

    # Parse key/value input (your code goes here)
#	try:
	value = value[1:-1].split(',')
	value[0] = int(value[0])
	value[1] = int(value[1])
	value[2] = float(value[2])
#	except ValueError:
#		continue


    # If we are still on the same key...
	if key == current_key:

        # Process key/value pair
		if value[0] == 0:
			current_A.append((value[1],value[2]))
		else:
			current_B.append((value[1],value[2]))
    # Otherwise, if this is a new key...
	else:
        #i If this is a new key and nmt(current_word, current_count))t the first key we've seen
		if current_key:
			curr_key = current_key[1:-1].split(',')
				
			current_A = sorted(current_A)
			current_B = sorted(current_B)
			Aind = 0
			Bind = 0
			Cij = 0
			while Aind < len(current_A) and Bind < len(current_B):
				k1 = current_A[Aind][0]
				k2 = current_B[Bind][0]
				if k1 == k2:
					Cij = Cij+current_A[Aind][1]*current_B[Bind][1]
					Aind = Aind +1
					Bind = Bind +1
				elif k1< k2:
					Aind = Aind+1
				else:
					Bind = Bind+1

			print(str(int(curr_key[0]))+','+str(int(curr_key[1]))+','+str(Cij))
        # Process input for new key
		current_key = key
		current_A = []
		current_B = []
#		print(value)
		if value[0] == 0:
			current_A.append((value[1],value[2]))
		else:
			current_B.append((value[1],value[2]))   

if current_key:
	curr_key = current_key[1:-1].split(',')
	current_A = sorted(current_A)
	current_B = sorted(current_B)
	Aind = 0
	Bind = 0
	Cij = 0
	while Aind < len(current_A) and Bind < len(current_B):
		k1 = current_A[Aind][0]
		k2 = current_B[Bind][0]
		if k1 == k2:
			Cij = Cij+current_A[Aind][1]*current_B[Bind][1]
			Aind = Aind +1
			Bind = Bind +1
		elif k1< k2:
			Aind = Aind+1
		else:
			Bind = Bind+1
	print(str(int(curr_key[0]))+','+str(int(curr_key[1]))+','+str(Cij))



#Compute/output result for the last key


