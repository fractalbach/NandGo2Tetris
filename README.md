# Nand GO2 Tetris Files

This repository contains some of the code I've written
while taking the course **Nand 2 Tetris**.  There are videos online from the [Nand2Tetris Coursera](https://www.coursera.org/learn/build-a-computer), and more information at the [Nand2Tetris Website](www.nand2tetris.org).

The hardware layers are mostly written in a [Hardware
Description Language](https://en.wikipedia.org/wiki/Hardware_description_language).
However, once you reach the [Assembler](https://en.wikipedia.org/wiki/Assembly_language#Assembler), and begin to climb up the Software Hierarchy, it makes more sense to write the programs in a high level language.

I choose to write these tools in the Go Programming Language, simply because it seemed like a good idea at the time.
Also, because it was a lot of fun ;)



# The Programs

- [x] [HACK Assembler](https://github.com/fractalbach/NandGo2Tetris/tree/master/hackasm) converts Assembly to Machine Code.
- [x] [HACK VM Translator](https://github.com/fractalbach/NandGo2Tetris/tree/master/hackvmslate) converts VM code into assembly.
- [ ] (***IN PROGRESS***) -  [HACK Compiler](https://github.com/fractalbach/NandGo2Tetris/tree/master/hackcompiler) converts code in a High Level Programming Language into VM code.
    - [x] Tokenization
    - [x] Semantic Analysis 
    - [ ] Code Generation (***IN PROGRESS***)
