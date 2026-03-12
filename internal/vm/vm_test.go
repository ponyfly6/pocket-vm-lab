package vm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"testing"
)

func encodeConst(val int64) []byte {
	buf := make([]byte, 9)
	buf[0] = byte(OpConst)
	binary.LittleEndian.PutUint64(buf[1:9], uint64(val))
	return buf
}

func TestConstInstruction(t *testing.T) {
	program := encodeConst(42)

	vm := New(program, nil)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stack := vm.Stack()
	if len(stack) != 1 {
		t.Fatalf("expected stack size 1, got %d", len(stack))
	}
	if stack[0] != 42 {
		t.Fatalf("expected 42, got %d", stack[0])
	}
}

func TestAddInstruction(t *testing.T) {
	// CONST 3; CONST 5; ADD; HALT
	program := append(encodeConst(3), encodeConst(5)...)
	program = append(program, byte(OpAdd), byte(OpHalt))

	vm := New(program, nil)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !vm.Halted() {
		t.Error("expected VM to be halted")
	}

	stack := vm.Stack()
	if len(stack) != 1 {
		t.Fatalf("expected stack size 1, got %d", len(stack))
	}
	if stack[0] != 8 {
		t.Fatalf("expected 8, got %d", stack[0])
	}
}

func TestAddStackUnderflow(t *testing.T) {
	// ADD without pushing anything
	program := []byte{byte(OpAdd)}

	vm := New(program, nil)
	err := vm.Run()
	if err == nil {
		t.Fatal("expected error for stack underflow, got nil")
	}
	if !errors.Is(err, ErrStackUnderflow) {
		t.Fatalf("expected ErrStackUnderflow, got %v", err)
	}
}

func TestPrintInstruction(t *testing.T) {
	var buf bytes.Buffer
	// CONST 100; PRINT; HALT
	program := append(encodeConst(100), byte(OpPrint), byte(OpHalt))

	vm := New(program, &buf)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "100\n"
	if buf.String() != expected {
		t.Fatalf("expected output %q, got %q", expected, buf.String())
	}
}

func TestPrintStackUnderflow(t *testing.T) {
	// PRINT without pushing anything
	program := []byte{byte(OpPrint)}

	vm := New(program, nil)
	err := vm.Run()
	if err == nil {
		t.Fatal("expected error for stack underflow, got nil")
	}
	if !errors.Is(err, ErrStackUnderflow) {
		t.Fatalf("expected ErrStackUnderflow, got %v", err)
	}
}

func TestHaltInstruction(t *testing.T) {
	// HALT
	program := []byte{byte(OpHalt)}

	vm := New(program, nil)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !vm.Halted() {
		t.Error("expected VM to be halted")
	}
}

func TestInvalidOpcode(t *testing.T) {
	// invalid opcode 0x99
	program := []byte{0x99}

	vm := New(program, nil)
	err := vm.Run()
	if err == nil {
		t.Fatal("expected error for invalid opcode, got nil")
	}
	if !errors.Is(err, ErrInvalidOpcode) {
		t.Fatalf("expected ErrInvalidOpcode, got %v", err)
	}
}

func TestMultipleConstAndAdd(t *testing.T) {
	// Compute (10 + 20) + (30 + 40) = 100
	program := append(encodeConst(10), encodeConst(20)...)
	program = append(program, byte(OpAdd))
	program = append(program, encodeConst(30)...)
	program = append(program, encodeConst(40)...)
	program = append(program, byte(OpAdd))
	program = append(program, byte(OpAdd), byte(OpHalt))

	vm := New(program, nil)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stack := vm.Stack()
	if len(stack) != 1 {
		t.Fatalf("expected stack size 1, got %d", len(stack))
	}
	if stack[0] != 100 {
		t.Fatalf("expected 100, got %d", stack[0])
	}
}

func TestConstTruncatedOperand(t *testing.T) {
	// CONST with only 3 bytes instead of 8
	program := []byte{byte(OpConst), 0x01, 0x02, 0x03}

	vm := New(program, nil)
	err := vm.Run()
	if err == nil {
		t.Fatal("expected error for truncated operand, got nil")
	}
	if !errors.Is(err, ErrProgramTooShort) {
		t.Fatalf("expected ErrProgramTooShort, got %v", err)
	}
}

func TestNegativeValues(t *testing.T) {
	// CONST -5; CONST 3; ADD; HALT
	program := append(encodeConst(-5), encodeConst(3)...)
	program = append(program, byte(OpAdd), byte(OpHalt))

	vm := New(program, nil)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stack := vm.Stack()
	if len(stack) != 1 {
		t.Fatalf("expected stack size 1, got %d", len(stack))
	}
	if stack[0] != -2 {
		t.Fatalf("expected -2, got %d", stack[0])
	}
}

func TestEndOfProgram(t *testing.T) {
	// No HALT - program just ends
	program := encodeConst(7)

	vm := New(program, nil)
	err := vm.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should not be halted since we didn't hit HALT
	if vm.Halted() {
		t.Error("expected VM not to be halted without explicit HALT")
	}

	stack := vm.Stack()
	if len(stack) != 1 {
		t.Fatalf("expected stack size 1, got %d", len(stack))
	}
	if stack[0] != 7 {
		t.Fatalf("expected 7, got %d", stack[0])
	}
}
