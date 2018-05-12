Hack Virtual Machine Translator
=====================================


Translates virtual machine language into assembly language.


## Hack VM Language

The commands that operate on the VM stack machine.
For the arithemtic commands, y refers to the top of the stack,
and x refers to the element below the top.

For the memeory access commands, the *segment* and *index*
refer to a specific memory register in the RAM.  This abstraction
allows the VM language to remain seperate from the computer architecture.


Arithemtic Stack Commands
- [x] add (x+y)
- [x] sub (x-y)
- [x] negate (!y)

Logical Stack Commands
- [x] eq (true if x = y, else false)
- [x] gt (true if x > y, else false)
- [x] lt (true if x < y, else false)

Bitwise Stack Commands
- [x] and (x && y)
- [x] or (x || y)
- [x] not (Not y)

Memory Access Stack Commands
- [x] Push segment index
- [x] Pop segment index

Segments
- argument
- local
- static
- constant
- this 
- that
- pointer
- temp


