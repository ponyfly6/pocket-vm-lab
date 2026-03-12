package vm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Common errors returned by the VM.
var (
	ErrInvalidOpcode   = errors.New("invalid opcode")
	ErrStackUnderflow  = errors.New("stack underflow")
	ErrProgramTooShort = errors.New("program too short")
)

const (
	maxStackDepth = 1024
)

// VM represents the state of the virtual machine.
type VM struct {
	program []byte
	pc      int
	stack   []int64
	halted  bool
	output  io.Writer
}

// New creates a new VM instance with the given bytecode program.
func New(program []byte, output io.Writer) *VM {
	if output == nil {
		output = io.Discard
	}
	return &VM{
		program: program,
		pc:      0,
		stack:   make([]int64, 0, 16),
		halted:  false,
		output:  output,
	}
}

// Run executes the program until HALT or an error.
func (vm *VM) Run() error {
	for !vm.halted {
		if vm.pc >= len(vm.program) {
			return nil
		}
		op := Opcode(vm.program[vm.pc])
		vm.pc++

		if err := vm.execute(op); err != nil {
			return err
		}
	}
	return nil
}

// execute dispatches to the appropriate instruction handler.
func (vm *VM) execute(op Opcode) error {
	switch op {
	case OpConst:
		return vm.execConst()
	case OpAdd:
		return vm.execAdd()
	case OpPrint:
		return vm.execPrint()
	case OpHalt:
		vm.halted = true
		return nil
	default:
		return fmt.Errorf("%w: 0x%02X at pc=%d", ErrInvalidOpcode, op, vm.pc-1)
	}
}

// execConst reads an int64 and pushes it onto the stack.
func (vm *VM) execConst() error {
	if vm.pc+8 > len(vm.program) {
		return fmt.Errorf("%w: need 8 bytes for CONST operand at pc=%d", ErrProgramTooShort, vm.pc)
	}

	val := int64(binary.LittleEndian.Uint64(vm.program[vm.pc : vm.pc+8]))
	vm.pc += 8

	if len(vm.stack) >= maxStackDepth {
		return errors.New("stack overflow")
	}
	vm.stack = append(vm.stack, val)
	return nil
}

// execAdd pops two values, pushes their sum.
func (vm *VM) execAdd() error {
	if len(vm.stack) < 2 {
		return fmt.Errorf("%w: need 2 values for ADD", ErrStackUnderflow)
	}

	b := vm.stack[len(vm.stack)-1]
	a := vm.stack[len(vm.stack)-2]
	vm.stack = vm.stack[:len(vm.stack)-2]
	vm.stack = append(vm.stack, a+b)
	return nil
}

// execPrint pops and prints the top of stack.
func (vm *VM) execPrint() error {
	if len(vm.stack) < 1 {
		return fmt.Errorf("%w: need 1 value for PRINT", ErrStackUnderflow)
	}

	val := vm.stack[len(vm.stack)-1]
	vm.stack = vm.stack[:len(vm.stack)-1]
	fmt.Fprintln(vm.output, val)
	return nil
}

// Halted returns true if the VM has halted.
func (vm *VM) Halted() bool {
	return vm.halted
}

// Stack returns a copy of the current stack.
func (vm *VM) Stack() []int64 {
	cpy := make([]int64, len(vm.stack))
	copy(cpy, vm.stack)
	return cpy
}
