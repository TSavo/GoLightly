package govirtual

//import "fmt"

func EmulationInstructions(i *InstructionSet) {
	i.Instruction("noop", func(p *Processor, args ...Pointer) Memory {
		return nil
	})
	//i.Instruction("print", func(p *Processor, args ...Pointer) Memory {
//		fmt.Println(args[0].Get())
//		return nil
//	})
	i.Movement("jump", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		p.Jump(args[0])
		return nil
	}, Argument{"label", "string"})
	i.Movement("jumpIfZero", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		if Cardinalize(args[0].Get()) == 0 {
			p.Jump(args[1].Get())
		} else {
			p.InstructionPointer++
		}
		return nil
	}, Argument{"compare", "int"}, Argument{"label", "string"})
	i.Movement("jumpIfNotZero", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		if Cardinalize(args[0].Get()) != 0 {
			p.Jump(args[1].Get())
		} else {
			p.InstructionPointer++
		}
		return nil
	}, Argument{"compare", "int"}, Argument{"label", "string"})
	i.Movement("jumpIfEquals", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		if Cardinalize(args[0].Get()) == Cardinalize(args[1].Get()) {
			p.Jump(args[2].Get())
		} else {
			p.InstructionPointer++
		}
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"label", "string"})
	i.Movement("jumpIfNotEquals", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		if Cardinalize(args[0].Get()) != Cardinalize(args[1].Get()) {
			p.Jump(args[2].Get())
		} else {
			p.InstructionPointer++
		}
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"label", "string"})
	i.Movement("jumpIfGreaterThan", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		if Cardinalize(args[0].Get()) > Cardinalize(args[1].Get()) {
			p.Jump(args[2].Get())
		} else {
			p.InstructionPointer++
		}
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"label", "string"})
	i.Movement("jumpIfLessThan", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		if Cardinalize(args[0].Get()) < Cardinalize(args[1].Get()) {
			p.Jump(args[2].Get())
		} else {
			p.InstructionPointer++
		}
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"label", "string"})
	i.Movement("call", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		p.Call(args[0].Get())
		return nil
	}, Argument{"label", "string"})
	i.Movement("return", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			if recover() != nil {
				p.InstructionPointer++
			}
		}()
		p.Return()
		return nil
	})
	i.Instruction("add", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			recover()
		}()
		args[2].Set(Cardinalize(args[0]) + Cardinalize(args[1]))
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"result", "int"})
	i.Instruction("subtract", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			recover()
		}()
		args[2].Set(Cardinalize(args[0]) - Cardinalize(args[1]))
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"result", "int"})
	i.Instruction("multiply", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			recover()
		}()
		args[2].Set(Cardinalize(args[0]) * Cardinalize(args[1]))
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"result", "int"})
	i.Instruction("divide", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			recover()
		}()
		args[2].Set(Cardinalize(args[0]) / Cardinalize(args[1]))
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"result", "int"})
	i.Instruction("modulos", func(p *Processor, args ...Pointer) Memory {
		defer func(){
			recover()
		}()
		args[2].Set(Cardinalize(args[0]) % Cardinalize(args[1]))
		return nil
	}, Argument{"left", "int"}, Argument{"right", "int"}, Argument{"result", "int"})

//	i.Infix("=", func(p *Processor, args ...Pointer) Memory {
//		args[0].Set(args[1].Get())
//		return nil
//	}, Argument{"left", "int"}, Argument{"right", "int"})

}
