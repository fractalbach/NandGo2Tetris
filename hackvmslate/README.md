Hack Virtual Machine Translator
=====================================


Translates virtual machine language into assembly language.

### About the VM language

The nand2tetris vm language is a [Stack Machine](https://en.wikipedia.org/wiki/Stack_machine),
which includes simple stack commands `pop` and `push`, and commands to use functions and labels.


## Implementing Functions

The most challenging part of making this VM translator was implementing function **calls** and **returns**.

When a function is **called**, the arguments, local variables, and pointers 
(stack pointer, local variable pointer, "this" and "that") have to pushed to the stack prior to the
execution of the function.  

After the function finishes its routine, it then **returns**.
All of those pointers need to be restored, and all that must remain on the stack is a single value:
the *return value*.



## Challenges and Solutions

I had a bug in the my implementation of the function **return**.
It took awhile even to identify the problem, because it would only happen in certain cases.
It's challenging to debug a problem with assembly because problems begin to propogate throughout the program.

To **identify where the problem was happening**, I had to look at **when things started to go wrong**.
Using simpler test programs helped to achieve this.
I noticed that the stack pointer would get thrown off in certain situations, but only after functions returned.

### Solution Commit 

https://github.com/fractalbach/NandGo2Tetris/commit/a32efda9829aa5bbd9f1c16a1fbe7640009cedaa
![capture](https://user-images.githubusercontent.com/32124562/41142281-f7615058-6aa9-11e8-8aa7-2f3539dc9456.PNG)
