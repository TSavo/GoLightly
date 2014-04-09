#GoVirtual

GoVirtual is a lightweight virtual machine toolkit implemented in Go, designed for flexibility and reuse.

It provides, out of the box, all the tools necessary to construct a working CPU architecture entirely in software, and run programs on that architecture. It can be used as a rapid prototyping environment for CPU architectures, emulation of existing architectures, or as a learning tool for CPU design. 

This includes
-Memory (including support for heap and stack based operations)
-Instruction Sets (defined by the user)
-Programs (and a lexer/compiler to transform your code into programs)
-Processors (which use programs to execute a series of instructions which operate on memory)
-Termination conditions (for halting the machine)

##Memory

The purpose of any machine is to perform computations on values. Ultimately, those values are stored as a series of bits in 'memory'. GoVirtual defines 'memory' as an array of integers ([]int). This simplifies the design because all values in memory are also valid pointers to memory (an int can deference a slice by cardinality). This doesn't mean that VMs can't perform floating point operations, it's just that all internal representations are in int format, so any floating point operations will require a conversion from int when loading from memory, and to int when storing... or separate storage (think map to struct keyed on processor).

The Memory class is defined as:

```go
type Memory []int
```

It also has a number of helper methods on it, the two most important of which are .Get(index int) (int),  and .Set(index int, value int) (void). Get and Set are operations that will always succeed, because they always ensure that the index is within the range of the length of the slice, and take the remainder of dividing by the length if it isn't, so they behave like circular buffers which wrap around if you keep incrementing index beyond the length of the slice.

Other helpful methods are .Resize(size int), which increases or decreases the length of the slice (elements are added to or removed from the end), .Zero() which sets all values to zero, .Reallocate(size int) which resizes AND zeros it, Append(), Prepend(), Delete(), and others.

##InstructionSet

An 'instruction set' is a set of operations which the able processor is able to perform. Instructions can be defined to be anything at all, and are entirely defined by the user. Since the instruction set is extensible, complex operations can be supported by the addition of an instruction by the user that invokes that complex operation, allowing the Processor to interact with the outside world in whatever way user intends.

Each instruction in the instruction set has a unique name, an associated function to execute, and a value to advance the processor onto the next instruction after execution called 'movement' (usually 0 in the case of an operation that changes the instruction pointer explicitly like a 'jump', or 1 to just execute the next instruction). The method signature for those functions accepts a pointer to a processor core, and a pointer to some memory for the operands, so it looks like this:

```go
func(*vm.ProcessorCore, *vm.Memory)
```

You can define a new instruction by calling .Define(name string, movement int, closure func(*ProcessorCore, *Memory)). For instance the definition of a 'do nothing' instruction would look like this:

```go
is := make(vm.InstructionSet)
is.Define("noop", 1, func(*vm.ProcessorCore, *vm.Memory){})
```

Here we've defined an instruction called 'noop', which will after execution advance the instruction pointer by 1, and will call the closure in the third argument when executed. 

In that closure, the first argument is the processor the instruction is being executed on, and the second argument is the 'operands' which are supplied with the instruction, which when put together is called an 'Operation'. 

There's helper methods on instruction set which do the same thing: .Standard(...) and .Movement(...). Both functions take a string and a closure, Standard() calls .Compile() with 1 movement, and .Movement() calls compile with 0 movement. It's expected that 'movement' operations will change the instruction pointer on the processor on it's own, otherwise the processor will never advance in the program. 

```go
is.Standard("noop", func(*vm.ProcessorCore, *vm.Memory){})
```

To assemble an instruction into an operation, we call upon the instruction set to combine an instruction with a set of operands using the 'Compile' method, which has the following signature:

```go
func (i *InstructionSet) Compile(name string, args ...int) (o *Operation)
```

So to compile that noop instruction into an operation, you would:

```go
op := is.Compile("noop")
```

Of course noop requires no operands, so nothing but the name is passed into .Compile(). Let's take a look at a slightly more complex instruction. Let's say you wanted an 'add' instruction which added the contents of one register to another register, and stored it in a third register. First, the definition of the instruction:

```go
addInstruction := is.Standard("add", func(processor *vm.Processor, memory *vm.Memory) {
	(*processor).Registers.Set((*memory).Get(2), (*processor).Registers.Get((*memory).Get(0)) + (*processor).Registers.Set((*memory).Get(1))) 
})
```

That's the definition. What it does is call Set on the registers with the third value from the operands, where the value to set is the value of registers as referenced in the operands.

Now to compile that instruction into an operation which adds the contents of register 0 and 1 together, and stores the result in register 2:

```go
op := is.Compile("add", 0, 1, 2)
```

In this way you can make your instruction set do whatever you want, including emulating another architecture. For instance, the 8088 architecture (the father of all modern CPUs) defined the "mov" instruction for moving a value in memory into a register.

>mov eax, [ebx] ; Move the 4 bytes in memory at the address contained in EBX into EAX

That would look like:

```go
movInstruction := is.Standard("mov", func(processor *vm.Processor, memory *vm.Memory) {
	(*processor).Registers.Set((*memory).Get(0), (*processor).Heap.Get((*memory).Get(1))) 
})

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

is.Standard("printHelloWorld", func(processor *vm.ProcessorCore, memory *vm.Memory){
	fmt.Println("Hello World!")
})

is.Movement("jump", 1, func(processor *vm.ProcessorCore, memory *vm.Memory){
	processor.Jump((*memory).Get(0))
})

program := make(vm.Program, 2)
program[0] = is.Compile("printHelloWorld")
program[0] = is.Compile("jump", 0)
```

Or you can write the program as a string, and have the instruction set compile it for you into a program. The CompileProgram() method takes a string with one instruction and it's operands per line:

```go
program := is.CompileProgram("printHelloWorld\njump 0\n")
```

Operands can be separated with a comma, space, or both. All the following forms are valid:

```go
noop              //No operands
decrement 1       //One operand
set 1, 255        //Operands seperated by comma
jumpIfNotZero 1 0 //Operands seperated by space
```

##Termination Conditions

Processors can (and will) execute 'forever'. If it reaches the end of the program, it will start over at the beginning (Programs act like memory/circular buffers too, so all values of 'instruction pointer' are valid). To get a processor to stop, a Termination Condition interface is defined with a single method, ShouldTerminate(*ProcessorCore) bool:

```go
type TerminationCondition interface {
	func ShouldTerminate(*ProcessorCore) bool
}
```

Of course a termination condition that never terminates is trivial to write:

```go
type NeverTerminate struct {}
func (never NeverTerminate) ShouldTerminate(processor *vm.ProcessorCore) bool {
	return false
}
```

A more interesting termination condition might wait until someone sends on a channel:

```go
type TerminateChan chan bool
func (termChan TerminateChan) ShouldTerminate(processor *vm.ProcessorCore) bool {
	select {
	case <- termChan:
		return true
	default:
		return false
	}
}
```


##Processors

A processor is composed of an 'instruction set' (the set of all operations supported on the processor), several blocks of 'memory', a program (a list of instructions and operands), an 'instruction pointer' (where in the program the processor is executing), and a 'termination condition' (a object which can tell it when to stop executing).

The processor has support for several blocks of memory: local or 'register' memory, 'heap' memory (which can be shared among processors, if you like), 'stack' memory which support push and pop operations in addition to cardinal access, and a separate 'call stack' memory, all of which can be accessed/modified by the instructions. However unlike most processor designs, the rules for how instructions operate on memory are left wide open, allowing the user to model the operation of the virtual machine however they want.

Support for compilation of textual representations of programs is included, but the supported grammar is limited to a single 'assembly-like' form. This allows users to write programs in separate files, and 'load' them into the processor's program space for execution.

To load a program, you can use .LoadProgram(program *Program):

```go
processor 
```

###Credit
The project was originally forked from https://github.com/feyeleanor/GoLightly, but has been changed drastically enough to warrant a separate project and name.

