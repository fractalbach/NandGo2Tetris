Hack Assembler
===========================

hackasm.go is a [Multi-Pass Assembler](https://en.wikipedia.org/wiki/Assembly_language#Number_of_passes),
which builds up a [Symbol table](https://en.wikipedia.org/wiki/Symbol_table) by scanning through 
the source assembly file multiple times.  There are two main types of symbols:  **Variables** and **Locations**.

Locations are used like "waypoints", creating a place to "jump" to (or "goto").  They specify the actual line of machine code that the [Program counter](https://en.wikipedia.org/wiki/Program_counter)  will change to.

Variables store arbitrary values in memory registers that are determined by the assembler.  By convention of this specific assembly language, the variable start at the 17th register, and each new variable is inserted after that (17, 18, 19, ... ).

The Assembler mostly just manipulates bytes and strings. The input is a file written in ASCII, filled with comments and spaces.  the output is a file with only 1s and 0s.  Each instruction is converted into a line of machine code, with a 1-to-1 correspondence.  This makes [Disassembly](https://en.wikipedia.org/wiki/Disassembler) possible, because you can then reverse the process, and retrieve most of the assembly source again.






How to Get
-------------

~~~
go get github.com/fractalbach/nandGo2tetris/hackasm
~~~


Usage 
----------

~~~
hackasm [-o OUTPUT][-v][-t] INPUT
~~~

Normal Usage
----------------

If output filename is specified, the default output location
is the input filename with a .hack file extension.  From the command line, the easiest usage looks like:

~~~
hackasm INPUT
~~~

The program will attempt to resolve absolute paths, and place the output file in the same folder as the input file.  The name will replace .asm with .hack, and write machine code to the file.

Special Usage
---------------

There are some other flags which are mostly used for debugging:

### Verbose (-v)

Prints some formatted machine code to standard output instead of creating a file.  This is formatted in a friendly way, with line numbers, the instruction type, and the equivalent assembly instruction alongside the machine code.


~~~
hackasm -v INPUT
~~~


### Table (-t)

Prints only the symbol table that is created AFTER the assembler does its first pass.  All custom symbols will be included in this table.  The table is internally stored in a Hash Map: (string) -> (int).  Only the symbol table, and no machine code, is printed to standard output.

~~~
hackasm -t INPUT
~~~
