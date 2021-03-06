package vmWriter

import (
	"fmt"
	"io"
)

type Segment int

const (
	CONST Segment = 1 << iota
	ARG
	LOCAL
	STATIC
	THIS
	THAT
	POINTER
	TEMP
)

type Command int

const (
	ADD Command = 1 << iota
	SUB
	NEG
	EQ
	GT
	LT
	AND
	OR
	NOT
	MULT
	DIV
)

var SegmentString = map[Segment]string{
	CONST:   "constant",
	ARG:     "argument",
	LOCAL:   "local",
	STATIC:  "static",
	THIS:    "this",
	THAT:    "that",
	POINTER: "pointer",
	TEMP:    "temp",
}

var CommandString = map[Command]string{
	ADD:  "add",
	SUB:  "sub",
	NEG:  "neg",
	EQ:   "eq",
	GT:   "gt",
	LT:   "lt",
	AND:  "and",
	OR:   "or",
	NOT:  "not",
	MULT: `call Math.multiply 2`,
	DIV:  `call Math.divide 2`,
}

var mapSymbolToCmd = map[string]Command{
	"+": ADD,
	"-": SUB,
	"~": NEG,
	"=": EQ,
	">": GT,
	"<": LT,
	"&": AND,
	"|": OR,
	"!": NOT,
	"*": MULT,
	"/": DIV,
}

type VMWriter interface {
	WritePush(segment Segment, n int)
	WritePop(segment Segment, n int)
	WriteArithmetic(cmd Command)
	WriteLabel(label string)
	WriteGoto(label string)
	WriteIf(label string)
	WriteCall(name string, nArgs int)
	WriteFunction(name string, nLocals int)
	WriteReturn()
}

func NewVMWriter(w io.Writer) VMWriter {
	return &vmWriter{
		w: w,
	}
}

type vmWriter struct {
	w io.Writer
}

func (vw *vmWriter) WritePush(seg Segment, n int) {
	fmt.Fprintln(vw.w, "push", SegmentString[seg], n)
}
func (vw *vmWriter) WritePop(seg Segment, n int) {
	fmt.Fprintln(vw.w, "pop", SegmentString[seg], n)
}
func (vw *vmWriter) WriteArithmetic(cmd Command) {
	fmt.Fprintln(vw.w, CommandString[cmd])
}
func (vw *vmWriter) WriteLabel(label string) {
	fmt.Fprintln(vw.w, "label", label)
}
func (vw *vmWriter) WriteGoto(label string) {
	fmt.Fprintln(vw.w, "goto", label)
}
func (vw *vmWriter) WriteIf(label string) {
	fmt.Fprintln(vw.w, "if-goto", label)
}
func (vw *vmWriter) WriteCall(name string, nArgs int) {
	fmt.Fprintln(vw.w, "call", name, nArgs)
}
func (vw *vmWriter) WriteFunction(name string, nLocals int) {
	fmt.Fprintln(vw.w, "function", name, nLocals)
}
func (vw *vmWriter) WriteReturn() {
	fmt.Fprintln(vw.w, "return")
}

func OpToCmd(op string) Command {
	cmd, ok := mapSymbolToCmd[op]
	if !ok {
		s := fmt.Sprintln("The given operator ", op, " cannot be converted to a command.")
		panic(s)
	}
	return cmd
}
