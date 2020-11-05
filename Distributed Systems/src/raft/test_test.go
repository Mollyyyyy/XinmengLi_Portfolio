package raft

//
// Raft tests.
//
// we will use the original test_test.go to test your code for grading.
// so, while you can modify this code to help you debug, please
// test with the original before submitting.
//

import "testing"
import "fmt"
import "time"
import "math/rand"
import "sync/atomic"
import "sync"

// The tester generously allows solutions to complete elections in one second
// (much more than the paper's range of timeouts).
const RaftElectionTimeout = 1000 * time.Millisecond

func TestInitialElection2A(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2A): initial election ...\n")

//	fmt.Printf(" is a leader elected?\n")
	cfg.checkOneLeader()

//	fmt.Printf(" does the leader+term stay the same if there is no network failure?\n")
	term1 := cfg.checkTerms()
	time.Sleep(2 * RaftElectionTimeout)
	term2 := cfg.checkTerms()
	if term1 != term2 {
		fmt.Printf("warning: term changed even though there were no failures")
	}

	fmt.Printf("  ... Passed\n")
}

func TestReElection2A(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2A): election after network failure ...\n")

	leader1 := cfg.checkOneLeader()

//	fmt.Printf(" ~~~~~~~~~~~if the leader %v disconnects, a new one should be elected.\n",leader1)
	cfg.disconnect(leader1)
	cfg.checkOneLeader()

//	fmt.Printf(" ~~~~~~~~~~if the old leader %v rejoins, that shouldn't  disturb the old leader.\n",leader1)
	cfg.connect(leader1)
	leader2 := cfg.checkOneLeader()

//	fmt.Printf(" ~~~~~~~~~~~if there's no quorum, no leader should be elected.\n")
//	fmt.Printf(" disconnect %v, %v\n",leader2,(leader2+1)%servers)
	cfg.disconnect(leader2)
	cfg.disconnect((leader2 + 1) % servers)
//	fmt.Printf(" ~~~~~~~~~~wait 2 election round\n")
	time.Sleep(2 * RaftElectionTimeout)
	cfg.checkNoLeader()

//	fmt.Printf(" ~~~~~~~~~~~if a quorum arises, it should elect a leader.\n")
//	fmt.Printf(" reconnect %v\n",(leader2+1)%servers)
	cfg.connect((leader2 + 1) % servers)
	cfg.checkOneLeader()

//	fmt.Printf(" ~~~~~~~~~~re-join of last node shouldn't prevent leader from existing.\n")
//	fmt.Printf(" reconnect %v\n",leader2)
	cfg.connect(leader2)
	cfg.checkOneLeader()

	fmt.Printf("  ... Passed\n")
}

func TestBasicAgree2B(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): basic agreement ...\n")

	iters := 3
	for index := 1; index < iters+1; index++ {
		nd, _ := cfg.nCommitted(index)
		if nd > 0 {
			t.Fatalf("some have committed before Start()")
		}

		xindex := cfg.one(index*100, servers)
		if xindex != index {
			t.Fatalf("got index %v but expected %v", xindex, index)
		}
	}

	fmt.Printf("  ... Passed\n")
}

func TestFailAgree2B(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): agreement despite follower disconnection ...\n")

	cfg.one(101, servers)

	leader := cfg.checkOneLeader()
//	fmt.Printf("--- follower %v network disconnection\n",(leader + 1) % servers)
	cfg.disconnect((leader + 1) % servers)

//	fmt.Printf("--- agree despite one disconnected server?\n")
	cfg.one(102, servers-1)
	cfg.one(103, servers-1)
	time.Sleep(RaftElectionTimeout)
	cfg.one(104, servers-1)
	cfg.one(105, servers-1)

//	fmt.Printf("--- re-connect %v\n",(leader + 1) % servers)
	cfg.connect((leader + 1) % servers)

//	fmt.Printf("--- agree with full set of servers?\n")
	cfg.one(106, servers)
	time.Sleep(RaftElectionTimeout)
	cfg.one(107, servers)

	fmt.Printf("  ... Passed\n")
}

func TestFailNoAgree2B(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): no agreement if too many followers disconnect ...\n")

	cfg.one(10, servers)

	leader := cfg.checkOneLeader()
//	fmt.Printf("----- 3 of 5 followers %v %v %v disconnect\n",(leader + 1) % servers,(leader + 2) % servers,(leader + 3) % servers)
	cfg.disconnect((leader + 1) % servers)
	cfg.disconnect((leader + 2) % servers)
	cfg.disconnect((leader + 3) % servers)

	index, _, ok := cfg.rafts[leader].Start(20)
	if ok != true {
		t.Fatalf("leader rejected Start()")
	}
	if index != 2 {
		t.Fatalf("expected index 2, got %v", index)
	}
//	fmt.Printf("----- sleep 2 election timeout\n")
	time.Sleep(2 * RaftElectionTimeout)

	n, _ := cfg.nCommitted(index)
	if n > 0 {
		t.Fatalf("%v committed but no majority", n)
	}

//	fmt.Printf("----- reconnect to three servers\n")
	cfg.connect((leader + 1) % servers)
	cfg.connect((leader + 2) % servers)
	cfg.connect((leader + 3) % servers)

//	fmt.Printf("----- the disconnected majority may have chosen a leader from among their own ranks, forgetting index 2. or perhaps\n")
	leader2 := cfg.checkOneLeader()
	index2, _, ok2 := cfg.rafts[leader2].Start(30)
	if ok2 == false {
		t.Fatalf("leader2 rejected Start()")
	}
	if index2 < 2 || index2 > 3 {
		t.Fatalf("unexpected index %v", index2)
	}

	cfg.one(1000, servers)

	fmt.Printf("  ... Passed\n")
}

func TestConcurrentStarts2B(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): concurrent Start()s ...\n")

	var success bool
loop:
	for try := 0; try < 5; try++ {
		if try > 0 {
		//	fmt.Printf("--- give solution some time 3s to settle\n")
			time.Sleep(3 * time.Second)
		}

		leader := cfg.checkOneLeader()
		_, term, ok := cfg.rafts[leader].Start(1)
		if !ok {
		//	fmt.Printf("--- leader %v moved on really quickly\n",leader)
			continue
		}

		iters := 5
		var wg sync.WaitGroup
		is := make(chan int, iters)
		for ii := 0; ii < iters; ii++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				i, term1, ok := cfg.rafts[leader].Start(100 + i)
				if term1 != term {
					return
				}
				if ok != true {
					return
				}
				is <- i
			}(ii)
		}

		wg.Wait()
		close(is)

		for j := 0; j < servers; j++ {
			if t, _ := cfg.rafts[j].GetState(); t != term {
			//	fmt.Printf("--- term changed -- can't expect low RPC counts\n")
				continue loop
			}
		}

		failed := false
		cmds := []int{}
		for index := range is {
			cmd := cfg.wait(index, servers, term)
			if ix, ok := cmd.(int); ok {
				if ix == -1 {
				//	fmt.Printf("--- peers have moved on to later terms so we can't expect all Start()s to have succeeded\n")
					failed = true
					break
				}
				cmds = append(cmds, ix)
			} else {
				t.Fatalf("value %v is not an int", cmd)
			}
		}

		if failed {
		//	fmt.Printf("--- avoid leaking goroutines\n")
			go func() {
				for range is {
				}
			}()
			continue
		}

		for ii := 0; ii < iters; ii++ {
			x := 100 + ii
			ok := false
			for j := 0; j < len(cmds); j++ {
				if cmds[j] == x {
					ok = true
				}
			}
			if ok == false {
				t.Fatalf("cmd %v missing in %v", x, cmds)
			}
		}

		success = true
		break
	}

	if !success {
		t.Fatalf("term changed too often")
	}

	fmt.Printf("  ... Passed\n")
}

func TestRejoin2B(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): rejoin of partitioned leader ...\n")

	cfg.one(101, servers)

	// leader network failure
	leader1 := cfg.checkOneLeader()
	//fmt.Printf("--- disconnect leader %v\n",leader1)
	cfg.disconnect(leader1)

	//fmt.Printf("--- make old leader try to agree on some entries, send 102, 103, 104\n")
	cfg.rafts[leader1].Start(102)
	cfg.rafts[leader1].Start(103)
	cfg.rafts[leader1].Start(104)

	//fmt.Printf("--- new leader commits, also for index=2\n")
	cfg.one(103, 2)

	leader2 := cfg.checkOneLeader()
	//fmt.Printf("--- new leader network failure, disconnect %v\n",leader2)
	cfg.disconnect(leader2)

	//fmt.Printf("--- old leader %v connected again\n",leader1)
	cfg.connect(leader1)

	cfg.one(104, 2)

	// all together now
	//fmt.Printf("--- leader %v reconnect\n",leader2)
	cfg.connect(leader2)

	cfg.one(105, servers)

	fmt.Printf("  ... Passed\n")
}

func TestBackup2B(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): leader backs up quickly over incorrect follower logs ...\n")

	cfg.one(rand.Int(), servers)

	leader1 := cfg.checkOneLeader()
	//fmt.Printf("--- put leader and one follower in a partition, disconnect %v %v %v\n",(leader1 + 2) % servers,(leader1 + 3) % servers,(leader1 + 4) % servers)
	cfg.disconnect((leader1 + 2) % servers)
	cfg.disconnect((leader1 + 3) % servers)
	cfg.disconnect((leader1 + 4) % servers)

	//fmt.Printf("--- submit lots of commands(50 intotal) that won't commit\n")
	for i := 0; i < 50; i++ {
		cfg.rafts[leader1].Start(rand.Int())
	}
	//fmt.Printf("--- wait half election timeout\n")
	time.Sleep(RaftElectionTimeout / 2)

//	fmt.Printf("--- disconnect %v, %v\n", (leader1 + 0) % servers,(leader1 + 1) % servers)
	cfg.disconnect((leader1 + 0) % servers)
	cfg.disconnect((leader1 + 1) % servers)

	//fmt.Printf("--- allow other partition to recover, reconnect %v %v %v\n",(leader1 + 2) % servers,(leader1 + 3) % servers,(leader1 + 4) % servers)
	cfg.connect((leader1 + 2) % servers)
	cfg.connect((leader1 + 3) % servers)
	cfg.connect((leader1 + 4) % servers)

	//fmt.Printf("--- lots of successful commands to new group.\n")
	for i := 0; i < 50; i++ {
		cfg.one(rand.Int(), 3)
	}

	leader2 := cfg.checkOneLeader()
	other := (leader1 + 2) % servers
	if leader2 == other {
		other = (leader2 + 1) % servers
	}
	//fmt.Printf("--- now another partitioned leader and one follower, disconnect %v\n",other)
	cfg.disconnect(other)

	//fmt.Printf("--- lots more commands that won't commit\n")
	for i := 0; i < 50; i++ {
		cfg.rafts[leader2].Start(rand.Int())
	}
	//fmt.Printf("--- sleep half election timeout\n")
	time.Sleep(RaftElectionTimeout / 2)

	//fmt.Printf("--- bring original leader back to life, disconnect all and then connect %v %v %v\n",(leader1 + 0) % servers,(leader1 + 1) % servers,other)
	for i := 0; i < servers; i++ {
		cfg.disconnect(i)
	}
	cfg.connect((leader1 + 0) % servers)
	cfg.connect((leader1 + 1) % servers)
	cfg.connect(other)

	//fmt.Printf("--- lots of successful commands to new group.\n")
	for i := 0; i < 50; i++ {
		cfg.one(rand.Int(), 3)
	}

	// now everyone
	//fmt.Printf("--- connect all servers\n")
	for i := 0; i < servers; i++ {
		cfg.connect(i)
	}
	cfg.one(rand.Int(), servers)

	fmt.Printf("  ... Passed\n")
}

func TestCount2B(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2B): RPC counts aren't too high ...\n")

	rpcs := func() (n int) {
		for j := 0; j < servers; j++ {
			n += cfg.rpcCount(j)
		}
		return
	}

	leader := cfg.checkOneLeader()

	total1 := rpcs()

	if total1 > 30 || total1 < 1 {
		t.Fatalf("too many or few RPCs (%v) to elect initial leader\n", total1)
	}

	var total2 int
	var success bool
loop:
	for try := 0; try < 5; try++ {
		if try > 0 {
			// give solution some time to settle
			//fmt.Printf("--- sleep 3s\n")
			time.Sleep(3 * time.Second)
		}

		leader = cfg.checkOneLeader()
		total1 = rpcs()

		iters := 10
		starti, term, ok := cfg.rafts[leader].Start(1)
		if !ok {
			// leader moved on really quickly
			continue
		}
		cmds := []int{}
		for i := 1; i < iters+2; i++ {
			x := int(rand.Int31())
			cmds = append(cmds, x)
			index1, term1, ok := cfg.rafts[leader].Start(x)
			if term1 != term {
				// Term changed while starting
				continue loop
			}
			if !ok {
				// No longer the leader, so term has changed
				continue loop
			}
			if starti+i != index1 {
				t.Fatalf("Start() failed")
			}
		}

		for i := 1; i < iters+1; i++ {
			cmd := cfg.wait(starti+i, servers, term)
			if ix, ok := cmd.(int); ok == false || ix != cmds[i-1] {
				if ix == -1 {
					// term changed -- try again
					continue loop
				}
				t.Fatalf("wrong value %v committed for index %v; expected %v\n", cmd, starti+i, cmds)
			}
		}

		failed := false
		total2 = 0
		for j := 0; j < servers; j++ {
			if t, _ := cfg.rafts[j].GetState(); t != term {
				// term changed -- can't expect low RPC counts
				// need to keep going to update total2
				failed = true
			}
			total2 += cfg.rpcCount(j)
		}

		if failed {
			continue loop
		}

		if total2-total1 > (iters+1+3)*3 {
			t.Fatalf("too many RPCs (%v) for %v entries\n", total2-total1, iters)
		}

		success = true
		break
	}

	if !success {
		t.Fatalf("term changed too often")
	}

	time.Sleep(RaftElectionTimeout)

	total3 := 0
	for j := 0; j < servers; j++ {
		total3 += cfg.rpcCount(j)
	}

	if total3-total2 > 3*20 {
		t.Fatalf("too many RPCs (%v) for 1 second of idleness\n", total3-total2)
	}

	fmt.Printf("  ... Passed\n")
}

func TestPersist12C(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2C): basic persistence ...\n")

	cfg.one(11, servers)

	//fmt.Printf("--- crash\n")
	for i := 0; i < servers; i++ {
		cfg.start1(i)
	}
	//fmt.Printf("---  re-start all\n")
	for i := 0; i < servers; i++ {
		//fmt.Printf("--- disconnect %v \n",i)
		cfg.disconnect(i)
		//fmt.Printf("--- connect %v \n",i)
		cfg.connect(i)
	}
	//fmt.Printf("--- send command 12\n")
	cfg.one(12, servers)

	leader1 := cfg.checkOneLeader()
	//fmt.Printf("--- disconnect %v \n",leader1)
	cfg.disconnect(leader1)
	cfg.start1(leader1)
	//fmt.Printf("--- connect %v \n",leader1)
	cfg.connect(leader1)
	//fmt.Printf("--- send command 13\n")
	cfg.one(13, servers)

	leader2 := cfg.checkOneLeader()
	//fmt.Printf("--- disconnect %v \n",leader2)
	cfg.disconnect(leader2)
	//fmt.Printf("--- send command 14\n")
	cfg.one(14, servers-1)
	cfg.start1(leader2)
	//fmt.Printf("--- connect %v \n",leader2)
	cfg.connect(leader2)

	cfg.wait(4, servers, -1) // wait for leader2 to join before killing i3

	i3 := (cfg.checkOneLeader() + 1) % servers
	//fmt.Printf("--- disconnect %v \n",i3)
	cfg.disconnect(i3)
	//fmt.Printf("--- send command 15\n")
	cfg.one(15, servers-1)
	cfg.start1(i3)
	//fmt.Printf("--- connect %v \n",i3)
	cfg.connect(i3)
	//fmt.Printf("--- send command 16\n")
	cfg.one(16, servers)

	fmt.Printf("  ... Passed\n")
}

func TestPersist22C(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2C): more persistence ...\n")

	index := 1
	for iters := 0; iters < 5; iters++ {
		//fmt.Printf("--- send command %v to %v\n",10+index,servers)
		cfg.one(10+index, servers)
		index++

		leader1 := cfg.checkOneLeader()
		//fmt.Printf("--- disconnect %v \n",(leader1 + 1) % servers)
		cfg.disconnect((leader1 + 1) % servers)
		//fmt.Printf("--- disconnect %v \n",(leader1 + 2) % servers)
		cfg.disconnect((leader1 + 2) % servers)
		//fmt.Printf("--- send command %v to %v\n",10+index,servers-2)
		cfg.one(10+index, servers-2)
		index++
		//fmt.Printf("--- disconnect %v %v %v \n",(leader1 + 0) % servers,(leader1 + 3) % servers,(leader1 + 4) % servers)
		cfg.disconnect((leader1 + 0) % servers)
		cfg.disconnect((leader1 + 3) % servers)
		cfg.disconnect((leader1 + 4) % servers)

		cfg.start1((leader1 + 1) % servers)
		cfg.start1((leader1 + 2) % servers)
		//fmt.Printf("--- connect %v %v \n",(leader1 + 1) % servers,(leader1 + 2) % servers)
		cfg.connect((leader1 + 1) % servers)
		cfg.connect((leader1 + 2) % servers)

		time.Sleep(RaftElectionTimeout)

		cfg.start1((leader1 + 3) % servers)
		//fmt.Printf("--- connect %v \n",(leader1 + 3) % servers)
		cfg.connect((leader1 + 3) % servers)
		//fmt.Printf("--- send command  %v to %v \n",10+index,servers-2)
		cfg.one(10+index, servers-2)
		index++
		//fmt.Printf("--- connect %v %v\n",(leader1 + 4) % servers,(leader1 + 0) % servers)
		cfg.connect((leader1 + 4) % servers)
		cfg.connect((leader1 + 0) % servers)
	}
	//fmt.Printf("--- send command 1000 to all\n")
	cfg.one(1000, servers)

	fmt.Printf("  ... Passed\n")
}

func TestPersist32C(t *testing.T) {
	servers := 3
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2C): partitioned leader and one follower crash, leader restarts ...\n")
	//fmt.Printf("--- send command 101 to 3\n")
	cfg.one(101, 3)

	leader := cfg.checkOneLeader()
	//fmt.Printf("--- disconnect %v \n",(leader + 2) % servers)
	cfg.disconnect((leader + 2) % servers)
	//fmt.Printf("--- send command 102 to 2\n")
	cfg.one(102, 2)
	//fmt.Printf("--- crash %v %v\n",(leader + 0) % servers,(leader + 1) % servers)
	cfg.crash1((leader + 0) % servers)
	cfg.crash1((leader + 1) % servers)
	//fmt.Printf("--- connect %v %v \n",(leader + 2) % servers,(leader + 0) % servers)
	cfg.connect((leader + 2) % servers)
	cfg.start1((leader + 0) % servers)
	cfg.connect((leader + 0) % servers)
	//fmt.Printf("--- send command 103 to 2\n")
	cfg.one(103, 2)

	cfg.start1((leader + 1) % servers)
	//fmt.Printf("--- connect %v \n",(leader + 1) % servers)
	cfg.connect((leader + 1) % servers)
	//fmt.Printf("--- send command 104 to all\n")
	cfg.one(104, servers)

	fmt.Printf("  ... Passed\n")
}

//
// Test the scenarios described in Figure 8 of the extended Raft paper. Each
// iteration asks a leader, if there is one, to insert a command in the Raft
// log.  If there is a leader, that leader will fail quickly with a high
// probability (perhaps without committing the command), or crash after a while
// with low probability (most likey committing the command).  If the number of
// alive servers isn't enough to form a majority, perhaps start a new server.
// The leader in a new term may try to finish replicating log entries that
// haven't been committed yet.
//
func TestFigure82C(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	fmt.Printf("Test (2C): Figure 8 ...\n")
	//fmt.Printf("--- send command random to 1\n")
	cfg.one(rand.Int(), 1)

	nup := servers
	for iters := 0; iters < 1000; iters++ {
		leader := -1
		for i := 0; i < servers; i++ {
			if cfg.rafts[i] != nil {
				//fmt.Printf("--- send command random to %v\n",i)
				_, _, ok := cfg.rafts[i].Start(rand.Int())
				if ok {
					leader = i
				}
			}
		}

		if (rand.Int() % 1000) < 100 {
			ms := rand.Int63() % (int64(RaftElectionTimeout/time.Millisecond) / 2)
		//	fmt.Printf("--- sleep %v ms\n",ms)
			time.Sleep(time.Duration(ms) * time.Millisecond)
		} else {
			ms := (rand.Int63() % 13)
		//	fmt.Printf("--- sleep %v ms\n",ms)
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}

		if leader != -1 {
		//	fmt.Printf("--- crash %v \n",leader)
			cfg.crash1(leader)
			nup -= 1
		}

		if nup < 3 {
			s := rand.Int() % servers
			if cfg.rafts[s] == nil {
				cfg.start1(s)
			//	fmt.Printf("--- connect %v \n",s)
				cfg.connect(s)
				nup += 1
			}
		}
	}

	for i := 0; i < servers; i++ {
		if cfg.rafts[i] == nil {
			cfg.start1(i)
			//fmt.Printf("--- connect %v \n",i)
			cfg.connect(i)
		}
	}
	//fmt.Printf("--- send command random to all\n")
	cfg.one(rand.Int(), servers)

	fmt.Printf("  ... Passed\n")
}

func TestUnreliableAgree2C(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, true)
	defer cfg.cleanup()

	fmt.Printf("Test (2C): unreliable agreement ...\n")

	var wg sync.WaitGroup

	for iters := 1; iters < 50; iters++ {
		for j := 0; j < 4; j++ {
			wg.Add(1)
			go func(iters, j int) {
				defer wg.Done()
				//fmt.Printf("--- send command %v to %v\n",(100*iters)+j,1)
				cfg.one((100*iters)+j, 1)
			}(iters, j)
		}
		//fmt.Printf("--- send command %v to %v\n",iters,1)
		cfg.one(iters, 1)
	}

	cfg.setunreliable(false)

	wg.Wait()
	//fmt.Printf("--- send command 100 to all\n")
	cfg.one(100, servers)

	fmt.Printf("  ... Passed\n")
}

func TestFigure8Unreliable2C(t *testing.T) {
	servers := 5
	cfg := make_config(t, servers, true)
	defer cfg.cleanup()

	fmt.Printf("Test (2C): Figure 8 (unreliable) ...\n")
	//fmt.Printf("--- send random command to 1\n")
	cfg.one(rand.Int()%10000, 1)

	nup := servers
	for iters := 0; iters < 1000; iters++ {
		if iters == 200 {
			cfg.setlongreordering(true)
		}
		leader := -1
		//fmt.Printf("--- send random command to all\n")
		for i := 0; i < servers; i++ {
			_, _, ok := cfg.rafts[i].Start(rand.Int() % 10000)
			if ok && cfg.connected[i] {
				leader = i
			}
		}

		if (rand.Int() % 1000) < 100 {
			ms := rand.Int63() % (int64(RaftElectionTimeout/time.Millisecond) / 2)
		//	fmt.Printf("--- sleep %v ms\n",ms)
			time.Sleep(time.Duration(ms) * time.Millisecond)
		} else {
			ms := (rand.Int63() % 13)
	//		fmt.Printf("--- sleep %v ms\n",ms)
			time.Sleep(time.Duration(ms) * time.Millisecond)
		}

		if leader != -1 && (rand.Int()%1000) < int(RaftElectionTimeout/time.Millisecond)/2 {
			//fmt.Printf("--- XXXXXX disconnect %v\n",leader)
			cfg.disconnect(leader)
			nup -= 1
		}

		if nup < 3 {
			s := rand.Int() % servers
			if cfg.connected[s] == false {
				//fmt.Printf("--- XXXXXXX connect %v\n",s)
				cfg.connect(s)
				nup += 1
			}
		}
	}

	for i := 0; i < servers; i++ {
		if cfg.connected[i] == false {
			//fmt.Printf("--- XXXXXXXX connect %v\n",i)
			cfg.connect(i)
		}
	}
	
	//fmt.Printf("--- send random command to all\n")
	cfg.one(rand.Int()%10000, servers)

	fmt.Printf("  ... Passed\n")
}

func internalChurn(t *testing.T, unreliable bool) {

	if unreliable {
		fmt.Printf("Test (2C): unreliable churn ...\n")
	} else {
		fmt.Printf("Test (2C): churn ...\n")
	}

	servers := 5
	cfg := make_config(t, servers, unreliable)
	defer cfg.cleanup()

	stop := int32(0)

	// create concurrent clients
	cfn := func(me int, ch chan []int) {
		var ret []int
		ret = nil
		defer func() { ch <- ret }()
		values := []int{}
		for atomic.LoadInt32(&stop) == 0 {
			x := rand.Int()
			index := -1
			ok := false
			for i := 0; i < servers; i++ {
				// try them all, maybe one of them is a leader
				cfg.mu.Lock()
				rf := cfg.rafts[i]
				cfg.mu.Unlock()
				if rf != nil {
					index1, _, ok1 := rf.Start(x)
					if ok1 {
						ok = ok1
						index = index1
					}
				}
			}
			if ok {
				// maybe leader will commit our value, maybe not.
				// but don't wait forever.
				for _, to := range []int{10, 20, 50, 100, 200} {
					nd, cmd := cfg.nCommitted(index)
					if nd > 0 {
						if xx, ok := cmd.(int); ok {
							if xx == x {
								values = append(values, x)
							}
						} else {
							cfg.t.Fatalf("wrong command type")
						}
						break
					}
					time.Sleep(time.Duration(to) * time.Millisecond)
				}
			} else {
				time.Sleep(time.Duration(79+me*17) * time.Millisecond)
			}
		}
		ret = values
	}

	ncli := 3
	cha := []chan []int{}
	for i := 0; i < ncli; i++ {
		cha = append(cha, make(chan []int))
		go cfn(i, cha[i])
	}

	for iters := 0; iters < 20; iters++ {
		if (rand.Int() % 1000) < 200 {
			i := rand.Int() % servers
			//fmt.Printf("--- XXXXXX disconnect %v\n",i)
			cfg.disconnect(i)
		}

		if (rand.Int() % 1000) < 500 {
			i := rand.Int() % servers
			if cfg.rafts[i] == nil {
				cfg.start1(i)
			}
			//fmt.Printf("--- XXXXXX connect %v\n",i)
			cfg.connect(i)
		}

		if (rand.Int() % 1000) < 200 {
			i := rand.Int() % servers
			if cfg.rafts[i] != nil {
				//fmt.Printf("--- XXXXXX crash %v\n",i)
				cfg.crash1(i)
			}
		}

		// Make crash/restart infrequent enough that the peers can often
		// keep up, but not so infrequent that everything has settled
		// down from one change to the next. Pick a value smaller than
		// the election timeout, but not hugely smaller.
		time.Sleep((RaftElectionTimeout * 7) / 10)
	}

	time.Sleep(RaftElectionTimeout)
	cfg.setunreliable(false)
	for i := 0; i < servers; i++ {
		if cfg.rafts[i] == nil {
			cfg.start1(i)
		}
		//fmt.Printf("--- XXXXXX connect %v\n",i)
		cfg.connect(i)
	}

	atomic.StoreInt32(&stop, 1)

	values := []int{}
	for i := 0; i < ncli; i++ {
		vv := <-cha[i]
		if vv == nil {
			t.Fatal("client failed")
		}
		values = append(values, vv...)
	}

	time.Sleep(RaftElectionTimeout)

	lastIndex := cfg.one(rand.Int(), servers)

	really := make([]int, lastIndex+1)
	for index := 1; index <= lastIndex; index++ {
		v := cfg.wait(index, servers, -1)
		if vi, ok := v.(int); ok {
			really = append(really, vi)
		} else {
			t.Fatalf("not an int")
		}
	}

	for _, v1 := range values {
		ok := false
		for _, v2 := range really {
			if v1 == v2 {
				ok = true
			}
		}
		if ok == false {
			cfg.t.Fatalf("didn't find a value")
		}
	}

	fmt.Printf("  ... Passed\n")
}

func TestReliableChurn2C(t *testing.T) {
	internalChurn(t, false)
}

func TestUnreliableChurn2C(t *testing.T) {
	internalChurn(t, true)
}
