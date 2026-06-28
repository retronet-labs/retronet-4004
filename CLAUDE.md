# CLAUDE.md — retronet-4004

Emulatore **Intel 4004** in Go (cartella locale `go-4004`, repo GitHub
`retronet-4004`), primo modulo dell'ecosistema RetroNet. CPU completa (46/46
istruzioni), ROM/RAM virtuali, I/O interattivo, CLI e tracing. Panoramica utente:
[readme.md](readme.md).

## Setup su una macchina nuova (handoff)

Clona i repo come cartelle **sibling** sotto la stessa radice (il 4004 va in
`go-4004`). Un clone pulito compila già dalle versioni pubblicate; per delegare
l'aritmetica alla ALU a gate serve il sibling `retronet-hardware`:

```
work/source/
├── retronet-logic/
├── retronet-hardware/   (bridge i4004)
└── go-4004/             (questo repo)
```

`go.work` (non versionato) per il co-sviluppo:

```sh
go work init . ../retronet-hardware ../retronet-logic
```

Lo script globale `retronet/scripts/migrate.sh` fa tutto questo in automatico.

## Comandi

- Test: `go test ./...`
- CLI: `go run ./cmd/retronet-4004 [-trace] [-max N] [-dump-ram] -io <rom>`
- ROM golden: `go run ./testdata/gen` (rigenera `testdata/*.rom`)
- Esempio interattivo: `echo 1.5+2.25= | go run ./cmd/retronet-4004 -io calc.rom`

## Componenti (`cpu/`)

- `cpu.go` (struct CPU4004 + `Step`), `instructions.go` (decoder/executor),
  `opcodes.go`, `helpers.go`, `disasm.go`, `rom.go`, `ram.go` (chip 4002 virtuale).
- **Aritmetica delegata ai gate**: ADD/SUB e le rotazioni RAL/RAR passano dal
  `bridge/i4004` di retronet-hardware (ALU + shifter a porte), non dagli operatori
  Go. Verificato da test di conformità.
- I/O dal vivo: callback `KeyboardFunc`/`DisplayFunc` + flag CLI `-io` (ponte
  stdio). Il `DisplayFunc` in `-io` mappa il nibble 15 → `.` e 11 → `-`.

## Stato

`go test ./...` verde. **46/46 istruzioni**, CLI, trace, Dockerfile multi-stage.
Pubblicato e taggato fino a **`v0.3.0`** (la delega completa ai gate è l'ultimo
lavoro). La calcolatrice da tavolo BCD a virgola fissa vive come firmware `.asm`
in **retronet-asm** (`examples/calcolatrice-*.asm`), eseguita qui con `-io` —
validazione end-to-end assembler ↔ emulatore.

## Note tecniche da ricordare

- **"Stessa pagina"**: `JCN`/`ISZ`/`JIN` usano solo 8 bit d'indirizzo; i 4 bit
  alti vengono dal PC → il salto resta nella pagina da 256 byte corrente.
- **ISZ**: salta quando il registro ≠ 0; per N iterazioni inizializza a `16-N`.
- **Stack** a 3 livelli ciclico (`SP % 3`), come l'hardware.
- **Sottrazione BCD robusta** = complemento a 10 (`M + comp9(S) + 1`); il riporto
  finale dà il confronto `M ≥ S` usato dalla divisione.
