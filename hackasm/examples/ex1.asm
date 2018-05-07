

// Easy Example
  // Has no symbols

@R0  // go to register 0.
M=1  // set RAM[0] to 1.

@R2  // go to register 2.
MD=A  // set RAM[2] and the D register equal to 2.

// Goes to the register 3,
// and sets RAM[3] equal to the value saved in D register
// D register should be equal to 2
@R3
M=D

// Result should be 6 lines of machine code.
// Like so: (without the spaces)
// 
// 000 0 000000 000 000
// 111 0 111111 001 000
//
// 000 0 000000 000 010
// 111 0 110000 011 000
// 
// 000 0 000000 000 011
// 111 0 001100 001 000
// 