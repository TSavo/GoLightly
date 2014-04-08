package main

//import (
//	"fmt"
//	"github.com/gorilla/websocket"
//	"github.com/tsavo/golightly/intutil"
//	"github.com/tsavo/golightly/vm"
//	"io"
//	"net"
//	"net/http"
//	"time"
//)
//
//const (
//	POPULATION_SIZE = 10000
//	BEST_OF_BREED   = 1000
//	PROGRAM_LENGTH  = 50
//	UNIVERSE_SIZE   = 9
//	ROUND_LENGTH    = 10
//)
//
//func DefineInstructions() (i *vm.InstructionSet) {
//	i = vm.NewInstructionSet()
//	i.Operator("noop", func(p *vm.ProcessorCore, m *vm.Memory) {
//	})
//	i.Movement("halt", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Running = false
//	})
//	i.Movement("jump", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Jump(p.Registers.Get(1))
//	})
//	i.Movement("jumpIfZero", func(p *vm.ProcessorCore, m *vm.Memory) {
//		if p.Registers.Get((*m).Get(0)) == 0 {
//			p.Jump(p.Registers.Get(1))
//		}
//	})
//	i.Movement("jumpIfNotZero", func(p *vm.ProcessorCore, m *vm.Memory) {
//		if p.Registers.Get((*m).Get(0)) != 0 {
//			p.Jump(p.Registers[1])
//		}
//	})
//	i.Movement("call", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Call(p.Registers.Get((*m).Get(0)))
//	})
//	i.Movement("return", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Return()
//	})
//	i.Operator("set", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), (*m).Get(1))
//	})
//	i.Operator("store", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Heap.Set(p.Registers.Get(1), p.Registers.Get(0))
//	})
//	i.Operator("load", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set(0, p.Heap.Get(p.Registers.Get(1)))
//	})
//	i.Operator("swap", func(p *vm.ProcessorCore, m *vm.Memory) {
//		x := p.Registers.Get((*m).Get(0))
//		p.Registers.Set((*m).Get(0), (*m).Get(1))
//		p.Registers.Set((*m).Get(1), x)
//	})
//	i.Operator("push", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Stack.Push(p.Registers.Get((*m).Get(0)))
//	})
//	i.Operator("pop", func(p *vm.ProcessorCore, m *vm.Memory) {
//		if x, err := p.Stack.Pop(); !err {
//			p.Registers.Set((*m).Get(0), x)
//		}
//	})
//	i.Operator("increment", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Increment((*m).Get(0))
//	})
//	i.Operator("decrement", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Decrement((*m).Get(0))
//	})
//	i.Operator("add", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))+p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("subtract", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))-p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("multiply", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))*p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("divide", func(p *vm.ProcessorCore, m *vm.Memory) {
//		defer func() {
//			recover()
//		}()
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))/p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("modulos", func(p *vm.ProcessorCore, m *vm.Memory) {
//		defer func() {
//			recover()
//		}()
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))%p.Registers.Get((*m).Get(1)))
//	})
//
//	i.Operator("and", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))&p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("or", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))|p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("xor", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))^p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("xand", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))&^p.Registers.Get((*m).Get(1)))
//	})
//	i.Operator("leftShift", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))<<uint(p.Registers.Get((*m).Get(1))))
//	})
//	i.Operator("rightShift", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.Registers.Set((*m).Get(0), p.Registers.Get((*m).Get(0))>>uint(p.Registers.Get((*m).Get(1))))
//	})
//	i.Operator("mutateFaster", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.ChanceOfMutation += 0.01
//	})
//	i.Operator("mutateSlower", func(p *vm.ProcessorCore, m *vm.Memory) {
//		p.ChanceOfMutation -= 0.01
//		if p.ChanceOfMutation < 0.01 {
//			p.ChanceOfMutation = 0.01
//		}
//	})
//
//	return
//}
//
//func evaluate(p *vm.Memory, cost int) int64 {
//	if p.Get(0) > 10000 || p.Get(0) < 0 ||
//		p.Get(1) > 10000 || p.Get(1) < 0 ||
//		p.Get(2) > 10000 || p.Get(2) < 0 ||
//		p.Get(3) > 10000 || p.Get(3) < 0 {
//		return 10000
//	}
//
//	var fit int64 = intutil.Abs64(255 - int64(p.Get(0)))
//	fit += intutil.Abs64(500 - int64(p.Get(1)))
//	fit += intutil.Abs64(255 - int64(p.Get(2)))
//	fit += intutil.Abs64(255 - int64(p.Get(3)))
//	return intutil.Abs64(int64(cost) + fit)
//}
//
//func CombinePrograms(s1 []*vm.Program, s2 []*vm.Program) []*vm.Program {
//	prog := []*vm.Program{}
//	prog = append(prog, s1[0:len(s1)/2]...)
//	prog = append(prog, s2[0:len(s2)/2]...)
//	return prog
//}

//type Member struct {
//	Solver       *vm.Solver
//	Fitness      int64
//	SolutionChan chan *vm.Solution
//	ControlChan  chan bool
//}
//
//type Population []*Member
//
//func handleConnection(conn *net.Conn, solutionChan *chan *vm.Solution) {
//	for {
//		var one []byte
//		solution := <-(*solutionChan)
//		(*conn).SetReadDeadline(time.Now())
//		if _, err := (*conn).Read(one); err == io.EOF {
//			fmt.Println("%s detected closed LAN connection")
//			(*conn).Close()
//			(*conn) = nil
//			break
//		} else {
//			//var zero time.Time
//			(*conn).SetReadDeadline(time.Time{})
//		}
//		for i, x := range solution.Champions {
//			fmt.Fprintf((*conn), "%d,%d,%d\n", solution.Id, i, x.Hashcode())
//		}
//	}
//}
//
//type hub struct {
//	// Registered connections.
//	connections map[*connection]bool
//
//	// Inbound messages from the connections.
//	broadcast chan []byte
//
//	// Register requests from the connections.
//	register chan *connection
//
//	// Unregister requests from connections.
//	unregister chan *connection
//}
//
//var h = hub{
//	broadcast:   make(chan []byte),
//	register:    make(chan *connection),
//	unregister:  make(chan *connection),
//	connections: make(map[*connection]bool),
//}
//
//func (h *hub) run() {
//	for {
//		select {
//		case c := <-h.register:
//			h.connections[c] = true
//		case c := <-h.unregister:
//			delete(h.connections, c)
//			close(c.send)
//		case m := <-h.broadcast:
//			for c := range h.connections {
//				select {
//				case c.send <- m:
//				default:
//					delete(h.connections, c)
//					close(c.send)
//					go c.ws.Close()
//				}
//			}
//		}
//	}
//}
//
//type connection struct {
//	// The websocket connection.
//	ws *websocket.Conn
//
//	// Buffered channel of outbound messages.
//	send chan []byte
//}
//
//func (c *connection) reader(incoming chan string) {
//	for {
//		_, message, err := c.ws.ReadMessage()
//		if err != nil {
//			break
//		}
//		incoming <- string(message)
//		//h.broadcast <- message
//	}
//	c.ws.Close()
//}
//
//func (c *connection) writer() {
//	for message := range c.send {
//		err := c.ws.WriteMessage(websocket.TextMessage, message)
//		if err != nil {
//			break
//		}
//	}
//	c.ws.Close()
//}
//
//func wsHandler(w http.ResponseWriter, r *http.Request) {
//	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
//	if _, ok := err.(websocket.HandshakeError); ok {
//		http.Error(w, "Not a websocket handshake", 400)
//		return
//	} else if err != nil {
//		return
//	}
//	c := &connection{send: make(chan []byte, 256), ws: ws}
//	h.register <- c
//	defer func() { h.unregister <- c }()
//	go c.writer()
//
//	c.reader(nil)
//}
//
func main() {
//	go h.run()
//	vis := make(chan *vm.Solution)
//
//	go func() {
//		http.HandleFunc("/ws", wsHandler)
//		if err := http.ListenAndServe(":3000", nil); err != nil {
//			fmt.Println("ListenAndServe:", err)
//		}
//	}()
//
//	go func() {
//		ln, err := net.Listen("tcp", ":8080")
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//		for {
//			conn, err := ln.Accept()
//			if err != nil {
//				// handle error
//				continue
//			}
//			go handleConnection(&conn, &vis)
//		}
//	}()
//
//	instructionSet := DefineInstructions()
//	//finish := make(chan *vm.ProcessorCore)
//	//results := make(chan int)
//	population := make(Population, 0)
//	for x := 0; x < UNIVERSE_SIZE; x++ {
//		solver := vm.NewSolver(x, POPULATION_SIZE, BEST_OF_BREED, PROGRAM_LENGTH, 4, 4, 0.1, evaluate, instructionSet)
//		solutionChan := make(chan *vm.Solution)
//		control := make(chan bool)
//		population = append(population, &Member{solver, 10000, solutionChan, control})
//		outChan := make(chan int)
//		go solver.Solve(solutionChan, outChan, control, nil)
//	}
//
//	count := 0
//	bestFit := int64(100000)
//	for {
//		count++
//		champs := make([]*vm.Program, 0)
//		for _, member := range population {
//			solution := <-member.SolutionChan
//			select {
//			case vis <- solution:
//			default:
//			}
//			if bestFit > solution.Fitness {
//				fmt.Printf("%d: %d", solution.Id, solution.Fitness)
//				h.broadcast <- []byte(fmt.Sprintf("%d: %d", solution.Id, solution.Fitness))
//				bestFit = solution.Fitness
//				fmt.Println(solution.Heaps[0])
//			}
//			if count < ROUND_LENGTH {
//				member.ControlChan <- true
//			} else {
//				member.ControlChan <- false
//				champs = append(champs, solution.Champions[0:len(solution.Champions)/10]...)
//			}
//		}
//		if count < ROUND_LENGTH {
//			continue
//		}
//		fmt.Printf("recombo: %d\n", len(champs))
//		count = 0
//		population = make(Population, UNIVERSE_SIZE)
//		for x := 0; x < UNIVERSE_SIZE; x++ {
//			solver := vm.NewSolver(x, POPULATION_SIZE, BEST_OF_BREED, PROGRAM_LENGTH, 4, 4, 0.1, evaluate, instructionSet)
//			solutionChan := make(chan *vm.Solution)
//			control := make(chan bool)
//			outChan := make(chan int)
//			population[x] = &Member{solver, 10000, solutionChan, control}
//			go solver.Solve(solutionChan, outChan, control, champs)
//		}
//	}
}
