#GoVirtual

[![Build Status](https://drone.io/github.com/TSavo/GoVirtual/status.png)](https://drone.io/github.com/TSavo/GoVirtual/latest) [![Coverage Status](https://coveralls.io/repos/TSavo/GoVirtual/badge.png)](https://coveralls.io/r/TSavo/GoVirtual)

GoVirtual is a virtual machine toolkit implemented in Go, designed for flexibility and reuse.

Out of the box it provides all the tools necessary to construct a working CPU architecture entirely in software and run programs on that architecture. It can be used as a rapid prototyping environment for CPU architectures, emulation of existing architectures, or as a learning tool for CPU design. 

It includes the necessary code for emulating:
1. Memory (including support for heap and stack based operations)
2. Instruction Sets (as supplied by the user)
3. Programs (including a lexer/compiler to transform your code into programs for your architecture)
4. Processors (which use programs to execute a series of instructions which operate on memory)
5. Termination conditions (for halting the machine)

##Wait, what?

>CPU architecture in software? Why on earth would you want to do such a thing?

First of all, this is coder's code. It's meant for extending, learning, and experimenting. It presumes you more or less understand how CPUs actually work, and want to make one yourself, or emulate an existing one. No apologies will be made for the design, and no claims of performance will be found here. This project is NOT about emulating any one architecture, or doing it the 'best' way, it's about providing a framework for emulating ANY architecture. Much leg work will be involved to actually see it do anything interesting. It will, in short, not change the world any time soon... unless of course you happen to be wanting to get into CPU architecture design using emulation in software, in which case this could rock your world.

But to actually answer the question... besides the obvious answer of, "Designing a CPU architecture is a fascinating exercise that every coder should experience.", it turns out that emulated architectures have a really interesting advantage: *they are written in software*. That means they can be debugged like any other piece of software, they can be designed using all the modern software development tools available, and can take full advantage of the entire technology stack from IDE to optimizing compiler to JIT VMs that optimize their actual flow of execution (VMception?). This means that processor design is no longer limited to people who are working with hardware, and people can experiment with alternative or obsolete and no longer available architectures, implement emulators quickly and easily, and through the power of the Go programming language, parallelize their execution by adding two letters to their program: 'go'.

Because they're written in software, they have another interesting advantage: They can literally do anything. Want a CPU with built in SIMD support? What about having 10,000 of them all linked together and executing in parallel? What about an architecture that changes it's configuration and inner workings during execution? The sky is the limit with regards to what you can make it do, because it's just software.

So, let's get started designing a CPU! 

It all starts with the memory...

##Memory

The purpose of any machine is to perform computations on values. Ultimately those values are stored as a series of bits in 'memory'. 
GoVirtual defines 'memory' as a map of `Address` to `Value`, where `Address` is an empty `interface{}` and Value is defined as:

```go
type Value interface {
	Get() interface{}
	Set(interface{})
}
```

What this means is that you're free to define your memory in terms of any architecture definition you want. 
Traditional architectures typically define their memory in terms of a fixed number of bits at each address space, which we usually call an 'int',
and this is in fact the most obvious implementation, but there's actually nothing stopping us from doing something really outside the box
in terms of declaring our memory addresses to be full objects with invokable interfaces and complex state. After all, the only requirement is that
they are interface{}.

More importantly, because the structure of memory is unbound, we can use this abstraction to define different logical segments of our memory, such as registers, RAM, ROM, and external devices, simply by using a different addressing scheme for each segment. We could, for instance, define the Value at address 'RAM' as an int[], or we could define 






This simplifies the design because all values in memory are also valid pointers to memory (an int can deference a slice by cardinality). This doesn't mean that VMs can't perform floating point operations, it's just that all internal representations are in int format, so any floating point operations will require a conversion from int when loading from memory, and to int when storing.

The Memory class is defined as:

```go
type Memory []int
```

It also has a number of helper methods on it, the two most important of which are .Get(index int) (int),  and .Set(index int, value int). Get and Set are operations that will always succeed, because they always ensure that the index is within the range of the length of the slice, and take the remainder of dividing by the length if it isn't, so they behave like circular buffers which wrap around if you keep incrementing index beyond the length of the slice.

Other helpful methods are .Resize(size int), which increases or decreases the length of the slice (elements are added to or removed from the end), .Zero() which sets all values to zero, .Reallocate(size int) which resizes AND zeros it, Append(), Prepend(), Delete(), and others.

##InstructionSet

An 'instruction set' is a set of operations which the processor is able to perform. Instructions can be defined to be anything at all, and are entirely supplied by you, the user. Since the instruction set is extensible, complex operations can be supported by the addition of an instruction that invokes that complex operation, allowing the processor to interact with the outside world in whatever way user intends.

Instructions come in two flavors: Infix, which can have at most two arguments, and prefix, which can have unlimited arguments.

For example, you might have an infix 'add' instruction that takes two literals, and adds them together:

```
2 add 2
```

Or perhaps, more familarly, the '+' infix instruction:

```
2 + 2
```

Prefix instructions look like function calls, with the arguments being comma seperated within the parentencies. The same add instruction in prefix form would look like:

```
add(2, 2)
```
or
```
+(2, 2)
```



Each instruction in the instruction set has a unique name, an associated function to execute, and a value to advance the processor onto the next instruction after execution called 'movement' (usually 0 in the case of an operation that changes the instruction pointer explicitly like a 'jump', or 1 to just execute the next instruction). The method signature for those functions accepts a pointer to a processor core, and a pointer to some memory for the operands, so it looks like this:

```go
func(*vm.Processor, *vm.Memory)
```

You can define a new instruction by calling .Define(name string, movement int, closure func(*Processor, *Memory)). For instance the definition of a 'do nothing' instruction would look like this:

```go
is := make(vm.InstructionSet)
is.Define("noop", 1, func(*vm.Processor, *vm.Memory){})
```

Here we've defined an instruction called 'noop', which will after execution advance the instruction pointer by 1, and will call the closure in the third argument when executed. 

In that closure the first argument is the processor the instruction is being executed on, and the second argument is the 'operands' which are supplied with the instruction, which when put together is called an 'Operation'. 

There's helper methods on instruction set which do the same thing: .Instuction(...) and .Movement(...). Both functions take a string and a closure, .Instruction() calls .Compile() with 1 movement, and .Movement() calls compile with 0 movement. It's expected that 'movement' operations will change the instruction pointer on the processor on it's own, otherwise the processor will never advance in the program. 

```go
is.Instruction("noop", func(*vm.Processor, *vm.Memory){})
```

So now our instruction set knows what a 'noop' instruction does. But we're still not ready to use it to make a program yet. Do do that, we need to pair it with it's operands, making it an operation.

To assemble an instruction into an operation, we call upon the instruction set to combine an instruction with a set of 0 or more operands using the 'Compile' method, which has the following signature:

```go
func (index *InstructionSet) Compile(name string, args... int) (o *Operation)
```

So to compile that noop instruction into an operation, you would:

```go
op := is.Compile("noop")
```

This step is required because the instructions are singular (only one instruction with the name 'noop' exists), but we need to represent the instructions as a series of operations which is plural. A series of operations is called a 'Program'.

Of course noop requires no operands, so nothing but the name is passed into .Compile(). Let's take a look at a slightly more complex instruction. Let's say you wanted an 'add' instruction which added the contents of one register to another register, and stored it in a third register. First, the definition of the instruction:

```go
addInstruction := is.Instruction("add", func(processor *vm.Processor, memory *vm.Memory) {
	processor.Registers[(*memory)[2]] = processor.Registers[(*memory)[0]] + processor.Registers[(*memory)[1]]) 
})
```

That's the definition. What it does is call Set on the registers with the third value from the operands, where the value to set is the value of the registers (deferenced by the first and second operands) added together.

Now to compile that instruction into an operation which adds the contents of register 0 and 1 together, and stores the result in register 2:

```go
op := is.Compile("add", 0, 1, 2)
```

So if your registers looked like this:

```go
[6, 4, 0]
```

After executing the above operation, they would look like this:

```go
[6, 4, 10]
```

In this way you can make your instruction set do whatever you want, including emulating another architecture. For instance, the 8088 architecture (the father of all modern CPUs) defined the "mov" instruction for moving a value in memory into a register.

>mov eax, [ebx] ; Move the 4 bytes in memory at the address contained in EBX into EAX

That would look like:

```go
is.Instruction("mov", func(processor *vm.Processor, memory *vm.Memory) {
	(*processor).Registers.Set((*memory).Get(0), (*processor).Heap.Get((*memory).Get(1))) 
})
```

And compiled:

```go
movOp = is.Compile("mov", 0, 100) // move the value in heap[100] into register[0]
```
 
##Program

A 'Program' is an array of operations. It's defined as:

```go
type Program []*Operation
```

So to make a simple program that would print, "Hello World!" to the console over and over, we might define two instructions, one to print the message, and one to 'jump' the processor back to the begining of the program:

```go
is := make(vm.InstructionSet)

is.Instruction("printHelloWorld", func(processor *vm.Processor, memory *vm.Memory){
	fmt.Println("Hello World!")
})

is.Movement("jump", 1, func(processor *vm.Processor, memory *vm.Memory){
	processor.Jump((*memory).Get(0))
})

program := make(vm.Program, 2)
program[0] = is.Compile("printHelloWorld")
program[0] = is.Compile("jump", 0)
```

Or you can write the program as a string, and have the instruction set compile it for you into a program. The CompileProgram() method takes a string with one instruction and it's operands per line:

```go
program := is.CompileProgram(`
printHelloWorld
jump 0
`)
```

Operands can be separated with a comma, space, or both. All the following forms are valid:

```go
noop              //No operands
decrement 1       //One operand
set 1, 255        //Operands seperated by comma
jumpIfNotZero 1 0 //Operands seperated by space
```

##Termination Conditions

Processors can (and will) execute 'forever'. If it reaches the end of the program, it will start over at the beginning (Programs act like memory/circular buffers too, so all values of 'instruction pointer' are valid). To get a processor to stop, a Termination Condition interface is defined with a single method, ShouldTerminate(*Processor) bool:

```go
type TerminationCondition interface {
	func ShouldTerminate(*Processor) bool
}
```

ShouldTerminate() gets called after every execution of the program, and if it returns true, the processor stops.

Of course a termination condition that never terminates is trivial to write:

```go
type NeverTerminate struct {}
func (never NeverTerminate) ShouldTerminate(processor *vm.Processor) bool {
	return false
}
```

With this example a good compiler will elide the entire branch. A more interesting termination condition might wait until someone sends on a channel:

```go
type TerminateChan chan bool
func (termChan TerminateChan) ShouldTerminate(processor *vm.Processor) bool {
	select {
	case <- termChan:
		return true
	default:
		return false
	}
}
```

You could then make this channel buffered, so anyone who wrote on it wouldn't block, and the processor will pick it up when it can:

```go
term := make(TerminateChan, 1)
```

##Processors

A processor is composed of an 'instruction set' (the set of all operations supported on the processor), several blocks of 'memory', a program (a list of instructions and operands), an 'instruction pointer' (where in the program the processor is executing), and a 'termination condition' (a object which can tell it when to stop executing).

The processor has support for several blocks of memory: local or 'register' memory, 'heap' memory (which can be shared among processors, if you like), 'stack' memory which support push and pop operations in addition to cardinal access, and a separate 'call stack' memory, all of which can be accessed/modified by the instructions. However unlike most processor designs, the rules for how instructions operate on memory are left wide open, allowing the user to model the operation of the virtual machine however they want.

Support for compilation of textual representations of programs is included, but the supported grammar is limited to a single 'assembly-like' form. This allows users to write programs in separate files, and 'load' them into the processor's program space for execution.

To load a program, you can use .LoadProgram(program *Program):

```go
program := is.CompileProgram("...")

processor := make(vm.Processor)
processor.LoadProgram(program)
```

Don't forget your termination condition:

```go
term := make(SomeTerminationCondition)
processor.TerminationCondition = &term
```

Give it some memory:

```go
processor.Registers = make(vm.Memory, 4)
```

And then start the program running with .Run():

```go
processor.Run()
```

Run() will enter a for loop which will keep executing your program until your termination condition returns true, so if you want to continue execution on this thread, add 'go' before it:

```go

go processor.Run()
```

So, let's put it all together. We'll make a program that executes 100 operations, and a program that computes a power of 2 series:

```go
//Define our Termination condition
type GoForAWhileThenStop struct {
	times int                     //State for our condition
}

//The TerminateCondition has a single method we need to implement to conform to the interface.
func (self *GoForAWhileThenStop) ShouldTerminate(processor *vm.Processor) bool {
	if(self.times >= 100) { //Have we been called more than 100 times?
		return true         //Stop the processor
	}
	self.times++            //Increment how many times we've been called
	return false            //Don't stop yet
}

//Our Terminate condition
term := make(GoForAWhileThenStop)

//Our instruction set
is := make(vm.InstructionSet) 

//Our first instruction writes a value to a register
is.Instruction("setRegister", func(processor *vm.Processor, memory *vm.Memory){
	processor.Registers[(*memory)[0]] = (*memory)[1]
})

//Our second instruction multiplies two registers, and sets the result back to the first register x86 style.
is.Instruction("multiply", func(processor *vm.Processor, memory *vm.Memory){
	processor.Registers[(*memory)[0]] = processor.Registers[(*memory)[0]] * processor.Registers[(*memory)[1]]
})

//Our third instruction changes the instruction pointer.  
is.Movement("jump", 1, func(processor *vm.Processor, memory *vm.Memory){
	processor.Jump((*memory)[0])
})

//Our program:
//  register[0] = 1
//  register[1] = 2
//:loop
//  register[0] = register[0] * register[1]
//  goto :loop
//
//Now compile it.
program := is.CompileProgram(`
set 0 1
set 1 2
multiply 0 1
jump 2
`)

//Our Processor
processor := make(vm.Processor)
processor.Registers = make(vm.Memory, 4)
//Set our terminate condition
processor.TerminateCondition = &term
//Load our program
processor.LoadProgram(program)
//And away we go!
processor.Run()
//When we reach here in our code, it's because the terminate condition returned true.
//Let's do it again, but this time let's not wait for it to finish.
processor.Reset()
go processor.Run()
```

###Credit
The project was originally forked from https://github.com/feyeleanor/GoLightly, but has been changed drastically enough to warrant a separate project and name.

