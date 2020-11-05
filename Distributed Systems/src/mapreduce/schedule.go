/*
This is the lab assignment for FALL19 CSCI-GA 2621-001 Distributed Systems 
http://www.news.cs.nyu.edu/~jinyang/fa19-ds/labs/. 
For lab2, lab3, and lab4 part A, if you run all of the test cases 100 times, the pass rate should be 100/100.
Written By Xinmeng Li.
PLEASE DO NOT REPLICATE OR FURTHER DISTRIBUTE THE CODE.
*/
package mapreduce

import "fmt"
//import "sync"

//
// schedule() starts and waits for all tasks in the given phase (Map
// or Reduce). the mapFiles argument holds the names of the files that
// are the inputs to the map phase, one per map task. nReduce is the
// number of reduce tasks. the registerChan argument yields a stream
// of registered workers; each item is the worker's RPC address,
// suitable for passing to call(). registerChan will yield all
// existing registered workers (if any) and new ones as they register.
//
func schedule(jobName string, mapFiles []string, nReduce int, phase jobPhase, registerChan chan string) {
	var ntasks int
	var n_other int // number of inputs (for reduce) or outputs (for map)
	switch phase {
	case mapPhase:
		ntasks = len(mapFiles)
		n_other = nReduce
	case reducePhase:
		ntasks = nReduce
		n_other = len(mapFiles)
	}

	fmt.Printf("Schedule: %v %v tasks (%d I/Os)\n", ntasks, phase, n_other)

	// All ntasks tasks have to be scheduled on workers, and only once all of
	// them have been completed successfully should the function return.
	// Remember that workers may fail, and that any given worker may finish
	// multiple tasks.
	//
	// TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO TODO
	//
	tasks := make(chan int)
	oks := make(chan bool)
	finished:=0
	//var wg sync.WaitGroup
	//wg.Add(1)
	go func(){
		for task:=0;task<ntasks;task++{
			//fmt.Printf("put task %v in chan %v in total\n",task,ntasks)
			tasks<-task
		}
		//wg.Done()
	}()
	//wg.Wait()
	//go func(){
loop:
	for{
		select{
		case task := <-tasks:
			go func(task int){
				wraddress:=<-registerChan
				//fmt.Printf("assign task %v to %v\n",task,wraddress)			
				arg:=&DoTaskArgs{
					Phase :phase,
					JobName : jobName,
					TaskNumber :task,
					NumOtherPhase: n_other,
				}
				if phase == mapPhase{arg.File = mapFiles[task]}
				ok := call(wraddress, "Worker.DoTask", arg, nil)
				if !ok {
					tasks<-task
					//fmt.Printf("%v task failed\n",arg.TaskNumber)
				}else{
					oks<-ok
					//fmt.Printf("finished %v\n",task)
					registerChan<-wraddress
				}
			}(task)
		case <-oks:
			finished++
			if finished == ntasks{break loop}
		}
	}
	//}()
	fmt.Printf("Schedule: %v phase done\n", phase)
}
