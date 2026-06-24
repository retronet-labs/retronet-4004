// Programma: calcolatrice BCD a cifra singola — calcola 7 + 5 con correzione DAA.
package main

import (
	"fmt"
	"github.com/retronet-labs/retronet-4004/cpu"
)

func main() {
	rom := cpu.NewROM(make([]byte, 256))
	ram := cpu.NewRAM()

	// --- SETUP ---
	rom.Data[0x000] = cpu.LDM(0) // A = 0 (serve perché DCL copia A in CL)
	rom.Data[0x001] = cpu.DCL()  // CL = 0 → banco RAM 0 attivo

	// --- scrivi operando A=7 nella cella [0][0][0] ---
	rom.Data[0x002] = cpu.FIM(cpu.R0) // R0=0 (registro RAM), R1=0 (carattere) → indirizzo 0x00
	rom.Data[0x003] = 0x00
	rom.Data[0x004] = cpu.SRC(cpu.R0) // invia 0x00 al chip RAM: "lavora sulla cella [0][0][0]"
	rom.Data[0x005] = cpu.LDM(7)      // A = 7
	rom.Data[0x006] = cpu.WRM()       // RAM[0][0][0] = 7

	// --- scrivi operando B=5 nella cella [0][0][1] ---
	rom.Data[0x007] = cpu.FIM(cpu.R0) // R0=0 (registro), R1=1 (carattere) → indirizzo 0x01
	rom.Data[0x008] = 0x01
	rom.Data[0x009] = cpu.SRC(cpu.R0) // seleziona la cella [0][0][1]
	rom.Data[0x00A] = cpu.LDM(5)      // A = 5
	rom.Data[0x00B] = cpu.WRM()       // RAM[0][0][1] = 5

	// --- calcolo: A = cella0 + cella1, corretto in BCD ---
	rom.Data[0x00C] = cpu.FIM(cpu.R0) // torna a puntare al carattere 0
	rom.Data[0x00D] = 0x00
	rom.Data[0x00E] = cpu.SRC(cpu.R0) // seleziona [0][0][0]
	rom.Data[0x00F] = cpu.RDM()       // A = RAM[0][0][0] = 7
	rom.Data[0x010] = cpu.CLC()       // C = 0 → somma "pulita", senza riporto residuo
	rom.Data[0x011] = cpu.FIM(cpu.R0) // punta al carattere 1
	rom.Data[0x012] = 0x01
	rom.Data[0x013] = cpu.SRC(cpu.R0) // seleziona [0][0][1]
	rom.Data[0x014] = cpu.ADM()       // A = A + RAM[0][0][1] + C = 7 + 5 + 0 = 12 (0xC), C=0
	rom.Data[0x015] = cpu.DAA()       // A = 2, C = 1  ← correzione BCD: 12 → cifra 2 + riporto 1

	// --- salva la cifra unità nella cella [0][0][2] e sulla porta di output ---
	rom.Data[0x016] = cpu.FIM(cpu.R0) // punta al carattere 2
	rom.Data[0x017] = 0x02
	rom.Data[0x018] = cpu.SRC(cpu.R0) // seleziona [0][0][2]
	rom.Data[0x019] = cpu.WRM()       // RAM[0][0][2] = 2  (cifra unità del risultato)
	rom.Data[0x01A] = cpu.WMP()       // Port[0] = 2       (uscita verso il "display")

	// --- HALT: loop infinito su se stesso ---
	rom.Data[0x01B] = cpu.JUN(0x0) // JUN 0x01B → salta a se stesso
	rom.Data[0x01C] = 0x1B

	const haltAddr = 0x01B

	c := cpu.NewCPU4004()
	c.Trace = true

	fmt.Println("=== Calcolatrice BCD: 7 + 5 (cifra singola) ===")
	fmt.Println()

	for {
		if err := c.Step(rom, ram); err != nil {
			fmt.Printf("Errore: %v\n", err)
			return
		}
		if c.PC == haltAddr {
			break
		}
	}

	fmt.Println()
	fmt.Printf("RAM[0][0][0] = %d  (operando A, atteso: 7)\n", ram.Data[0][0][0])
	fmt.Printf("RAM[0][0][1] = %d  (operando B, atteso: 5)\n", ram.Data[0][0][1])
	fmt.Printf("RAM[0][0][2] = %d  (cifra unità, atteso: 2)\n", ram.Data[0][0][2])
	fmt.Printf("riporto (C)  = %v  (atteso: true — la decina)\n", c.C)
	fmt.Printf("Port[0]      = %d  (uscita display, atteso: 2)\n", ram.Port[0])

	// 7 + 5 = 12  →  cifra unità 2, riporto (decina) 1
	if ram.Data[0][0][2] == 2 && c.C && ram.Port[0] == 2 {
		fmt.Println("✓ Corretto! 7 + 5 = 12 (cifra 2, riporto 1)")
	} else {
		fmt.Println("✗ Errore!")
	}
}
