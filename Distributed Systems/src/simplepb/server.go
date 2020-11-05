/*
This is the lab assignment for FALL19 CSCI-GA 2621-001 Distributed Systems 
http://www.news.cs.nyu.edu/~jinyang/fa19-ds/labs/. 
For lab2, lab3, and lab4 part A, if you run all of the test cases 100 times, the pass rate should be 100/100.
Written By Xinmeng Li.
PLEASE DO NOT REPLICATE OR FURTHER DISTRIBUTE THE CODE.
*/
package simplepb

//
// This is a outline of primary-backup replication based on a simplifed version of Viewstamp replication.
//
//
//

import (
	"sync"
//	"fmt"
	"labrpc"
)

// the 3 possible server status
const (
	NORMAL = iota
	VIEWCHANGE
	RECOVERING
)

// PBServer defines the state of a replica server (either primary or backup)
type PBServer struct {
	mu             sync.Mutex          // Lock to protect shared access to this peer's state
	peers          []*labrpc.ClientEnd // RPC end points of all peers
	me             int                 // this peer's index into peers[]
	currentView    int                 // what this peer believes to be the current active view
	status         int                 // the server's current status (NORMAL, VIEWCHANGE or RECOVERING)
	lastNormalView int                 // the latest view which had a NORMAL status

	log         []interface{} // the log of "commands"
	commitIndex int           // all log entries <= commitIndex are considered to have been committed.
	req         []chan *PrepareArgs
	wai	    []chan *PrepareArgs
	//	waitsev     []int
	//	waitreq     []*PrepareArgs
	// ... other state that you might need ...
}

// Prepare defines the arguments for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC args struct
type PrepareArgs struct {
	View          int         // the primary's current view
	PrimaryCommit int         // the primary's commitIndex
	Index         int         // the index position at which the log entry is to be replicated on backups
	Entry         interface{} // the log entry to be replicated
}

// PrepareReply defines the reply for the Prepare RPC
// Note that all field names must start with a capital letter for an RPC reply struct
type PrepareReply struct {
	View    int  // the backup's current view
	Success bool // whether the Prepare request has been accepted or rejected
}

// RecoverArgs defined the arguments for the Recovery RPC
type RecoveryArgs struct {
	View   int // the view that the backup would like to synchronize with
	Server int // the server sending the Recovery RPC (for debugging)
}

type RecoveryReply struct {
	View          int           // the view of the primary
	Entries       []interface{} // the primary's log including entries replicated up to and including the view.
	PrimaryCommit int           // the primary's commitIndex
	Success       bool          // whether the Recovery request has been accepted or rejected
}

type ViewChangeArgs struct {
	View int // the new view to be changed into
}

type ViewChangeReply struct {
	LastNormalView int           // the latest view which had a NORMAL status at the server
	Log            []interface{} // the log at the server
	Success        bool          // whether the ViewChange request has been accepted/rejected
}

type StartViewArgs struct {
	View int           // the new view which has completed view-change
	Log  []interface{} // the log associated with the new new
}

type StartViewReply struct {
}

// GetPrimary is an auxilary function that returns the server index of the
// primary server given the view number (and the total number of replica servers)
func GetPrimary(view int, nservers int) int {
	return view % nservers
}

// IsCommitted is called by tester to check whether an index position
// has been considered committed by this server
func (srv *PBServer) IsCommitted(index int) (committed bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	//fmt.Printf("COMMIT INDEX %v INDEX %v\n",srv.commitIndex,index)
	if srv.commitIndex >= index {
		return true
	}
	return false
}

// ViewStatus is called by tester to find out the current view of this server
// and whether this view has a status of NORMAL.
func (srv *PBServer) ViewStatus() (currentView int, statusIsNormal bool) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	return srv.currentView, srv.status == NORMAL
}

// GetEntryAtIndex is called by tester to return the command replicated at
// a specific log index. If the server's log is shorter than "index", then
// ok = false, otherwise, ok = true
func (srv *PBServer) GetEntryAtIndex(index int) (ok bool, command interface{}) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	//fmt.Printf("SERVER %v has log %v, SHOULD greater than %v\n",srv.me,len(srv.log),index)
	if len(srv.log) > index {
		return true, srv.log[index]
	}
	return false, command
}

// Kill is called by tester to clean up (e.g. stop the current server)
// before moving on to the next test
func (srv *PBServer) Kill() {
	// Your code here, if necessary
}

// Make is called by tester to create and initalize a PBServer
// peers is the list of RPC endpoints to every server (including self)
// me is this server's index into peers.
// startingView is the initial view (set to be zero) that all servers start in
func Make(peers []*labrpc.ClientEnd, me int, startingView int) *PBServer {
	srv := &PBServer{
		peers:          peers,
		me:             me,
		currentView:    startingView,
		lastNormalView: startingView,
		status:         NORMAL,
	}
	// all servers' log are initialized with a dummy command at index 0
	var v interface{}
	srv.log = append(srv.log, v)
	// Your other initialization code here, if there's any
	for i:=0;i<len(srv.peers);i++{
		srv.req = append(srv.req,make(chan *PrepareArgs,20))
		srv.wai = append(srv.req,make(chan *PrepareArgs,20))
	}
	//fmt.Printf("initialized array of %v channels\n",len(srv.req))
	return srv
}

// Start() is invoked by tester on some replica server to replicate a
// command.  Only the primary should process this request by appending
// the command to its log and then return *immediately* (while the log is being replicated to backup servers).
// if this server isn't the primary, returns false.
// Note that since the function returns immediately, there is no guarantee that this command
// will ever be committed upon return, since the primary
// may subsequently fail before replicating the command to all servers
//
// The first return value is the index that the command will appear at
// *if it's eventually committed*. The second return value is the current
// view. The third return value is true if this server believes it is
// the primary.
func (srv *PBServer) Start(command interface{}) (
	index int, view int, ok bool) {
		srv.mu.Lock()
		defer srv.mu.Unlock()
		// do not process command if status is not NORMAL
		// and if i am not the primary in the current view
		if srv.status != NORMAL {
			return -1, srv.currentView, false
		} else if GetPrimary(srv.currentView, len(srv.peers)) != srv.me {
			return -1, srv.currentView, false
		}
		// Your code here
		srv.log = append(srv.log,command)
	//	fmt.Printf("%v servers in total, %v is primary with log %v and commit %v\n",len(srv.peers),srv.me,len(srv.log),srv.commitIndex)
		comm:=len(srv.log)-1
		preArgs := &PrepareArgs{
			Entry: command,
			View: srv.currentView,
			PrimaryCommit: srv.commitIndex,
			Index: len(srv.log)-1,
		}
	//	fmt.Printf("len of peers %v, command %v, view %v\n",len(srv.peers),command,srv.currentView)
		for i:=0;i<len(srv.peers);i++{
			if i==srv.me{continue}
			srv.req[i] <- preArgs
		}
		preReplyChan := srv.send(srv.req)
		srv.collect(preReplyChan,preArgs)
	return comm, srv.currentView, true

}
func (srv *PBServer) send(chanarr []chan*PrepareArgs)  chan*PrepareReply{
	preReplyChan :=make(chan *PrepareReply,len(srv.peers))
	for i := 0; i < len(srv.peers); i++ {
		if i==srv.me{ continue}
		go func(server int) {
			var ok bool
			var reply PrepareReply
			ok = srv.sendPrepare(server,chanarr[server],&reply)
			if ok {
				preReplyChan <- &reply
			} else{
				preReplyChan <- nil
			}
		}(i)
	//	fmt.Printf("server %v finished\n",i)
	}
	return preReplyChan

}

func (srv *PBServer) collect(preReplyChan chan  *PrepareReply,preArgs *PrepareArgs) {
	successReplies :=0
	majority := len(srv.peers)/2
	precommit:=srv.commitIndex
	go func() {
		var nReplies int
		for r := range preReplyChan {
			nReplies++
			if r != nil && r.Success {
				successReplies++
				//successReplies = append(successReplies, r)
			} else if r!=nil && r.View>srv.currentView{
	//			fmt.Printf("THIS server %v is not primary for %v, there is larger view\n",srv.me,preArgs.Entry)
				return
			}
			if successReplies == majority{
	//			fmt.Printf("Most replicas prepared\n")
				srv.commitIndex++
	//			fmt.Printf("Current Commit Index %v\n",srv.commitIndex)
				break
			}
			if nReplies == len(srv.peers)-1 {
				break
			}
		}
	//	fmt.Printf("The commit index of server %v now is %v, before prepare is %v\n",srv.me,srv.commitIndex,precommit)
		if srv.commitIndex==precommit{
	//		fmt.Printf("Need to recall prepare\n")
			for i:=0;i<len(srv.peers);i++{
				if i==srv.me{continue}
				preArgs.View= srv.currentView
				//preArgs.PrimaryCommit= srv.commitIndex,
				srv.wai[i] <- preArgs
	//			fmt.Printf("add command %v to server %v wait channel\n",preArgs.Entry,i)
			}
			waitReplyChan:=srv.send(srv.wai)
			srv.collect(waitReplyChan,preArgs)
		}
	}()

}
// exmple code to send an AppendEntries RPC to a server.
// server is the index of the target server in srv.peers[].
// expects RPC arguments in args.
// The RPC library fills in *reply with RPC reply, so caller should pass &reply.
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
func (srv *PBServer) sendPrepare(server int, argschan chan *PrepareArgs, reply *PrepareReply) bool {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	args:= <-argschan
	//fmt.Printf("%v is going to add command %v on %v\n",server,args.Entry,args.Index)
	ok := srv.peers[server].Call("PBServer.Prepare", args, reply)
	return ok
}

// Prepare is the RPC handler for the Prepare RPC
func (srv *PBServer) Prepare(args *PrepareArgs, reply *PrepareReply) {
	//srv.mu.Lock()
	//defer srv.mu.Unlock()
	//fmt.Printf("Server %v with log %v view %v tries to view %v replicate %v at %v\n",srv.me,len(srv.log),srv.currentView,args.View,args.Entry, args.Index)
	if args.View != srv.currentView || args.Index != len(srv.log){
		reply.View = srv.currentView
		if args.View < srv.currentView{
	//		fmt.Printf("server %v Return false since view %v is smaller for %v\n",srv.me,args.View,args.Entry)
			reply.Success = false
			return
		}else if len(srv.log)>args.Index && srv.currentView==args.View{

			srv.log[args.Index] = args.Entry
			reply.Success = true
			return
		}else{
			reply.Success = false
		}
	//	fmt.Printf("Server %v with View %v should be %v. num log is %v, should be %v\n",srv.me,srv.currentView,args.View,len(srv.log),args.Index)
		if  args.View > srv.currentView || args.Index > len(srv.log){
			if srv.status != NORMAL{
	//			fmt.Printf("server %v in status %v, cannot recover to log %v\n",srv.me,srv.status,args.Index)
				return
			}
	//		fmt.Printf("Server %v starts recovery.\n",srv.me)
			srv.lastNormalView = srv.currentView
			srv.status = RECOVERING
			//majority := len(srv.peers)/2+1
			rcArgs := &RecoveryArgs{
				View: args.View,
				Server: srv.me,
			}
	//		fmt.Printf("-----recovery------%v send request to primary %v\n",srv.me,GetPrimary(args.View,len(srv.peers)))
			var rcreply RecoveryReply
			ok:=srv.peers[ GetPrimary(args.View,len(srv.peers))].Call("PBServer.Recovery", rcArgs, &rcreply)
			if ok && rcreply.Success && rcreply.View==args.View{
				srv.currentView = args.View
				srv.commitIndex = rcreply.PrimaryCommit
				srv.log = rcreply.Entries
				srv.status = NORMAL
	//			fmt.Printf("server %v finished recovery with view %v and log %v\n",srv.me,srv.currentView,srv.log)
				return

			}else{
	//			fmt.Printf("server %v Failed recovery with view %v and log %v for log %v\n",srv.me,srv.currentView,len(srv.log),args.Index)
			}
		}
	} else{
		srv.log = append(srv.log,args.Entry)
	//	fmt.Printf("Replicate succeed. Server %v now has %v entries.\n",srv.me,len(srv.log))
		reply.Success = true
		reply.View = srv.currentView
	}
}
// Recovery is the RPC handler for the Recovery RPC
func (srv *PBServer) Recovery(args *RecoveryArgs, reply *RecoveryReply) {
	// Your code here
	//		srv.mu.Lock()
	//		defer srv.mu.Unlock()
//	fmt.Printf("---recovery of server %v---------------server %v status %v\n",args.Server,srv.me,srv.status)
	if srv.status == NORMAL {
		reply.Entries = srv.log
		reply.PrimaryCommit = srv.commitIndex
		reply.View  = args.View
		reply.Success = true

	} else{
		reply.Success = false
	}
	return
}
// Some external oracle prompts the primary of the newView to
// switch to the newView.
// PromptViewChange just kicks start the view change protocol to move to the newView
// It does not block waiting for the view change process to complete.
func (srv *PBServer) PromptViewChange(newView int) {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	newPrimary := GetPrimary(newView, len(srv.peers))

	if newPrimary != srv.me { //only primary of newView should do view change
		return
	} else if newView <= srv.currentView {
		return
	}
	vcArgs := &ViewChangeArgs{
		View: newView,
	}
	vcReplyChan := make(chan *ViewChangeReply, len(srv.peers))
	// send ViewChange to all servers including myself
	for i := 0; i < len(srv.peers); i++ {
		go func(server int) {
			var reply ViewChangeReply
			ok := srv.peers[server].Call("PBServer.ViewChange", vcArgs, &reply)
			//fmt.Printf("node-%d received reply ok=%v reply=%v\n", srv.me,ok, r.reply)
			if ok {
				vcReplyChan <- &reply
			}
		}(i)
	}

	// wait to receive ViewChange replies
	// if view change succeeds, send StartView RPC
	go func() {
		var successReplies []*ViewChangeReply
		var nReplies int
		majority := len(srv.peers)/2 + 1
		for r := range vcReplyChan {
			nReplies++
			if r != nil && r.Success {
				successReplies = append(successReplies, r)
	//			fmt.Printf("new primary %v recieved %v viewchange replies, %v succeed\n",srv.me,nReplies,len(successReplies))
			}//else if r!=nil{fmt.Printf("new primary %v recieved %v viewchange replies, %v failed\n",srv.me,nReplies,nReplies-len(successReplies))}
			if nReplies == len(srv.peers) || len(successReplies) == majority {
				break
			}
		}
		ok, log := srv.determineNewViewLog(successReplies)
		if !ok {
			return
		}
		srv.log = log
	//	fmt.Printf("primary %v previous commitIndex %v\n",srv.me,srv.commitIndex)
		srv.commitIndex = len(srv.log)-1
	//	fmt.Printf("primary %v current commitIndex %v\n",srv.me,srv.commitIndex)
		svArgs := &StartViewArgs{
			View: vcArgs.View,
			Log:  log,
		}
		// send StartView to all servers including myself
		for i := 0; i < len(srv.peers); i++ {
			var reply StartViewReply
			go func(server int) {
	//			fmt.Printf("node-%d sending StartView v=%d to node-%d\n", srv.me, svArgs.View, server)
				srv.peers[server].Call("PBServer.StartView", svArgs, &reply)
			}(i)
		}
	}()

}

// determineNewViewLog is invoked to determine the log for the newView based on
// the collection of replies for successful ViewChange requests.
// if a=--ol,uhnn  uccessful replies exist, then ok is set to true.
// otherwise, ok = false.
func (srv *PBServer) determineNewViewLog(successReplies []*ViewChangeReply) (
	ok bool, newViewLog []interface{}) {
		// Your code here
		ok=false
		majority := len(srv.peers)/2 + 1
		if len(successReplies) == majority{
			ok=true
		}else{
			return ok,newViewLog
		}
		lastnormalarr:=[][]interface{}{successReplies[0].Log}
		if len(successReplies)==1 {return ok,successReplies[0].Log}
		max := successReplies[0].LastNormalView
//		fmt.Printf("0th reply has last normal view %v and log %v\n",max,successReplies[0].Log[len(successReplies[0].Log)-1])
		for i:=1;i<len(successReplies);i++{
			r:=successReplies[i]
//			fmt.Printf("%vth reply has last normal view %v and log %v\n",i,r.LastNormalView,r.Log[len(r.Log)-1])
			if r.LastNormalView > max{
				max = r.LastNormalView
				lastnormalarr =[][]interface{}{r.Log}
			}else if r.LastNormalView == max{
				lastnormalarr = append(lastnormalarr,r.Log)
			}
//			fmt.Printf("current max last normal view %v\n",max)
		}
		if len(lastnormalarr) == 1{
			newViewLog = lastnormalarr[0]
		}else{
			maxlog := lastnormalarr[0]
			for i:=0;i<len(lastnormalarr);i++{
				lo := lastnormalarr[i]
//				fmt.Printf("maxlog is %v, log is %v\n",maxlog,lo)
				if len(lo)>len(maxlog){maxlog = lo}
			}
			newViewLog = maxlog
		}
	//	fmt.Printf("newViewLog is %v\n",newViewLog)
		return ok, newViewLog
	}

	// ViewChange is the RPC handler to process ViewChange RPC.
	func (srv *PBServer) ViewChange(args *ViewChangeArgs, reply *ViewChangeReply) {
		// Your code here
		srv.mu.Lock()
        	defer srv.mu.Unlock()
		if args.View>srv.currentView {
	//		fmt.Printf("server %v starts view change, set view %v to %v\n",srv.me,srv.currentView,args.View)
			srv.lastNormalView = srv.currentView
			srv.currentView = args.View
			srv.status = VIEWCHANGE
			reply.Success = true
			reply.LastNormalView = srv.lastNormalView
			reply.Log = srv.log
		}else{
	//		fmt.Printf("server %v with view %v greater than %v failed to view change\n",srv.me,srv.currentView,args.View)
			reply.Success = false
		}

	}

	// StartView is the RPC handler to process StartView RPC.
	func (srv *PBServer) StartView(args *StartViewArgs, reply *StartViewReply) {
		// Your code here
		srv.mu.Lock()
        	defer srv.mu.Unlock()
		if srv.currentView <= args.View{
	//		fmt.Printf("server %v change view %v to %v in startview\n",srv.me,srv.currentView,args.View)
			srv.currentView = args.View
	//		fmt.Printf("server %v back to normal\n",srv.me)
			srv.status = NORMAL
		}
	}
