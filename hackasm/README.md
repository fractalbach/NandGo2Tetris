Hack Assembly in Go
===========================

- Converts assembly code into machine code instructions.
- Resolves custom symbols to memory addresses.
- Disassemble machine code back into assembly.




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
