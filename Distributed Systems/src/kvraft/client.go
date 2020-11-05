/*
This is the lab assignment for FALL19 CSCI-GA 2621-001 Distributed Systems 
http://www.news.cs.nyu.edu/~jinyang/fa19-ds/labs/. 
For lab2, lab3, and lab4 part A, if you run all of the test cases 100 times, the pass rate should be 100/100.
Written By Xinmeng Li.
PLEASE DO NOT REPLICATE OR FURTHER DISTRIBUTE THE CODE.
*/
package raftkv

import "labrpc"
import "crypto/rand"
import "math/big"
//import "fmt"
import "time"
type Clerk struct {
	servers []*labrpc.ClientEnd
	lastleader int
	me         int64
	opnum      int
	// You will have to modify this struct.
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(servers []*labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.servers = servers
	ck.lastleader = 0
	ck.me = nrand()
	ck.opnum = 0
	// You'll have to add code here.
	return ck
}

//
// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) Get(key string) string {
	// You will have to modify this function.
	ck.opnum++
	server:=ck.lastleader
	for{
		var reply GetReply
		args:=GetArgs{Key:key,Client:ck.me,Opnum:ck.opnum}
		//fmt.Printf("GET: client %v sends req to server %v key %v Opnum %v\n",ck.me,server,key,args.Opnum)
		ok := ck.servers[server].Call("RaftKV.Get", &args, &reply)
		if !ok || reply.WrongLeader == true{
			server = (server+1)%len(ck.servers)
		}else if reply.Err == ErrNoKey{
			ck.lastleader = server
			return ""
		}else{
			ck.lastleader = server
			return reply.Value
		}
		time.Sleep(100 * time.Millisecond)
	}
}

//
// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.servers[i].Call("RaftKV.PutAppend", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
//
func (ck *Clerk) PutAppend(key string, value string, op string) {
	// You will have to modify this function.
	ck.opnum++
	server:=ck.lastleader
	for{
		args:=PutAppendArgs{Key:key,Value:value,Op:op,Client:ck.me,Opnum:ck.opnum}
		//fmt.Printf("%v: client %v sends req to server %v key %v value %v Opnum %v\n",op,ck.me,server,key,value,args.Opnum)
		var reply PutAppendReply
		ok := ck.servers[server].Call("RaftKV.PutAppend", &args, &reply)
		if !ok || reply.WrongLeader == true ||reply.Err!=OK{
			server = (server+1)%len(ck.servers)
		}else{
			ck.lastleader = server
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}
func (ck *Clerk) Append(key string, value string) {
	ck.PutAppend(key, value, "Append")
}
