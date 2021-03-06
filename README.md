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
- [x] [HACK Compiler](https://github.com/fractalbach/NandGo2Tetris/tree/master/hackcompiler) converts code in a High Level Programming Language into VM code.
    - [x] Tokenization
    - [x] Semantic Analysis 
    - [x] Code Generation

## Example

Here is what it looks like to compile a program written in the "Jack Programming language" into the VM language, then the Assembly Language, and finally into binary.

![copilers](https://user-images.githubusercontent.com/32124562/41803426-920254d8-763d-11e8-83ea-75515d91e5ed.PNG)
