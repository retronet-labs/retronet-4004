package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	// Demo: WRM — scrive valori BCD in RAM.
	//
	// Programma:
	//   LDM 0     → A=0
	//   DCL       → CL=0 (banco RAM 0)
	//   FIM R0, 0x05 → R0=0 (registro 0), R1=5 (carattere 5)
	//   SRC R0    → SRCAddr = 0x05
	//   LDM 7     → A=7
	//   WRM       → ram.Data[0][0][5] = 7
	//   FIM R0, 0x06 → R0=0, R1=6 (carattere 6)
	//   SRC R0    → SRCAddr = 0x06
	//   LDM 3     → A=3
	//   WRM       → ram.Data[0][0][6] = 3

	rom := cpu.NewROM(make([]byte, 4096))
	ram := cpu.NewRAM()

	rom.Data[0x000] = cpu.LDM(0)
	rom.Data[0x001] = cpu.DCL()
	rom.Data[0x002] = cpu.FIM(cpu.R0)
	rom.Data[0x003] = 0x05 // R0=0, R1=5
	rom.Data[0x004] = cpu.SRC(cpu.R0)
	rom.Data[0x005] = cpu.LDM(7)
	rom.Data[0x006] = cpu.WRM()
	rom.Data[0x007] = cpu.FIM(cpu.R0)
	rom.Data[0x008] = 0x06 // R0=0, R1=6
	rom.Data[0x009] = cpu.SRC(cpu.R0)
	rom.Data[0x00A] = cpu.LDM(3)
	rom.Data[0x00B] = cpu.WRM()

	c := cpu.NewCPU4004()
	fmt.Println("=== Demo WRM ===")

	for i := 0; i < 12; i++ {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			break
		}
	}

	fmt.Printf("ram.Data[0][0][5] = %d (atteso 7)\n", ram.Data[0][0][5])
	fmt.Printf("ram.Data[0][0][6] = %d (atteso 3)\n", ram.Data[0][0][6])
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
