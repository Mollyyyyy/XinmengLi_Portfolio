package simplepb

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func majority(nservers int) int {
	return nservers/2 + 1
}

func Test1ABasicPB(t *testing.T) {
	servers := 3                        //3 servers
	primaryID := GetPrimary(0, servers) //primary ID is determined by view=0
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	for index := 1; index <= 10; index++ {
		xindex := cfg.replicateOne(primaryID, 1000+index, servers) // replicate command 1000+index, expected successful replication to all servers
//		fmt.Printf("XXXXXXXXXXXXXX INDEX %v finished xindex %v XXXXXXXXXXXXX\n",index,xindex)
		if xindex != index {
			t.Fatalf("got index %v but expected %v", xindex, index)
		}
	}
	fmt.Printf(" ... Passed\n")
}

func Test1AConcurrentPB(t *testing.T) {
	servers := 3                        //3 servers
	primaryID := GetPrimary(0, servers) //primary ID is determined by view=0
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	tries := 5
	for try := 0; try < tries; try++ {
		var wg sync.WaitGroup
		iters := 5
		for i := 0; i < iters; i++ {
			wg.Add(1)
			go func(x int) {
				defer wg.Done()
				val := 2000 + 100*try + x
				if _, _, ok := cfg.pbservers[primaryID].Start(val); !ok {
					t.Fatalf("node-%d rejected command %v\n", primaryID, val)
				}
			}(i)
		}
		wg.Wait()

		// wait for index (try + 1) * iters to be considered committed
		cfg.waitCommitted(primaryID, (try+1)*iters)

		// check that committed indexes [try*iters, (try+1)*iters] are identical at all servers
		var command interface{}
		for index := 1 + try*iters; index <= (try+1)*iters; index++ {
			cfg.checkCommittedIndex(index, command, majority(servers))
		}
	}
	fmt.Printf(" ... Passed\n")
}

func Test1AFailButCommitPB(t *testing.T) {
	servers := 3 //3 servers
	primaryID := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(primaryID, 3001, servers)

	var wg sync.WaitGroup
	for i := 0; i < 20; i += 4 {
		// disconnect a non-primary server
		cfg.disconnect((primaryID + 1) % servers)

		wg.Add(2)
		go func() {
			defer wg.Done()
			// agree despite replicate disconnected server?
			if _, _, ok := cfg.pbservers[primaryID].Start(3002 + i); !ok {
				t.Fatalf("node-%d rejected command %d\n", primaryID, 3002+i)
			}
			if _, _, ok := cfg.pbservers[primaryID].Start(3003 + i); !ok {
				t.Fatalf("node-%d rejected command %d\n", primaryID, 3003+i)
			}
		}()

		go func() {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)
			// re-connect
			cfg.connect((primaryID + 1) % servers)

			if _, _, ok := cfg.pbservers[primaryID].Start(3004 + i); !ok {
				t.Fatalf("node-%d rejected command %d\n", primaryID, 3004+i)
			}
		}()

		wg.Wait()
		cfg.replicateOne(primaryID, 3005, servers)
		// check that all servers replicate the same sequence of commands
		var command interface{}
		for index := 1; index <= 5+i; index++ {
			cfg.checkCommittedIndex(index, command, servers)
		}
//		fmt.Printf("iteration i=%d finished\n", i)
	}

	fmt.Printf("  ... Passed\n")
}

func Test1AFailNoCommitPB(t *testing.T) {
	servers := 3 //3 servers
	primaryID := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(primaryID, 4001, servers)

	// disconnect 2 out of 3 servers, both of which are backups
	cfg.disconnect((primaryID + 1) % servers)
//	fmt.Printf("server %v DISConnected\n",(primaryID + 1) % servers)
	cfg.disconnect((primaryID + 2) % servers)
//	fmt.Printf("server %v DISConnected\n",(primaryID + 2) % servers)
	// try to replicate command 4002
	index, _, ok := cfg.pbservers[primaryID].Start(4002)
	if !ok {
		t.Fatalf("primary rejected the command\n")
	}
	if index != 2 {
		t.Fatalf("expected index 2, got %v\n", index)
	}
	time.Sleep(2 * time.Second)

	committed := cfg.pbservers[primaryID].IsCommitted(index)
	if committed {
		t.Fatalf("index %d is incorrectly considered to have been committed\n", index)
	}

	// reconnect backups
	cfg.connect((primaryID + 1) % servers)
	cfg.connect((primaryID + 2) % servers)
//	fmt.Printf("server %v Connected\n",(primaryID + 1) % servers)
//	fmt.Printf("server %v Connected\n",(primaryID + 2) % servers)
	cfg.replicateOne(primaryID, 4003, servers)
//	fmt.Printf("going to SEND 4004\n")
	index = cfg.replicateOne(primaryID, 4004, servers)

	// disconnect the primary
	cfg.disconnect(primaryID)
//	fmt.Printf("server %v DISConnected\n",(primaryID))
	index2, _, ok := cfg.pbservers[primaryID].Start(4005)
	if !ok {
		t.Fatalf("primary rejected command\n")
	}
	if index2 != (index + 1) {
		t.Fatalf("primary put command at unexpected pos %d\n", index2)
	}
	time.Sleep(2 * time.Second)
	committed = cfg.pbservers[primaryID].IsCommitted(index2)
	if committed {
		t.Fatalf("index %d is incorrectly considered to have been committed\n", index2)
	}

	// reconnect primary
	cfg.connect(primaryID)
//	fmt.Printf("server %v DISConnected\n",(primaryID ))
	cfg.replicateOne(primaryID, 4006, servers)
	cfg.replicateOne(primaryID, 4007, servers)

	fmt.Printf(" ... Passed\n")
}

func Test1BSimpleViewChange(t *testing.T) {
	servers := 3 //3 servers
	oldPrimary := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(oldPrimary, 5001, servers)
	cfg.checkCommittedIndex(1, 5001, servers)

//	fmt.Printf( "disconnect one backup\n")
	transientBackup := (oldPrimary + 1) % servers
	cfg.disconnect(transientBackup)
//	fmt.Printf("replicate 5002 at oldPrimary and the remaining connected backup\n")
	cfg.replicateOne(oldPrimary, 5002, majority(servers))
	cfg.checkCommittedIndex(2, 5002, majority(servers))

//	fmt.Printf( "disconnect oldPrimary %v\n",oldPrimary)
	cfg.disconnect(oldPrimary)

//	fmt.Printf("reconnect the previously disconnected backup %v\n",transientBackup)
	cfg.connect(transientBackup)

//	fmt.Printf(" change to a new view 1\n")
	v1 := 1
	cfg.viewChange(v1)
	newPrimary := GetPrimary(v1, servers)

	cfg.replicateOne(newPrimary, 5003, majority(servers))
	cfg.replicateOne(newPrimary, 5004, majority(servers))

	for i := 1; i <= 4; i++ {
		cfg.checkCommittedIndex(i, 5000+i, majority(servers))
	}

//	fmt.Printf( "try to replicate 10 commands 5002...5011 at old disconnected primary\n")
	for i := 0; i < 10; i++ {
		_, _, ok := cfg.pbservers[oldPrimary].Start(5002 + i)
		if !ok {
			t.Fatalf("old primary %d rejected command\n", oldPrimary)
		}
	}

//	fmt.Printf( "reconnect old primary\n")
	cfg.connect(oldPrimary)

//	fmt.Printf("replicate 5005 through newPrimary to all 3 servers\n")
	cfg.replicateOne(newPrimary, 5005, servers)
//	fmt.Printf(" check that all 5001...5005 have been replicated at the correct place at all servers\n")
	for i := 1; i <= 5; i++ {
		cfg.checkCommittedIndex(i, 5000+i, servers)
	}
}

func Test1BConcurrentViewChange(t *testing.T) {
	servers := 3 //3 servers
	v0Primary := GetPrimary(0, servers)
	cfg := make_config(t, servers, false)
	defer cfg.cleanup()

	cfg.replicateOne(v0Primary, 6001, servers)
	cfg.checkCommittedIndex(1, 6001, servers)

//	fmt.Printf("DISconnect node0\n")
	cfg.disconnect(v0Primary)

//	fmt.Printf("TRY to commit command 6002 through disconnected v0Primary, should not succeed\n")
	index, _, ok := cfg.pbservers[v0Primary].Start(5999)
	if !ok {
		t.Fatalf("primary rejected the command\n")
	}
	if index != 2 {
		t.Fatalf("expected index 2, got %v\n", index)
	}
	time.Sleep(2 * time.Second)
	committed := cfg.pbservers[v0Primary].IsCommitted(2)
	if committed {
		t.Fatalf("index 2 is incorrectly considered to have been committed\n")
	}

//	fmt.Printf( "CONcurrent view change\n")
	var wg sync.WaitGroup
	newView := 2
	for v := 1; v <= newView; v++ {
		wg.Add(1)
		go func(view int) {
			defer wg.Done()
			cfg.viewChange(view)
		}(v)
	}
	wg.Wait()

//	fmt.Printf( "REConnect v0Primary\n")
	cfg.connect(v0Primary)

	newView = 5
	for v := 3; v <= newView; v++ {
		wg.Add(1)
		go func(view int) {
			defer wg.Done()
			cfg.viewChange(view)
		}(v)
	}
	wg.Wait()

	newPrimary := GetPrimary(newView, servers)
	cfg.replicateOne(newPrimary, 6002, servers)

	for i := 1; i <= 2; i++ {
		cfg.checkCommittedIndex(i, 6000+i, servers)
	}
}
