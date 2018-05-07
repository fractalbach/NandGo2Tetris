

// This is a random example file.
// Comments don't go into machine code.

@10   // remove this comment
M=1  //Here's another one / / /// / strange one.

@15  //\/\/\/\//\/\\///\/\\/\\/\
M=0 
 
@12
M=A

(HELLO)
@13
D=A

@14
M=D

// The first time you get here, RAM[15] should be 0.
// The second time, it should be 1.
// jumps to the end if RAM[15] > 15

@15
D=M

@END
M; JGT

(THERE)
@15
D=M

@HELLO
D; JEQ

(END)
@END
0; JMP



//empty space