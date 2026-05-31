package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	c := cpu.NewCPU4004()

	// salva il carry di un'addizione in R0
	program := []byte{
		cpu.LDM(0xF),    // A = 15
		cpu.XCH(cpu.R0), // R0 = 15
		cpu.LDM(0xF),    // A = 15
		cpu.ADD(cpu.R0), // A = 14, C = true (overflow)
		cpu.TCC(),       // A = 1 (carry salvato), C = false
		cpu.XCH(cpu.R1), // R1 = 1 (overflow memorizzato)
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
