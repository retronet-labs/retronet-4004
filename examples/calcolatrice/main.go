// Esempio: calcolatrice (stadio A) — somma di due cifre lette dalla tastiera.
//
// Il firmware in ROM realizza il ciclo della calcolatrice in miniatura:
//
//	leggi tasto → leggi tasto → somma (BCD) → mostra il risultato
//
// La tastiera e il display "reali" sono i callback Go (KeyboardFunc/DisplayFunc),
package main

import (
	"fmt"
	"github.com/retronet-labs/retronet-4004/cpu"
)

func main() {
	rom := cpu.NewROM(make([]byte, 256))
	ram := cpu.NewRAM()

	// --- FIRMWARE: A + B su una cifra, risultato a 2 cifre (decine, unità) ---
	rom.Data[0x000] = cpu.LDM(0)
	rom.Data[0x001] = cpu.DCL()       // banco RAM 0
	rom.Data[0x002] = cpu.RDR()       // A = primo tasto
	rom.Data[0x003] = cpu.XCH(cpu.R0) // R0 = primo operando
	rom.Data[0x004] = cpu.RDR()       // A = secondo tasto
	rom.Data[0x005] = cpu.XCH(cpu.R1) // R1 = secondo operando
	rom.Data[0x006] = cpu.LD(cpu.R0)  // A = R0
	rom.Data[0x007] = cpu.CLC()       // nessun riporto in ingresso
	rom.Data[0x008] = cpu.ADD(cpu.R1) // A = R0 + R1
	rom.Data[0x009] = cpu.DAA()       // correzione BCD: A = unità, C = riporto
	rom.Data[0x00A] = cpu.XCH(cpu.R2) // R2 = cifra unità (salvata)
	rom.Data[0x00B] = cpu.TCC()       // A = riporto → cifra decine
	rom.Data[0x00C] = cpu.WMP()       // display ← decine
	rom.Data[0x00D] = cpu.LD(cpu.R2)  // A = unità
	rom.Data[0x00E] = cpu.WMP()       // display ← unità
	const haltAddr = 0x00F
	rom.Data[0x00F] = cpu.JUN(0x0) // halt (salto su se stesso)
	rom.Data[0x010] = 0x0F

	c := cpu.NewCPU4004()

	// --- Tastiera virtuale: la sequenza di tasti "premuti" ---
	tasti := []uint8{7, 5}
	i := 0
	c.KeyboardFunc = func() uint8 {
		k := tasti[i]
		i++
		return k
	}

	// --- Display virtuale: mostra ogni cifra inviata dal firmware ---
	c.DisplayFunc = func(n uint8) {
		fmt.Printf("display ← %d\n", n)
	}

	fmt.Println("=== Calcolatrice (stadio A): somma di due cifre ===")
	fmt.Printf("tasti premuti: %v\n\n", tasti)
	for {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			return
		}
		if c.PC == haltAddr {
			break
		}
	}
	fmt.Println("\n→ il display ha ricevuto le cifre del risultato (decine, unità)")
}
