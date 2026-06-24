// Esempio: I/O virtuale — tastiera e display "dal vivo" via callback
// Il firmware legge 3 tasti con RDR e li rimanda al display con WMP; i callback
// in Go forniscono i tasti e stampano le cifre nel momento in cui escono.
package main

import (
	"fmt"
	"github.com/retronet-labs/retronet-4004/cpu"
)

func main() {
	rom := cpu.NewROM(make([]byte, 256))
	ram := cpu.NewRAM()

	rom.Data[0x000] = cpu.LDM(0)
	rom.Data[0x001] = cpu.DCL() // banco 0
	rom.Data[0x002] = cpu.RDR() // A = tasto
	rom.Data[0x003] = cpu.WMP() // display = A
	rom.Data[0x004] = cpu.RDR()
	rom.Data[0x005] = cpu.WMP()
	rom.Data[0x006] = cpu.RDR()
	rom.Data[0x007] = cpu.WMP()
	const haltAddr = 0x008
	rom.Data[0x008] = cpu.JUN(0x0) // JUN 0x008 → halt
	rom.Data[0x009] = 0x08

	c := cpu.NewCPU4004()

	// Tastiera virtuale: fornisce i tasti 7, 5, 3 in sequenza.
	tasti := []uint8{7, 5, 3}
	i := 0
	c.KeyboardFunc = func() uint8 {
		k := tasti[i]
		i++
		return k
	}
	// Display virtuale: stampa ogni cifra appena il firmware la invia con WMP.
	c.DisplayFunc = func(n uint8) {
		fmt.Printf("display ← %d\n", n)
	}

	fmt.Println("=== I/O virtuale: il firmware legge 3 tasti e li mostra ===")
	for {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			return
		}
		if c.PC == haltAddr {
			break
		}
	}
}
