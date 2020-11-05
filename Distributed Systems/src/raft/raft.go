/*
This is the lab assignment for FALL19 CSCI-GA 2621-001 Distributed Systems 
http://www.news.cs.nyu.edu/~jinyang/fa19-ds/labs/. 
For lab2, lab3, and lab4 part A, if you run all of the test cases 100 times, the pass rate should be 100/100.
Written By Xinmeng Li.
PLEASE DO NOT REPLICATE OR FURTHER DISTRIBUTE THE CODE.
*/
package raft

//
// this is an outline of the API that raft must expose to
// the service (or tester). see comments below for
// each of these functions for more details.
//
// rf = Make(...)
//   create a new Raft server.
// rf.Start(command interface{}) (index, term, isleader)
//   start agreement on a new log entry
// rf.GetState() (term, isLeader)
//   ask a Raft for its current term, and whether it thinks it is leader
// ApplyMsg
//   each time a new entry is committed to the log, each Raft peer
//   should send an ApplyMsg to the service (or tester)
//   in the same server.
//

import "sync"
import "labrpc"
//import "fmt"
import "math/rand"
import "time"
import "bytes"
import "encoding/gob"



//
// as each Raft peer becomes aware that successive log entries are
// committed, the peer should send an ApplyMsg to the service (or
// tester) on the same server, via the applyCh passed to Make().
//
type ApplyMsg struct {
	Index       int
	Command     interface{}
	UseSnapshot bool   // ignore for lab2; only used in lab3
	Snapshot    []byte // ignore for lab2; only used in lab3
}

type LogEntry struct{
	Command interface{}
	Term int
}

//
// A Go object implementing a single Raft peer.
//
type Raft struct {
	mu        sync.Mutex          // Lock to protect shared access to this peer's state
	peers     []*labrpc.ClientEnd // RPC end points of all peers
	persister *Persister          // Object to hold this peer's persisted state
	me        int                 // this peer's index into peers[]
	CurrentTerm int // Persist state on all 
	VotedFor	int // Persist state on all 
	Log	[]LogEntry // Persist state on all 
	commitIndex int //NonPersist state on all
	lastApplied int //NonPersist state on all
	nextIndex []int //NonPersist state on leader
	matchIndex []int //NonPersist state on leader
	role	int //0 for followers, 1 for candidate, 2 for leader
	elapsetime int
	votenum int
	electchan chan bool
	appchan chan bool
	grantvote chan bool
	applyCh chan ApplyMsg
	//killed bool
	killchan chan bool
	// Your data here (2A, 2B, 2C).
	// Look at the paper's Figure 2 for a description of what
	// state a Raft server must maintain.

}

// return CurrentTerm and whether this server
// believes it is the leader.
func (rf *Raft) GetState() (int, bool) {
	var term int
	var isleader bool
	// Your code here (2A).
	term = rf.currTerm()
	if rf.getRole() ==2{
		isleader = true
	}else{
		isleader = false
	}
	return term, isleader
}

//
// save Raft's persistent state to stable storage,
// where it can later be retrieved after a crash and restart.
// see paper's Figure 2 for a description of what should be persistent.
//
func (rf *Raft) persist() {
	// Your code here (2C).
	// Example:
	 w := new(bytes.Buffer)
	 e := gob.NewEncoder(w)
	 rf.mu.Lock()
	 e.Encode(rf.CurrentTerm)
	 e.Encode(rf.VotedFor)
	 e.Encode(rf.Log)
	 rf.mu.Unlock()
	 data := w.Bytes()
	 rf.persister.SaveRaftState(data)
}

//
// restore previously persisted state.
//
func (rf *Raft) readPersist(data []byte) {
	// Your code here (2C).
	// Example:
	if data == nil || len(data) < 1 { // bootstrap without any state?
		return
	}
	r := bytes.NewBuffer(data)
	d := gob.NewDecoder(r)
	rf.mu.Lock()
	d.Decode(&rf.CurrentTerm)
	d.Decode(&rf.VotedFor)
	d.Decode(&rf.Log)
	rf.mu.Unlock()
}




//
// example RequestVote RPC arguments structure.
// field names must start with capital letters!
//
type RequestVoteArgs struct {
	// Your data here (2A, 2B).
	Term int
	CandidateID int
	LastLogIndex int
	LastLogTerm int
}

//
// example RequestVote RPC reply structure.
// field names must start with capital letters!
//
type RequestVoteReply struct {
	// Your data here (2A).
	Term int
	VoteGranted bool
}
type AppendEntriesArgs struct{
	Term int
	LeaderID int
	PrevLogIndex int
	PrevLogTerm int
	Entries	[]LogEntry
	LeaderCommit int
}
type AppendEntriesReply struct{
	Term int
	Success bool
	Inconsistent bool
	TermMatchi int
}
//
// example RequestVote RPC handler.
//
func (rf *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {
	// Your code here (2A, 2B).
	//rf.reqchan <- args
	//rf.mu.Lock()
	//defer rf.mu.Unlock()
	/*if rf.checkKill() == true{
		return
	}*/
	rf.toFoll(args.Term)
	//fmt.Printf("Server %v received a request vote\n",rf.me)
	reply.Term = rf.currTerm()
	if args.Term<rf.currTerm(){
		reply.VoteGranted = false
		return
	}else{
		rf.mu.Lock()
		v:=rf.VotedFor
		rf.mu.Unlock()
		if v == -1||v ==args.CandidateID{
			l:=rf.loglen()
			if l==0{
				reply.VoteGranted = true
			}else if args.LastLogTerm>rf.LogTerm(l-1){
				reply.VoteGranted = true
			}else if args.LastLogTerm==rf.LogTerm(l-1) && args.LastLogIndex>=(l-1){
				reply.VoteGranted = true
			}else{reply.VoteGranted = false}
		}else{reply.VoteGranted = false}
	}
	if reply.VoteGranted == true /*&& rf.checkKill() == false*/{
		rf.mu.Lock()
		rf.VotedFor = args.CandidateID
		rf.mu.Unlock()
		rf.persist()
		rf.grantvote<-true
	}
}

func (rf *Raft) AppendEntries(args *AppendEntriesArgs, reply *AppendEntriesReply){
	/*if rf.checkKill() == true{
		reply.Success = false
		return
	}*/
	//fmt.Printf("Server %v received heartbeat from %v\n",rf.me,args.LeaderID)
	r:=rf.getRole()
	if r== 0 /*&& rf.checkKill() == false*/{
		rf.appchan <-true
	}else if r == 1 && args.Term>=rf.currTerm() /*&& rf.checkKill() == false*/{
		//fmt.Printf("Candidate %v received append entries and becomes follower\n",rf.me)
		rf.appchan <-true
		rf.mu.Lock()
		rf.role = 0
		rf.VotedFor = -1
		rf.mu.Unlock()
		rf.persist()
	}
	rf.toFoll(args.Term)
	reply.Term = rf.currTerm()
	reply.Inconsistent = false
	l:=rf.loglen()
	if args.Term<rf.currTerm(){
		reply.Success = false
		return
	}else if l-1<args.PrevLogIndex{
		reply.Inconsistent = true
		reply.Success = false
		reply.TermMatchi = rf.loglen()
		//fmt.Printf("Prev log entry of server %v missing, has prev %v, change next index from %v to %v\n",rf.me,rf.Log[rf.loglen()-1],args.PrevLogIndex+1,reply.TermMatchi)
		return
	}else if args.PrevLogIndex>0 && rf.LogTerm(args.PrevLogIndex) != args.PrevLogTerm{
		reply.Inconsistent = true
		reply.Success = false
		rf.mu.Lock()
		con:=rf.Log[args.PrevLogIndex]
		rf.Log = rf.Log[:args.PrevLogIndex]
		rf.mu.Unlock()
		rf.persist()
		termMatchi:=rf.loglen()-1
		ll:=rf.loglen()-1
		for i:=ll;i>=0;i--{
			if con.Term != rf.LogTerm(i){
				rf.mu.Lock()
				termMatchi = i+1
				rf.mu.Unlock()
				break
			}
			if i == 0{
				rf.mu.Lock()
				termMatchi = 1
				rf.mu.Unlock()
			}
		}
		reply.TermMatchi = termMatchi
		//fmt.Printf("Prev log entry conflits of server %v, change nextindex from %v prev %v to %v prev %v\n",rf.me,args.PrevLogIndex+1,con,reply.TermMatchi,rf.Log[reply.TermMatchi-1])
		return
	}else{
		reply.Success = true
		/*
		rf.mu.Lock()
		rf.Log = rf.Log[:args.PrevLogIndex+1]
		rf.mu.Unlock()
		rf.persist()*/
		//fmt.Printf("Server %v Append entries %v of leader\n",rf.me,args.Entries)
		for i:=0;i<len(args.Entries);i++{
			if args.PrevLogIndex+1+i >= rf.loglen(){
				rf.mu.Lock()
				rf.Log = append(rf.Log,args.Entries[i])
				rf.mu.Unlock()
			}else if rf.LogTerm(args.PrevLogIndex+1+i) == args.Entries[i].Term{
				rf.mu.Lock()
				rf.Log[args.PrevLogIndex+1+i] = args.Entries[i]
				rf.mu.Unlock()
			}else{
					rf.mu.Lock()
					rf.Log = rf.Log[:args.PrevLogIndex+1+i]
					rf.Log = append(rf.Log,args.Entries[i])
					rf.mu.Unlock()
					//fmt.Printf("__________term conflicts\n")
			}
		}
		rf.persist()
		if args.LeaderCommit>rf.getCommit(){
			cI:= rf.loglen()-1
			if args.LeaderCommit<cI{
				cI = args.LeaderCommit
			}
			rf.mu.Lock()
			rf.commitIndex = cI
			rf.mu.Unlock()
			applied:=rf.getApplied()
			committed:=rf.getCommit()
			if committed>applied{
				go func(){
					//applied=rf.getApplied()+1
					//committed=rf.getCommit()
					for i:=rf.getApplied()+1;i<=rf.getCommit();i++{
						if i<rf.loglen(){
							rf.mu.Lock()
							appMsg:= ApplyMsg{Index:i,Command:rf.Log[i].Command}
							rf.mu.Unlock()
						/*if rf.checkKill() == false{*/
							rf.applyCh <-appMsg
						}else{
							//fmt.Printf("WWWWWWWWWWWWWWWWWWWWWWARNING: loglen is %v committed is %v, applied is %v\n",rf.loglen(),committed,applied)
						}
					}
					rf.mu.Lock()
					rf.lastApplied = rf.commitIndex
					rf.mu.Unlock()
					//fmt.Printf("Follower %v lastApplied %v committedIndex %v\n",rf.me,rf.getApplied(),rf.getCommit())
				}()
			}
		}
	}
}
//
// example code to send a RequestVote RPC to a server.
// server is the index of the target server in rf.peers[].
// expects RPC arguments in args.
// fills in *reply with RPC reply, so caller should
// pass &reply.
// the types of the args and reply passed to Call() must be
// the same as the types of the arguments declared in the
// handler function (including whether they are pointers).
//
// The labrpc package simulates a lossy network, in which servers
// may be unreachable, and in which requests and replies may be lost.
// Call() sends a request and waits for a reply. If a reply arrives
// within a timeout interval, Call() returns true; otherwise
// Call() returns false. Thus Call() may not return for a while.
// A false return can be caused by a dead server, a live server that
// can't be reached, a lost request, or a lost reply.
//
// Call() is guaranteed to return (perhaps after a delay) *except* if the
// handler function on the server side does not return.  Thus there
// is no need to implement your own timeouts around Call().
//
// look at the comments in ../labrpc/labrpc.go for more details.
//
// if you're having trouble getting RPC to work, check that you've
// capitalized all field names in structs passed over RPC, and
// that the caller passes the address of the reply struct with &, not
// the struct itself.
//
func (rf *Raft) sendRequestVote(server int, args *RequestVoteArgs, reply *RequestVoteReply) bool {
	//rf.mu.Lock()
	//defer rf.mu.Unlock()
	ok := rf.peers[server].Call("Raft.RequestVote", args, reply)
	return ok
}

func (rf *Raft) sendAppendEntries(server int, args *AppendEntriesArgs, reply *AppendEntriesReply) bool{
	//rf.mu.Lock()
	//defer rf.mu.Lock()
//	fmt.Printf("---Leader %v send heartbeat to server %v\n",args.LeaderID,server)
	ok:= rf.peers[server].Call("Raft.AppendEntries",args,reply)
	return ok
}
//
// the service using Raft (e.g. a k/v server) wants to start
// agreement on the next command to be appended to Raft's log. if this
// server isn't the leader, returns false. otherwise start the
// agreement and return immediately. there is no guarantee that this
// command will ever be committed to the Raft log, since the leader
// may fail or lose an election.
//
// the first return value is the index that the command will appear at
// if it's ever committed. the second return value is the current
// term. the third return value is true if this server believes it is
// the leader.
//
func (rf *Raft) Start(command interface{}) (int, int, bool) {
	rf.mu.Lock()
	index := -1
	term := -1
	isLeader := true
	rf.mu.Unlock()
	// Your code here (2B).
	if rf.getRole()!=2{return index,term,false}
	rf.mu.Lock()
	index = len(rf.Log)
    term = rf.CurrentTerm
	rf.Log = append(rf.Log,LogEntry{Term:rf.CurrentTerm,Command:command})
	rf.mu.Unlock()
	rf.persist()
	//l:=rf.loglen()
	//fmt.Printf("________________________________________________________new command is sent to leader %v, current log %v %v term %v\n",rf.me,l,rf.Log[l-1],rf.currTerm())
	return index, term, isLeader
}

//
// the tester calls Kill() when a Raft instance won't
// be needed again. you are not required to do anything
// in Kill(), but it might be convenient to (for example)
// turn off debug output from this instance.
//
func (rf *Raft) Kill() {
	// Your code here, if desired.
	rf.killchan<-true
	//fmt.Printf("___________________Kill is called. %v should be disconnected\n",rf.me)
	/*go func(){
		for{
			fmt.Printf("___________________Kill is called. %v should be disconnected\n",rf.me)
			select{
			case <-rf.electchan:
				fmt.Printf("Server %v electchan is not empty\n",rf.me)
			case <-rf.grantvote:
				fmt.Printf("Server %v grantchan is not empty\n",rf.me)
			case <-rf.appchan:
				fmt.Printf("Server %v appchan is not empty\n",rf.me)
			case <-time.After(time.Duration(60)*time.Millisecond):
				fmt.Printf("Server %v chan all empty\n",rf.me)
				return
			}
		}
		return
	}()*/
	//rf.electchan = make(chan bool)
	//rf.appchan = make(chan bool)
	//rf.grantvote = make(chan bool)
	//close(rf.electchan)
	//close(rf.appchan)
	//close(rf.grantvote)
}

func(rf *Raft) toCand(){
	//rf.mu.Lock()
	//defer rf.mu.Unlock()
	rf.mu.Lock()
	rf.role = 1
	rf.CurrentTerm++
	//fmt.Printf("Candidate %v now has term %v, log %v\n",rf.me,rf.CurrentTerm,len(rf.Log))
	rf.VotedFor = rf.me
	rf.votenum = 1
	rf.mu.Unlock()
	rf.persist()
	args:=&RequestVoteArgs{
		Term : rf.currTerm(),
		CandidateID : rf.me,
		LastLogIndex :rf.loglen()-1,
		LastLogTerm : rf.Log[rf.loglen()-1].Term,
	}
	majority := len(rf.peers)/2+1
	//fmt.Printf("Candidate %v starts to request vote\n",rf.me)
	for i:=0;i<len(rf.peers);i++{
		if rf.getRole()==0{return}
		if i==rf.me{continue}
		var reply RequestVoteReply
		go func(server int){
			if rf.getRole()==0 /*|| rf.checkKill()==true*/{return}
			ok:= rf.sendRequestVote(server,args,&reply)
			if ok{rf.toFoll(reply.Term)}
			if rf.getRole() == 0/*||rf.checkKill()==true*/{return}
			if ok && rf.getRole()==1 && args.Term == rf.currTerm()/* && rf.checkKill()==false*/{
				if reply.VoteGranted == true{
					//fmt.Printf("Candidate %v get a vote from %v\n",rf.me,server)
					rf.mu.Lock()
					rf.votenum++
					v:=rf.votenum
					rf.mu.Unlock()
					if v==majority{
						//if rf.checkKill()==false{
						rf.mu.Lock()
						rf.role=2
						rf.mu.Unlock()
						   // fmt.Printf("chan: Send to elect chan of Candidate %v\n",rf.me)
						rf.electchan<-true

					}
				//	fmt.Printf("Candidate %v now has %v vote\n",rf.me,rf.votenum)
				}
			}
		}(i)
	}
}
func(rf *Raft) toLead(){
	/*if rf.checkKill() == true{
		fmt.Printf("Server %v Cannot be a leader, HAS BEEN KILLED\n",rf.me)
		return
	}*/
	rf.mu.Lock()
	rf.role=2
	for i:=0;i<len(rf.peers);i++{
		rf.nextIndex[i] = len(rf.Log)
	}
	rf.mu.Unlock()
	//fmt.Printf("RAFT: Candidate %v becomes Leader\n",rf.me)
}

func(rf *Raft) toFoll(term int){
	changed :=0
	if term>rf.currTerm(){
		rf.mu.Lock()
		//fmt.Printf("Server %v has term %v, should be %v\n",rf.me,rf.CurrentTerm,term)
		rf.CurrentTerm = term
		rf.VotedFor = -1
		rf.role = 0
		changed=1
		//fmt.Printf("RAFT: Server %v becomes Follower\n",rf.me)
		rf.mu.Unlock()
	}
	if changed == 1{rf.persist()}

}
func(rf *Raft) loglen() int{
	rf.mu.Lock()
	l:=len(rf.Log)
	rf.mu.Unlock()
	return l
}
func(rf *Raft) currTerm() int{
	rf.mu.Lock()
	c:= rf.CurrentTerm
	rf.mu.Unlock()
	return c
}
func(rf *Raft) LogTerm(i int) int{
	rf.mu.Lock()
	c:= rf.Log[i].Term
	rf.mu.Unlock()
	return c
}
func(rf *Raft) getRole() int{
	rf.mu.Lock()
	r:=rf.role
	rf.mu.Unlock()
	return r
}
func(rf *Raft) getCommit() int{
	rf.mu.Lock()
	r:=rf.commitIndex
	rf.mu.Unlock()
	return r
}
func(rf *Raft) getApplied() int{
	rf.mu.Lock()
	r:=rf.lastApplied
	rf.mu.Unlock()
	return r
}
func(rf *Raft) getMatched(i int) int{
	rf.mu.Lock()
	r:=rf.matchIndex[i]
	rf.mu.Unlock()
	return r
}

func(rf *Raft) heartbeat(server int,args *AppendEntriesArgs){
	if rf.getRole()!=2 /*|| rf.checkKill() == true*/{return}
	var reply AppendEntriesReply
	//fmt.Printf("Before heartbeat Leader %v has Role %v, Term %v, commitIndex %v, LastApplied %v and Log %v\n",rf.me,rf.role,rf.currTerm(),rf.commitIndex,rf.lastApplied,len(rf.Log))
	ok := rf.sendAppendEntries(server, args, &reply)
	if ok {rf.toFoll(reply.Term)}
	if rf.getRole()!=2 /*|| rf.checkKill() == true*/{return}
	if ok && rf.getRole()==2 && rf.currTerm() == args.Term{
		rf.mu.Lock()
		if reply.Success == true && args.Entries !=nil{
			rf.nextIndex[server] = 1+len(args.Entries)+args.PrevLogIndex
			rf.matchIndex[server] = len(args.Entries)+args.PrevLogIndex
		}else if reply.Inconsistent == true{
			rf.nextIndex[server] = reply.TermMatchi
			//rf.nextIndex[server]--
		}
		rf.mu.Unlock()
		ll:=rf.loglen()-1
		cc:=rf.getCommit()
		for N:=ll;N>cc;N--{
			num:=1
			for matchi:=0;matchi<len(rf.matchIndex);matchi++{
				if N<=rf.getMatched(matchi){num++}
			}
		//	fmt.Printf("Leader %v term %v has %v peers with match index >= %v where entry term %v\n",rf.me,rf.currTerm(),num,N,rf.LogTerm(N))
			if num>=len(rf.peers)/2+1 && rf.LogTerm(N)==rf.currTerm(){
			//	fmt.Printf("Leader %v starts applying lastApplied %v commitIndex %v\n",rf.me,rf.getApplied(),rf.getCommit())
				rf.mu.Lock()
				rf.commitIndex = N
				rf.mu.Unlock()
				go func(){
					//applied:=rf.getApplied()+1
					//committed:=rf.getCommit()
					for appi:=rf.getApplied()+1;appi<=rf.getCommit();appi++{
						//go func(appi int){
						if appi < rf.loglen(){
						//	fmt.Printf("Leader %v tries to apply entry at index %v commited %v, log len %v\n",rf.me,appi,rf.getCommit(),rf.loglen())
							rf.mu.Lock()
							appMsg:= ApplyMsg{Index:appi,Command:rf.Log[appi].Command}
							rf.mu.Unlock()
						/*if rf.checkKill() == false{*/
							//fmt.Printf("Leader apply %v\n",appMsg)
							rf.applyCh <-appMsg
						//	fmt.Printf("Leader %v applied at index %v\n",rf.me,appi)
						}else{
							//fmt.Printf("WWWWWWWWWWWWWWWWWWWWWWARNING: leader loglen is %v committed is %v, applied is %v\n",rf.loglen(),rf.getCommit(),rf.getApplied()+1)
						}
							//return
						//}(appi)
					}
					rf.mu.Lock()
					rf.lastApplied = rf.commitIndex
					rf.mu.Unlock()
				//	fmt.Printf("Leader %v finished applying lastApplied %v commitIndex %v\n",rf.me,rf.getApplied(),rf.getCommit())
				}()
				break
			}
		}
	}
	//fmt.Printf("After heartbeat Leader %v has Role %v, Term %v, commitIndex %v, LastApplied %v and Log %v\n",rf.me,rf.role,rf.currTerm(),rf.commitIndex,rf.lastApplied,rf.Log)
	return
}
/*func(rf *Raft) checkKill() bool{
	rf.mu.Lock()
	k:=rf.killed
	rf.mu.Unlock()
	//fmt.Printf("check kill being called and return %v\n",k)
	return k
}*/
//
// the service or tester wants to create a Raft server. the ports
// of all the Raft servers (including this one) are in peers[]. this
// server's port is peers[me]. all the servers' peers[] arrays
// have the same order. persister is a place for this server to
// save its persistent state, and also initially holds the most
// recent saved state, if any. applyCh is a channel on which the
// tester or service expects Raft to send ApplyMsg messages.
// Make() must return quickly, so it should start goroutines
// for any long-running work.
//
func Make(peers []*labrpc.ClientEnd, me int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {
	rf := &Raft{}
	rf.peers = peers
	rf.persister = persister
	rf.me = me
	rf.elapsetime = rand.Intn(200)+1000
	rf.VotedFor = -1
	rf.CurrentTerm = 0
	rf.commitIndex = 0
	rf.lastApplied = 0
	rf.role = 0
	rf.votenum=0
	rf.Log = append(rf.Log,LogEntry{Term:0,Command:nil})
	rf.nextIndex = make([]int,len(rf.peers))
	for i:=0;i<len(rf.peers);i++{ 
		rf.mu.Lock()
		rf.nextIndex[i] = 1
		rf.mu.Unlock()
	}
	rf.matchIndex = make([]int,len(rf.peers))
	rf.electchan = make(chan bool)
	rf.appchan = make(chan bool)
	rf.grantvote = make(chan bool)
	rf.killchan = make(chan bool)
	rf.applyCh = applyCh
	// initialize from state persisted before a crash
	rf.readPersist(persister.ReadRaftState())
	rf.persist()
	//fmt.Printf("Initialization: %v now has term %v, log %v\n",rf.me,rf.CurrentTerm,rf.loglen())
	// Your initialization code here (2A, 2B, 2C).
	go func(){
		for{
			/*if rf.checkKill() == true{
				fmt.Printf("goroutine for server %v is done\n",rf.me)
				return
			}*/
			//fmt.Printf("%v is %v now with term %v, log %v\n",rf.me,rf.role,rf.CurrentTerm,rf.Log)
			select{
			case <-rf.killchan:
				return
			default:
				switch rf.getRole(){
				case 0:
					select{
					case <- rf.appchan:
						//l:=rf.loglen()
					//	fmt.Printf("     HHHHHHHH      chan : Follower %v received hearbeat, now has term %v log %v\n",rf.me,rf.CurrentTerm,rf.loglen())
					case <- rf.grantvote:
						//fmt.Printf("chan: Follower %v granted vote to others\n",rf.me)
					case <- time.After(time.Duration(rf.elapsetime)*time.Millisecond):
						rf.mu.Lock()
						rf.role =1
						rf.mu.Unlock()
						//fmt.Printf("Follower %v time out and becomes Candidate\n",rf.me)
					}
					//rf.mu.Unlock()
				case 1:
			//		fmt.Printf("Candidate %v before call Cand() has term %v\n",rf.me,rf.currTerm())
					rf.toCand()
					select{
					case <- rf.electchan:
						//fmt.Printf("chan: Candidate %v received from elect win channel\n",rf.me)
						//if rf.checkKill()==false{
						rf.toLead()
					case <- rf.appchan:
						//fmt.Printf("chan: Candidate %v received from heartbeat channel\n",rf.me)
						//rf.mu.Lock()
						//rf.role = 0
						//rf.VotedFor = -1
						//rf.mu.Unlock()
			//			fmt.Printf("Candidate %v received heartbeat and becomes Follower, role is %v now\n",rf.me,rf.role)
					case <-time.After(time.Duration(rand.Intn(200)+1000)*time.Millisecond):
						//fmt.Printf("Candidate %v time out and starts new election\n",rf.me)
					}
				case 2:
					//rf.mu.Lock()
					for i:=0;i<len(rf.peers);i++{
						if rf.getRole() !=2/*||rf.checkKill()==true */{break}
						if rf.getRole() ==2 && i!=rf.me{
							rf.mu.Lock()
							args:=&AppendEntriesArgs{
								Term :rf.CurrentTerm,
								LeaderID : rf.me,
								PrevLogIndex :rf.nextIndex[i]-1,
								PrevLogTerm : rf.Log[rf.nextIndex[i]-1].Term,
								Entries : nil,
								LeaderCommit : rf.commitIndex,
							}
							if len(rf.Log)-1>=rf.nextIndex[i]{
								args.Entries = rf.Log[rf.nextIndex[i]:]
							}
							rf.mu.Unlock()
							go rf.heartbeat(i,args)
						}
					}
					time.Sleep(100*time.Millisecond)
					//rf.mu.Unlock()
				}
			}
		}

	}()

	return rf
}
