/*
This is the lab assignment for FALL19 CSCI-GA 2621-001 Distributed Systems 
http://www.news.cs.nyu.edu/~jinyang/fa19-ds/labs/. 
For lab2, lab3, and lab4 part A, if you run all of the test cases 100 times, the pass rate should be 100/100.
Written By Xinmeng Li.
PLEASE DO NOT REPLICATE OR FURTHER DISTRIBUTE THE CODE.
*/
package mapreduce

import (
	"sort"
	"os"
	"encoding/json"
	"fmt"
)
// doReduce manages one reduce task: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
	//
	// You will need to write this function.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTaskNumber) yields the file
	// name from map task m.
	//
	files := make([] *os.File, nMap) 
	for m:=0;m<nMap;m++{
		filename:=reduceName(jobName, m, reduceTaskNumber)
		file, err := os.Open(filename) // For read access.
		if err != nil {
			fmt.Printf("ERR open file %v\n",err)
		}else{files[m] = file}
	}
	keyvalues:=make(map[string][]string)
	for m:=0;m<nMap;m++{
		dec := json.NewDecoder(files[m])
		for{
			var kv KeyValue
			if err := dec.Decode(&kv); err != nil {
	        	//fmt.Printf("ERR decode file %v\n",err)
	        	break
	        }
	        keyvalues[kv.Key]=append(keyvalues[kv.Key],kv.Value)
	    }
	    if err := files[m].Close(); err != nil {
	        fmt.Printf("ERR close file %v\n",err)
	    }
	}
    keys := make([]string, 0, len(keyvalues))
	for k := range keyvalues {
		keys = append(keys, k)
	}
	sort.Strings(keys)
    f, err := os.Create(outFile)
    if err != nil {
        fmt.Printf("ERR open file %v\n",err)
    }
	enc := json.NewEncoder(f)
	var count string
	for _, k := range keys {
		count=reduceF(k, keyvalues[k])
		enc.Encode(KeyValue{k,count})
		//fmt.Printf("length of values %v, reduced result %v \n",len(keyvalues[k]),count)
	}
	if err := f.Close(); err != nil {
        fmt.Printf("ERR close file %v\n",err)
    }
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
}
