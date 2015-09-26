package chesspool

import (
	"log"
	"os/exec"
	"sync"
)

type Move struct {
	Id       int
	Position string
	Depth    int
	Result   string
}

type Pool struct {
	wait   sync.WaitGroup
	lock   sync.Mutex
	Input  chan Move
	Output map[int]Move // Maps Move ID to Move.
}

func (p *Pool) NewPool(num int) {

}

func (p *Pool) runEngine(filename string) {
	engine := exec.Command(filename)
	engine.Start()
	defer p.engineJob(engine)
	engine.Wait()
}

func (p *Pool) engineJob(engine *exec.Cmd) {
}
