package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	c := cpu.NewCPU4004()

	// DCL: seleziona il banco RAM attivo tramite il registro CL
	// Seleziona banco RAM 2, poi torna al banco 0
	program := []byte{
		cpu.LDM(2), // A = 2
		cpu.DCL(),  // CL = 2 — banco RAM 2 attivo
		cpu.LDM(0), // A = 0
		cpu.DCL(),  // CL = 0 — banco RAM 0 attivo
	}

	fmt.Println("=== BEFORE ===")
	printCPU(c)

	for i, op := range program {
		fmt.Printf("\nSTEP %d\n", i)
		fmt.Printf("Executing opcode: 0x%02X\n", op)

		if err := c.Execute(op); err != nil {
			panic(err)
		}

		printCPU(c)
	}

	fmt.Println("\n=== FINAL STATE ===")
	printCPU(c)
}

func printCPU(c *cpu.CPU4004) {
	fmt.Printf("A=%d C=%v\n", c.A, c.C)

	for i := 0; i < 16; i++ {
		fmt.Printf("R%X=%d ", i, c.R[i])

		if (i+1)%4 == 0 {
			fmt.Println()
		}
	}
}
