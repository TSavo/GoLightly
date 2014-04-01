package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tsavo/golightly/intutil"
	"github.com/tsavo/golightly/vm"
	"io"
	"math/rand"
	"net"
	"net/http"
	"runtime"
	"time"
)

const (
	POPULATION_SIZE = 100
	BEST_OF_BREED   = 10
	PROGRAM_LENGTH  = 50
	UNIVERSE_SIZE   = 9
	ROUND_LENGTH    = 10
)

func DefineInstructions(flapChan chan bool) (i *vm.InstructionSet) {
	i = vm.NewInstructionSet()
	i.Operator("noop", func(p *vm.ProcessorCore, m *vm.Memory) {
	})
	i.Movement("jump", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Jump(p.Registers.Get(1))
	})
	i.Movement("jumpIfZero", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) == 0 {
			p.Jump(p.Registers.Get(1))
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfNotZero", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) != 0 {
			p.Jump(p.Registers[1])
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfEquals", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) == p.Registers.Get((*m).Get(1)) {
			p.Jump(p.Registers[1])
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfNotEquals", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) != p.Registers.Get((*m).Get(1)) {
			p.Jump(p.Registers[1])
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfGreaterThan", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) > p.Registers.Get((*m).Get(1)) {
			p.Jump(p.Registers[1])
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("jumpIfLessThan", func(p *vm.ProcessorCore, m *vm.Memory) {
		if p.Registers.Get((*m).Get(0)) < p.Registers.Get((*m).Get(1)) {
			p.Jump(p.Registers[1])
		} else {
			p.InstructionPointer++
		}
	})
	i.Movement("call", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Call(p.Registers.Get((*m).Get(0)))
	})
	i.Movement("return", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Return()
	})
	i.Operator("set", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), (*m).Get(1))
	})
	i.Operator("store", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Heap.Set(p.Registers.Get(1), p.Registers.Get(0))
	})
	i.Operator("load", func(p *vm.ProcessorCore, m *vm.Memory) {
		//fmt.Println(p.Heap.Get(p.Registers.Get(1)))
		p.Registers.Set(0, p.Heap.Get(p.Registers.Get(1)))
	})
	i.Operator("swap", func(p *vm.ProcessorCore, m *vm.Memory) {
		x := p.Registers.Get((*m).Get(0))
		p.Registers.Set((*m).Get(0), (*m).Get(1))
		p.Registers.Set((*m).Get(1), x)
	})
	i.Operator("push", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Stack.Push(p.Registers.Get((*m).Get(0)))
	})
	i.Operator("pop", func(p *vm.ProcessorCore, m *vm.Memory) {
		if x, err := p.Stack.Pop(); !err {
			p.Registers.Set((*m).Get(0), x)
		}
	})
	i.Operator("increment", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Increment((*m).Get(0))
	})
	i.Operator("decrement", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Decrement((*m).Get(0))
	})
	i.Operator("add", func(p *vm.ProcessorCore, m *vm.Memory) {
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))+p.Registers.Get((*m).Get(1)))
	})
	i.Operator("subtract", func(p *vm.ProcessorCore, m *vm.Memory) {
		//fmt.Println(p.Registers.Get((*m).Get(0))-p.Registers.Get((*m).Get(1)))
		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))-p.Registers.Get((*m).Get(1)))
	})
	i.Operator("flap", func(p *vm.ProcessorCore, m *vm.Memory) {
		flapChan <- true
		time.Sleep(50 * time.Millisecond)
	})
	i.Operator("sleep", func(p *vm.ProcessorCore, m *vm.Memory) {
		time.Sleep(50 * time.Millisecond)
	})

	return
}

func CombinePrograms(s1 []*vm.Program, s2 []*vm.Program) []*vm.Program {
	prog := []*vm.Program{}
	prog = append(prog, s1[0:len(s1)/2]...)
	prog = append(prog, s2[0:len(s2)/2]...)
	return prog
}

func handleConnection(conn *net.Conn, solutionChan *chan *vm.Solution) {
	for {
		var one []byte
		(*conn).SetReadDeadline(time.Now())
		if _, err := (*conn).Read(one); err == io.EOF {
			fmt.Println("%s detected closed LAN connection")
			(*conn).Close()
			(*conn) = nil
			break
		} else {
			//var zero time.Time
			(*conn).SetReadDeadline(time.Time{})
		}
	}
}

type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			delete(h.connections, c)
			close(c.send)
		}
	}
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

func (c *connection) reader(incoming chan string) {
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		incoming <- string(message)
		//h.broadcast <- message
	}
	c.ws.Close()
}

func (c *connection) writer() {
	for message := range c.send {
		err := c.ws.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
	c.ws.Close()
}

type FlappyEvaluator struct {
	reward int64
}

func (eval *FlappyEvaluator) Evaluate(*vm.ProcessorCore) int64 {
	x := eval.reward
	eval.reward = 0
	return x
}

type FlappyGenerator struct {
	InstructionSet *vm.InstructionSet
}



func (gen *FlappyGenerator) GenerateProgram() *vm.Program {
	pro := make(vm.Program, 11)
	if rand.Int()%100 > 50 {
		pro[0] = gen.InstructionSet.Compile("set", 4, 0)
		pro[1] = gen.InstructionSet.Compile("set", 2, 3)
		pro[2] = gen.InstructionSet.Compile("set", 1, 5)
		pro[3] = gen.InstructionSet.Compile("set", 3, rand.Int()%2000)
		pro[4] = gen.InstructionSet.Compile("load")
		pro[5] = gen.InstructionSet.Compile("subtract", 3, 0)
		pro[6] = gen.InstructionSet.Compile("set", 1, 0)
		pro[7] = gen.InstructionSet.Compile("jumpIfGreaterThan", 3, 1)
		pro[8] = gen.InstructionSet.Compile("flap")
		pro[9] = gen.InstructionSet.Compile("sleep")
		pro[10] = gen.InstructionSet.Compile("jump")
	} else {
		for x := 0; x < len(pro); x++ {
			pro[x] = gen.InstructionSet.Encode(&vm.Memory{rand.Int() % 2000, rand.Int() % 2000, rand.Int() % 2000})
		}
	}

	return &pro
}

var id = 0

func wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	outChan := make(chan bool, 1)
	flappyEval := new(FlappyEvaluator)
	is := DefineInstructions(outChan)
	flappyGen := &FlappyGenerator{is}
	solver := vm.NewSolver(id, POPULATION_SIZE, BEST_OF_BREED, 4, 0.1, is, flappyGen, flappyEval)
	id++
	go func() {
		for {
			<-outChan
			go func() {
				c.send <- []byte("1")
			}()
		}
	}()
	inFlap := make(chan string, 1000)
	heap := make(vm.Memory, 8)
	control := make(chan bool, 1)
	go func() {
		for {
			flap := <-inFlap
			if flap == "DEAD" {
				control <- false
				continue
			}
			myX := 0
			y := 0
			center := 0
			fmt.Sscanf(flap, "%d,%d,%d", &myX, &y, &center)
			heap.Set(5, myX)
			heap.Set(6, y)
			heap.Set(7, center)
			flappyEval.reward += int64(1000 - intutil.Abs(1000-myX))
		}
	}()
	stopChan := make(chan bool, 1)
	go func() {
		solver.SolveOneAtATime(&heap, nil, nil, control, stopChan, nil, nil)
	}()

	c.reader(inFlap)
	stopChan <- false
}

func main() {
	runtime.LockOSThread()
	go h.run()
	vis := make(chan *vm.Solution)

	go func() {
		http.HandleFunc("/ws", wsHandler)
		if err := http.ListenAndServe(":3000", nil); err != nil {
			fmt.Println("ListenAndServe:", err)
		}
	}()

	go func() {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			fmt.Println(err)
			return
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				// handle error
				continue
			}
			go handleConnection(&conn, &vis)
		}
	}()

	<-make(chan int)
}
