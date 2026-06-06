package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	// Demo: WR0 + RD0 — salva e rileggi un flag di stato dalla RAM.
	//
	//   LDM 0 / DCL       → banco 0
	//   FIM R0, 0x00 / SRC R0
	//   LDM 0xC           → A = 12 (flag da salvare)
	//   WR0               → status[0][0][0] = 12
	//   LDM 0             → azzera A
	//   RD0               → A = status[0][0][0] = 12

	rom := cpu.NewROM(make([]byte, 4096))
	ram := cpu.NewRAM()

	rom.Data[0x000] = cpu.LDM(0)
	rom.Data[0x001] = cpu.DCL()
	rom.Data[0x002] = cpu.FIM(cpu.R0)
	rom.Data[0x003] = 0x00
	rom.Data[0x004] = cpu.SRC(cpu.R0)
	rom.Data[0x005] = cpu.LDM(0xC)
	rom.Data[0x006] = cpu.WR0()
	rom.Data[0x007] = cpu.LDM(0)
	rom.Data[0x008] = cpu.RD0()

	c := cpu.NewCPU4004()
	fmt.Println("=== Demo WR0 + RD0 ===")

	for i := 0; i < 9; i++ {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			break
		}
	}

	fmt.Printf("A = %X (atteso C: round-trip WR0→RD0)\n", c.A)
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
