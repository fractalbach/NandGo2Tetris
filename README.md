# Nand GO2 Tetris Files

This repository contains some of the code I've written
while taking the course **Nand 2 Tetris**.  There are videos online from the [Nand2Tetris Coursera](https://www.coursera.org/learn/build-a-computer), and more information at the [Nand2Tetris Website](www.nand2tetris.org).

The hardware layers are mostly written in a [Hardware
Description Language](https://en.wikipedia.org/wiki/Hardware_description_language).
However, once you reach the [Assembler](https://en.wikipedia.org/wiki/Assembly_language#Assembler), and begin to climb up the Software Hierarchy, it makes more sense to write the programs in a high level language.

I choose to write these tools in the Go Programming Language, simply because it seemed like a good idea at the time.
Also, because it was a lot of fun ;)



# The Programs

- [x] (HACK Assembler)[hackasm] converts Assembly to Machine Code.
- [ ] ***In Progress*** - (HACK VM Translator)[hackvmslate] converts VM code into assembly.
- [ ] (HACK Compiler) converts code in a High Level Programming Language into VM code.



## Assembler: hackasm

hackasm.go is a [Multi-Pass Assembler](https://en.wikipedia.org/wiki/Assembly_language#Number_of_passes),
which builds up a [Symbol table](https://en.wikipedia.org/wiki/Symbol_table) by scanning through 
the source assembly file multiple times.  There are two main types of symbols:  **Variables** and **Locations**.

Locations are used like "waypoints", creating a place to "jump" to (or "goto").  They specify the actual line of machine code that the [Program counter](https://en.wikipedia.org/wiki/Program_counter)  will change to.

Variables store arbitrary values in memory registers that are determined by the assembler.  By convention of this specific assembly language, the variable start at the 17th register, and each new variable is inserted after that (17, 18, 19, ... ).

The Assembler mostly just manipulates bytes and strings. The input is a file written in ASCII, filled with comments and spaces.  the output is a file with only 1s and 0s.  Each instruction is converted into a line of machine code, with a 1-to-1 correspondence.  This makes [Disassembly](https://en.wikipedia.org/wiki/Disassembler) possible, because you can then reverse the process, and retrieve most of the assembly source again.


## VM Translator: hackvmslate

Translates the Hack Virtual Machine language into assembly instructions. 