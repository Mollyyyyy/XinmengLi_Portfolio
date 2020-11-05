/*
This is the lab assignment for FALL19 CSCI-GA 2621-001 Distributed Systems 
http://www.news.cs.nyu.edu/~jinyang/fa19-ds/labs/. 
For lab2, lab3, and lab4 part A, if you run all of the test cases 100 times, the pass rate should be 100/100.
Written By Xinmeng Li.
PLEASE DO NOT REPLICATE OR FURTHER DISTRIBUTE THE CODE.
*/
package raftkv

import (
	"encoding/gob"
	"labrpc"
	"log"
	"raft"
	"sync"
	"time"
//	"fmt"
)

const Debug = 0

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug > 0 {
		log.Printf(format, a...)
	}
	return
}


type Op struct {
	Command		int
	K 		string
	V 		string
	Client  int64
	Opnum   int
	// Your definitions here.
	// Field names must start with capital letters,
	// otherwise RPC will break.
}

type resultk struct {
  	Client int64
  	Opnum  int
} 
//type resultv struct{
//	Term int
//	Index int
//}
type RaftKV struct {
	mu      sync.Mutex
	me      int
	rf      *raft.Raft
	applyCh chan raft.ApplyMsg
	killchan chan bool
	resultmap map[resultk]chan int
	maxraftstate int // snapshot if log grows this big
	m        map[string]string
	latestOp map[int64]int
	latestInd int
	// Your definitions here.
}


func (kv *RaftKV) Get(args *GetArgs, reply *GetReply) {
	// Your code here. 
	term, isleader:=kv.rf.GetState()
	if isleader == false{
		reply.WrongLeader = true
		return
	}
	//fmt.Printf("GET: kvserver %v is leader\n",kv.me)
	reply.WrongLeader = false
//	fmt.Printf("GET: kvserver %v start lock1\n",kv.me)
	kv.mu.Lock()
//	fmt.Printf("GET: kvserver %v locking1\n",kv.me)
	kv.resultmap[resultk{Client:args.Client,Opnum:args.Opnum}] = make(chan int)
	kv.mu.Unlock()
//	fmt.Printf("GET: kvserver %v Unlock1\n",kv.me)
	Index,Term,isLeader:=kv.rf.Start(Op{Command:0,K:args.Key,Client:args.Client,Opnum:args.Opnum})
	if isLeader == false{
		reply.WrongLeader = true
		return
	}
	i:=0
	for{
		i++
		//fmt.Printf("GET: kvserver %v start lock2\n",kv.me)
		kv.mu.Lock()
		//fmt.Printf("GET: kvserver %v locking2\n",kv.me)
		readchan := kv.resultmap[resultk{Client:args.Client,Opnum:args.Opnum}]
		kv.mu.Unlock()
		//fmt.Printf("GET: kvserver %v Unlock2\n",kv.me)
		select{
		case ind:=<-readchan:
		//	fmt.Printf("GET:  kvserver %v index %v replied to client %v op %v \n",kv.me,ind,args.Client,args.Opnum)
			if ind != Index{
			//	fmt.Printf("GET: kvserver %v has changed leadership -- index diff %v\n",kv.me,ind-Index)
				reply.WrongLeader = true
				return
			}
			kv.mu.Lock()
			v, ok := kv.m[args.Key]
			kv.mu.Unlock()
			if ok {
				reply.Err = OK
				reply.Value = v
			//	fmt.Printf("GET: kvserver %v responds OK key %v Value %v\n",kv.me,args.Key,v)
			}else{
			//	fmt.Printf("GET: kvserver %v responds ErrNoKey key %v\n",kv.me,args.Key)
				reply.Err = ErrNoKey
			}
			return
		case <- time.After(10 * time.Millisecond):
			if i%50==0{
			//	fmt.Printf("GET: kvserver %v checking leadership %v times key %v \n",kv.me,i,args.Key)
			}
			term, isleader=kv.rf.GetState()
			if isleader == false || Term!=term{
			//	fmt.Printf("GET:  kvserver %v has changed leadership --isleader %v term diff %v\n",kv.me,isleader,term-Term)
				reply.WrongLeader = true
				return
			}
			if i==3000{
			//	fmt.Printf("GET: kvserver %v timeout\n",kv.me)
				reply.WrongLeader = true
				return
			}
		}
		//time.Sleep(10 * time.Millisecond)
	}

}

func (kv *RaftKV) PutAppend(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	op :=1
	if args.Op == "Append"{op=2}
	term, isleader:=kv.rf.GetState()
	if isleader == false{
		reply.WrongLeader = true
		return
	}
//	fmt.Printf("PUTAPPEND: kvserver %v is leader\n",kv.me)
	reply.WrongLeader = false
	kv.mu.Lock()
	kv.resultmap[resultk{Client:args.Client,Opnum:args.Opnum}] = make(chan int)
	kv.mu.Unlock()
	Index,Term,isLeader:=kv.rf.Start(Op{Command:op,K:args.Key,V:args.Value,Client:args.Client,Opnum:args.Opnum})
	if isLeader == false{
		reply.WrongLeader = true
		return
	}
	i:=0
	for{
		i++
		kv.mu.Lock()
		readchan := kv.resultmap[resultk{Client:args.Client,Opnum:args.Opnum}]
		kv.mu.Unlock()
		select{
		case ind:=<-readchan:
			//fmt.Printf("PUTAPPEND:  kvserver %v index %v replied to client %v op %v \n",kv.me,ind,args.Client,args.Opnum)
			if ind != Index{
				//fmt.Printf("PUTAPPEND:  kvserver %v has changed leadership -- index diff %v\n",kv.me,ind-Index)
				reply.WrongLeader = true
				return
			}
			reply.Err = OK
			//fmt.Printf("PUTAPPEND: kvserver %v succeed\n",kv.me)
			return
		case <- time.After(10 * time.Millisecond):
			if i%50==0{
				//fmt.Printf(": PUTAPPEND kvserver %v checking leadership %v times key %v \n",kv.me,i,args.Key)
			}
			term, isleader=kv.rf.GetState()
			if isleader == false || Term!=term{
			//	fmt.Printf("PUTAPPEND:  kvserver %v has changed leadership --isleader %v term diff %v\n",kv.me,isleader,term-Term)
				reply.WrongLeader = true
				return
			}
		}
		//time.Sleep(10 * time.Millisecond)
	}

}

//
// the tester calls Kill() when a RaftKV instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (kv *RaftKV) Kill() {
	kv.rf.Kill()
	kv.killchan<-true
	// Your code here, if desired.
}

//
// servers[] contains the ports of the set of
// servers that will cooperate via Raft to
// form the fault-tolerant key/value service.
// me is the index of the current server in servers[].
// the k/v server should store snapshots with persister.SaveSnapshot(),
// and Raft should save its state (including log) with persister.SaveRaftState().
// the k/v server should snapshot when Raft's saved state exceeds maxraftstate bytes,
// in order to allow Raft to garbage-collect its log. if maxraftstate is -1,
// you don't need to snapshot.
// StartKVServer() must return quickly, so it should start goroutines
// for any long-running work.
//
func StartKVServer(servers []*labrpc.ClientEnd, me int, persister *raft.Persister, maxraftstate int) *RaftKV {
	// call gob.Register on structures you want
	// Go's RPC library to marshall/unmarshall.
	gob.Register(Op{})

	kv := new(RaftKV)
	kv.me = me
	kv.maxraftstate = maxraftstate

	// You may need initialization code here.

	kv.applyCh = make(chan raft.ApplyMsg)
	kv.killchan = make(chan bool)
	kv.m = make(map[string]string)
	kv.latestOp = make(map[int64]int)
	kv.latestInd = 0
	kv.resultmap = make(map[resultk]chan int)
	kv.rf = raft.Make(servers, me, persister, kv.applyCh)

	// You may need initialization code here.
	go func(){
		for{
			select{
			case <-kv.killchan:
				return
			case msg:=<-kv.applyCh:
				op := msg.Command.(Op)
				ind:= msg.Index
			//	fmt.Printf("kv server %v received index %v from apply chan\n",kv.me,ind)
				kv.mu.Lock()
				latest:=kv.latestOp[op.Client]
				kv.mu.Unlock()
				for ind > kv.latestInd+1{
			//		fmt.Printf("kv server %v with index %v wait for index %v\n",kv.me,ind,kv.latestInd+1)
					time.Sleep(10 * time.Millisecond)
				}
				kv.latestInd++
				if latest < op.Opnum{
					if op.Command ==1{
			//			fmt.Printf("PUT: kvserver %v applied key %v value %v\n",kv.me,op.K,op.V)
						kv.mu.Lock()
						kv.m[op.K]  = op.V
						kv.mu.Unlock()
					}else if op.Command ==2{
			//			fmt.Printf("Append: kvserver %v applied key %v value %v\n",kv.me,op.K,op.V)
						kv.mu.Lock()
						kv.m[op.K]  += op.V
						kv.mu.Unlock()
					}
					kv.mu.Lock()
					kv.latestOp[op.Client] = op.Opnum
					kv.mu.Unlock()
				}
				_, isleader:=kv.rf.GetState()
				if isleader == true{
					kv.mu.Lock()
					resultchan, ok := kv.resultmap[resultk{Client:op.Client,Opnum:op.Opnum}]
					kv.mu.Unlock()
				//	fmt.Printf("Main: kvserver %v check chan for client %v op %v exists %v\n",kv.me,op.Client,op.Opnum,ok)
					if ok {
						go func(){
						//	fmt.Printf("Main: kvserver %v put index %v to chan client %v op %v\n",kv.me,ind,op.Client,op.Opnum)
							resultchan<-ind
							return
						}()
					}
				}
			default:
			}
		}
	}()
	return kv
}
