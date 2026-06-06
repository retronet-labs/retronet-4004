package main

import (
	"fmt"
	"go-4004/cpu"
)

func main() {
	// Demo: WRR + RDR — simula un ciclo di scansione tastiera.
	//
	// Il firmware attiva una riga della tastiera (WRR), poi legge
	// le colonne attive (RDR) e decodifica il tasto premuto (KBP).
	//
	//   LDM 0b0001 → A = 1 (riga 1)
	//   WRR        → rom.Port = 1 (attiva riga 1)
	//   RDR        → A = rom.Port (leggi colonne)
	//   KBP        → A = numero tasto (decodifica one-hot)
	//
	// Simuliamo il tasto nella colonna 3 premuto (bit 2 attivo = 0b0100).

	rom := cpu.NewROM(make([]byte, 4096))
	ram := cpu.NewRAM()

	rom.Data[0x000] = cpu.LDM(0b0001) // attiva riga 1
	rom.Data[0x001] = cpu.WRR()        // invia sulla porta ROM
	rom.Data[0x002] = cpu.RDR()        // leggi risposta colonne
	rom.Data[0x003] = cpu.KBP()        // decodifica one-hot → numero tasto

	// Simuliamo il tasto in colonna 3 premuto
	rom.Port = 0b0100

	c := cpu.NewCPU4004()
	fmt.Println("=== Demo WRR + RDR (scansione tastiera) ===")

	for i := 0; i < 4; i++ {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			break
		}
	}

	fmt.Printf("A = %d (atteso 3: tasto colonna 3 premuto)\n", c.A)
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
