package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	c := cpu.NewCPU4004()

	program := []byte{
		cpu.LDM(2),      // A = 2
		cpu.XCH(cpu.R0), // R0 = 2, A = 0
		cpu.LDM(3),      // A = 3
		cpu.ADD(cpu.R0), // A = A + R0
	}

	for _, op := range program {
		if err := c.Execute(op); err != nil {
			panic(err)
		}
	}

	fmt.Println("A =", c.A)
	fmt.Println("Carry =", c.C)
	fmt.Println("R0 =", c.R[cpu.R0])
}
