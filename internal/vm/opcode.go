// Package vm implements a minimal stack-based virtual machine.
package vm

// Opcode represents a single instruction in the VM.
type Opcode byte

const (
	// OpConst pushes a constant onto the stack.
	// Format: CONST <int64 little-endian>
	OpConst Opcode = 0x01

	// OpAdd pops two values from the stack and pushes their sum.
	OpAdd Opcode = 0x02

	// OpPrint pops the top value from the stack and prints it.
	OpPrint Opcode = 0x03

	// OpHalt stops execution.
	OpHalt Opcode = 0xFF
)

// String returns a human-readable name for the opcode.
func (op Opcode) String() string {
	switch op {
	case OpConst:
		return "CONST"
	case OpAdd:
		return "ADD"
	case OpPrint:
		return "PRINT"
	case OpHalt:
		return "HALT"
	default:
		return "UNKNOWN"
	}
}
