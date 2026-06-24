// Comando gen genera i file .rom di esempio in testdata/, costruendoli con gli
// helper del package cpu (così i byte sono sempre corretti, senza conteggi a mano).
//
// Ogni ROM corrisponde a un programma in examples/, con in più un halt finale
// (JUN su se stesso) dove l'esempio originale si fermava per conteggio di step,
// così la CLI retronet-4004 può eseguirla e fermarsi da sola.
//
// Uso (dalla radice del repo):
//
//	go run ./testdata/gen
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/retronet-labs/retronet-4004/cpu"
)

func main() {
	roms := map[string][]byte{
		// 3 × 4 = 12 tramite addizioni ripetute (loop ISZ). Halt @ 0x00D.
		"moltiplicazione.rom": {
			cpu.LDM(0), cpu.DCL(),
			cpu.FIM(cpu.R0), 0x03, // R1 = 3 (addendo)
			cpu.FIM(cpu.R2), 0x00, // indirizzo RAM 0x00
			cpu.SRC(cpu.R2),
			cpu.LDM(12), cpu.XCH(cpu.R4), // R4 = 12 (contatore 16-4)
			cpu.ADD(cpu.R1),         // LOOP @ 0x09
			cpu.ISZ(cpu.R4), 0x09,   // R4++; se !=0 → 0x09
			cpu.WRM(),               // RAM[0][0][0] = 12
			cpu.JUN(0x0), 0x0D,      // halt @ 0x00D
		},

		// 3 + 5 = 8 chiamando una subroutine (JMS/BBL). Halt @ 0x00B.
		"subroutine.rom": {
			cpu.LDM(0), cpu.DCL(),
			cpu.FIM(cpu.R2), 0x00, cpu.SRC(cpu.R2),
			cpu.FIM(cpu.R0), 0x35,    // R0=3, R1=5
			cpu.JMS(0x0), 0x0D,       // chiama SOMMA @ 0x00D
			cpu.LD(cpu.R5), cpu.WRM(), // recupera risultato, scrivi in RAM
			cpu.JUN(0x0), 0x0B,       // halt @ 0x00B
			cpu.LD(cpu.R0), cpu.ADD(cpu.R1), cpu.XCH(cpu.R5), cpu.BBL(0), // SOMMA @ 0x00D
		},

		// Riempire un array in RAM con 1,2,3,4. Halt @ 0x00F.
		"ram-array.rom": {
			cpu.LDM(0), cpu.DCL(),
			cpu.FIM(cpu.R2), 0x00, // indirizzo iniziale
			cpu.FIM(cpu.R0), 0x01, // primo valore in R1
			cpu.FIM(cpu.R4), 0xC0, // contatore = 12 (16-4)
			cpu.SRC(cpu.R2), cpu.LD(cpu.R1), cpu.WRM(), // LOOP @ 0x08
			cpu.INC(cpu.R1), cpu.INC(cpu.R3),
			cpu.ISZ(cpu.R4), 0x08, // ripeti
			cpu.JUN(0x0), 0x0F,    // halt @ 0x00F
		},

		// Calcolatrice BCD a cifra singola: 7 + 5. Halt @ 0x01B.
		"somma-bcd.rom": {
			cpu.LDM(0), cpu.DCL(),
			cpu.FIM(cpu.R0), 0x00, cpu.SRC(cpu.R0), cpu.LDM(7), cpu.WRM(), // A=7 in cella 0
			cpu.FIM(cpu.R0), 0x01, cpu.SRC(cpu.R0), cpu.LDM(5), cpu.WRM(), // B=5 in cella 1
			cpu.FIM(cpu.R0), 0x00, cpu.SRC(cpu.R0), cpu.RDM(), cpu.CLC(),
			cpu.FIM(cpu.R0), 0x01, cpu.SRC(cpu.R0), cpu.ADM(), cpu.DAA(),
			cpu.FIM(cpu.R0), 0x02, cpu.SRC(cpu.R0), cpu.WRM(), cpu.WMP(), // cifra unità → cella 2 + porta
			cpu.JUN(0x0), 0x1B, // halt @ 0x01B
		},

		// Addizione BCD multi-cifra: 47 + 58 = 105. Halt @ 0x02E.
		"somma-multicifra.rom": {
			cpu.LDM(0), cpu.DCL(),
			cpu.FIM(cpu.R0), 0x00, cpu.SRC(cpu.R0), cpu.LDM(7), cpu.WRM(), // A unità
			cpu.FIM(cpu.R0), 0x01, cpu.SRC(cpu.R0), cpu.LDM(4), cpu.WRM(), // A decine
			cpu.FIM(cpu.R2), 0x10, cpu.SRC(cpu.R2), cpu.LDM(8), cpu.WRM(), // B unità
			cpu.FIM(cpu.R2), 0x11, cpu.SRC(cpu.R2), cpu.LDM(5), cpu.WRM(), // B decine
			cpu.FIM(cpu.R0), 0x00, cpu.FIM(cpu.R2), 0x10, cpu.FIM(cpu.R4), 0x20, cpu.FIM(cpu.R6), 0xE0, cpu.CLC(),
			cpu.SRC(cpu.R0), cpu.RDM(), cpu.SRC(cpu.R2), cpu.ADM(), cpu.DAA(), // LOOP @ 0x1F
			cpu.SRC(cpu.R4), cpu.WRM(), cpu.INC(cpu.R1), cpu.INC(cpu.R3), cpu.INC(cpu.R5),
			cpu.ISZ(cpu.R6), 0x1F,
			cpu.TCC(), cpu.SRC(cpu.R4), cpu.WRM(), // riporto finale → centinaia
			cpu.JUN(0x0), 0x2E,                    // halt @ 0x02E
		},
	}

	outDir := "testdata"
	for name, code := range roms {
		path := filepath.Join(outDir, name)
		if err := os.WriteFile(path, code, 0o644); err != nil {
			fmt.Fprintf(os.Stderr, "errore scrittura %s: %v\n", path, err)
			os.Exit(1)
		}
		fmt.Printf("scritto %s (%d byte)\n", path, len(code))
	}
}
