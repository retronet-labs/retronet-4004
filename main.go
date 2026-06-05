package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	// Demo: WRM + RDM — write/read round-trip sulla RAM virtuale.
	//
	//   LDM 0     → A=0, poi DCL → CL=0 (banco 0)
	//   FIM R0, 0x05 → R0=0, R1=5 (registro 0, carattere 5)
	//   SRC R0    → SRCAddr=0x05
	//   LDM 7     → A=7
	//   WRM       → ram.Data[0][0][5] = 7
	//   LDM 0     → azzera A
	//   RDM       → A = ram.Data[0][0][5] = 7

	rom := cpu.NewROM(make([]byte, 4096))
	ram := cpu.NewRAM()

	rom.Data[0x000] = cpu.LDM(0)
	rom.Data[0x001] = cpu.DCL()
	rom.Data[0x002] = cpu.FIM(cpu.R0)
	rom.Data[0x003] = 0x05
	rom.Data[0x004] = cpu.SRC(cpu.R0)
	rom.Data[0x005] = cpu.LDM(7)
	rom.Data[0x006] = cpu.WRM()
	rom.Data[0x007] = cpu.LDM(0)
	rom.Data[0x008] = cpu.RDM()

	c := cpu.NewCPU4004()
	fmt.Println("=== Demo WRM + RDM ===")

	for i := 0; i < 9; i++ {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			break
		}
	}

	fmt.Printf("A = %d (atteso 7 — letto da RAM dopo WRM)\n", c.A)
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
