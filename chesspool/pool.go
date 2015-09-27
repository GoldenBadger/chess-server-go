package chesspool

import (
	"bufio"
	"log"
	"os/exec"
	"strings"
	"sync"
)

type Move struct {
	Id       int
	Position string
	Depth    int
	Result   string
}

type Pool struct {
	wait         sync.WaitGroup
	lock         sync.Mutex
	numProcesses int
	done         chan bool
	Input        chan Move
	Output       map[int]Move // Maps Move ID to Move.
}

func NewPool(filename string, num int) *Pool {
	log.Printf("INFO: Starting pool with %d processes.", num)
	p := &Pool{
		numProcesses: num,
		done:         make(chan bool, num),
		Input:        make(chan Move, 256),
		Output:       make(map[int]Move),
	}
	p.wait.Add(num)
	for i := 0; i < num; i++ {
		go p.runEngine(filename)
	}
	return p
}

func (p *Pool) Stop() {
	log.Printf("INFO: Stopping pool.")
	for i := 0; i < p.numProcesses; i++ {
		p.done <- true
	}
	p.wait.Wait()
	log.Printf("INFO: Pool stopped.")
}

func (p *Pool) runEngine(filename string) {
	engine := exec.Command(filename)
	engine.Start()
	go p.engineJob(engine)
	engine.Wait()
	p.wait.Done()
}

func (p *Pool) engineJob(engine *exec.Cmd) {
	in, err := engine.StdinPipe()
	inWriter := bufio.NewWriter(in)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	out, err := engine.StdoutPipe()
	outReader := bufio.NewReader(out)
	if err != nil {
		log.Println("ERROR: ", err)
	}
	inWriter.WriteString("uci\nisready\n")
	select {
	case move := <-p.Input:
		inWriter.WriteString("position fen " + move.Position + "\n")
		inWriter.WriteString("go depth " + string(move.Depth) + "\n")
		engineOutput, err := outReader.ReadString('\n')
		if err != nil {
			log.Println("ERROR: ", err)
		}
		for !strings.HasPrefix(engineOutput, "bestmove") {
			engineOutput, err = outReader.ReadString('\n')
			if err != nil {
				log.Println("ERROR: ", err)
			}
		}
		if len(engineOutput) >= 14 {
			move.Result = strings.TrimSpace(engineOutput[9:14])
		} else if len(engineOutput) == 13 {
			move.Result = strings.TrimSpace(engineOutput[9:13])
		} else {
			log.Println("ERROR: Engine returned invalid move. Output: ", engineOutput)
		}
		p.lock.Lock()
		p.Output[move.Id] = move
		p.lock.Unlock()
	case <-p.done:
		err := engine.Process.Kill()
		if err != nil {
			log.Println("ERROR: Could not kill engine process: ", err)
		}
		return
	}
}
