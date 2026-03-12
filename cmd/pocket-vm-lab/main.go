package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"os"

	"github.com/ponyfly6/pocket-vm-lab/internal/vm"
)

func main() {
	fileFlag := flag.String("file", "", "path to bytecode file to execute")
	demoFlag := flag.Bool("demo", false, "run a hard-coded demo program")
	flag.Parse()

	if *fileFlag == "" && !*demoFlag {
		fmt.Fprintln(os.Stderr, "Usage: pocket-vm-lab -file <program.bin> or -demo")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var program []byte

	if *demoFlag {
		program = createDemoProgram()
	} else {
		var err error
		program, err = os.ReadFile(*fileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file: %v\n", err)
			os.Exit(1)
		}
	}

	machine := vm.New(program, os.Stdout)
	if err := machine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "execution error: %v\n", err)
		os.Exit(1)
	}
}

// createDemoProgram returns a bytecode program that computes 3 + 5 = 8.
func createDemoProgram() []byte {
	// Program:
	//   CONST 3
	//   CONST 5
	//   ADD
	//   PRINT
	//   HALT
	program := make([]byte, 0, 21)

	// CONST 3
	program = append(program, byte(vm.OpConst))
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], 3)
	program = append(program, buf[:]...)

	// CONST 5
	program = append(program, byte(vm.OpConst))
	binary.LittleEndian.PutUint64(buf[:], 5)
	program = append(program, buf[:]...)

	// ADD
	program = append(program, byte(vm.OpAdd))

	// PRINT
	program = append(program, byte(vm.OpPrint))

	// HALT
	program = append(program, byte(vm.OpHalt))

	return program
}
