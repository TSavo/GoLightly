#GoVirtual

GoVirtual is a lightweight virtual machine toolkit implemented in Go, designed for flexibility and reuse.

##Processors

GoVirtual provides a general purpose 'Processor' which supports a user extensible instruction set. An instruction set is a set of instructions which the processor is able to execute.

A processor is composed of an 'instruction set' (the set of all operations supported on the processor), several blocks of 'memory', a program (a list of instructions and operands), an 'instruction pointer' (where in the program the processor is executing), and a 'termination condition' (a object which can tell it when to stop executing).

Memory is defined as an array of 'int's, with several methods to ensure it's safe operation. Memory can be grown or shrunk dynamically, and all read/write operations are abs(index) % len(memory), so it's impossible to 'go outside the bounds' of a block of memory, since indexes greater than the length will just wrap around to the start of the memory. Memory of 'size == 0' is not supported.

The processor has support for several blocks of memory: local or 'register' memory, 'heap' memory (which can be shared among processors, if you like), 'stack' memory which support push and pop operations in addition to cardinal access, and a separate 'call stack' memory, all of which can be accessed/modified by the instructions. However unlike most processor designs, the rules for how instructions operate on memory are left wide open, allowing the user to model the operation of the virtual machine however they want.

Support for compilation of textual representations of programs is included, but the supported grammar is limited to a single 'assembly-like' form. This allows users to write programs in separate files, and 'load' them into the processor's program space for execution.

##InstructionSet

Since the instruction set is extensible, complex operations can be supported by the addition of an instruction by the user, allowing the Processor to interact with the outside world in whatever way user intends.



Firstly GoLightly provides support for vector instructions allowing each virtual processor to handle
large data sets more efficiently than in traditional VM designs. Not only is this a boon for fast
maths, it also allows complex data structures to be transferred between virtual processors with only
a couple of instructions.

And speaking of multiple virtual processors, once a processor has been initialised it can be cloned
with a single instruction and both processors are hooked together with a communications channel.

Each processor has its own dynamic instruction definition and dispatch mechanism providing a flexible
model during development. However function-based dynamic dispatch is more expensive than the
switch-based dispatch used in traditional VM designs so the latter is also supported via type
embedding of the ProcessorCore and method hiding of InlinedInstructions().

One interesting consequence of making each virtual processor configurable in this manner is that it
allows complex problems to be modelled in terms of several different processors running relatively
simple programs and arranged in a pipeline.

Likewise the ability to clone virtual processors can be used as a building block for transactional
and speculative execution.


== Tests and Benchmarks ==

Golightly includes a large collection of tests and benchmarks. A simple set of micro-benchmarks
is available to judge the relative cost of various standard Go instructions and idioms which can
be invoked from the command-line:

	cd go
	gotest -benchmarks="Benchmark"


Likewise the test suite for the VM primitives can be executed with:

	cd vm
	gotest

and the optional benchmarks invoked with:

	gotest -benchmarks="Benchmark"

These include extensive micro-benchmarks focusing on individual primitive operation performance which
will hopefully act as a useful guide to relative performance for those writing compilers or assemblers
targeting GoLightly virtual machines.