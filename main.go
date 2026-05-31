package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	c := cpu.NewCPU4004()

	// TCS: valore di correzione BCD per la sottrazione
	program := []byte{
		cpu.STC(),       // forza C = true (simula un borrow)
		cpu.TCS(),       // A = 10, C = false
		cpu.XCH(cpu.R0), // salva 10 in R0

		cpu.CLC(),       // C = false (nessun borrow)
		cpu.TCS(),       // A = 9, C = false
		cpu.XCH(cpu.R1), // salva 9 in R1
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
