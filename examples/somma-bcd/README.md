# Esempio: Calcolatrice BCD a cifra singola (7 + 5)

Primo **firmware reale** del progetto: il primo mattone della calcolatrice.
Somma due cifre BCD tenute in RAM e produce il risultato corretto in formato
decimale — esattamente come farebbe la Busicom 141-PF.

---

## Perché 7 + 5 e non 3 + 4

Il punto di tutto l'esempio è l'istruzione **DAA**. Con `3 + 4 = 7` il risultato
è già una cifra BCD valida e DAA non farebbe nulla.

`7 + 5 = 12` invece **rompe** la cifra singola: 12 non è una cifra decimale (0–9).
Questo costringe la CPU a fare ciò che fa una calcolatrice vera — produrre
**cifra unità = 2** e **riporto = 1** verso le decine. È il caso che giustifica
l'esistenza di DAA, e fa da ponte verso lo Step 11 (numeri multi-cifra), dove
quel riporto verrà propagato.

---

## Algoritmo

```
setup:
  CL = 0                       (banco RAM 0)
  RAM[0][0][0] = 7             (operando A)
  RAM[0][0][1] = 5             (operando B)

calcolo:
  A = RAM[0][0][0]             (RDM → A = 7)
  C = 0                        (CLC: somma pulita)
  A = A + RAM[0][0][1] + C     (ADM → A = 12 = 0xC, non valido BCD)
  DAA                          (A = 2, C = 1: cifra 2 + riporto 1)

output:
  RAM[0][0][2] = A             (cifra unità)
  Port[0]      = A             (uscita display)
  ← il riporto resta nel flag C (la decina)
```

---

## Layout RAM (banco 0, registro 0)

| Cella | Contenuto |
|-------|-----------|
| char 0 | operando A = 7 |
| char 1 | operando B = 5 |
| char 2 | risultato, cifra unità = 2 |
| Port[0] | risultato, cifra unità = 2 (uscita `WMP`) |

Il **riporto** (la decina) non sta in RAM: vive nel flag `C`. Una cifra BCD
singola non può contenerlo — ed è proprio il limite che lo Step 11 supererà.

---

## Layout ROM

```
       ── SETUP ──
0x000  LDM 0          A = 0 (serve per DCL)
0x001  DCL            CL = 0 (banco RAM 0)
       ── scrivi A=7 nella cella 0 ──
0x002  FIM R0, 0x00   R0=0 (registro), R1=0 (carattere)
0x004  SRC R0         seleziona [0][0][0]
0x005  LDM 7          A = 7
0x006  WRM            RAM[0][0][0] = 7
       ── scrivi B=5 nella cella 1 ──
0x007  FIM R0, 0x01   R1=1 → carattere 1
0x009  SRC R0         seleziona [0][0][1]
0x00A  LDM 5          A = 5
0x00B  WRM            RAM[0][0][1] = 5
       ── calcolo ──
0x00C  FIM R0, 0x00   torna al carattere 0
0x00E  SRC R0         seleziona [0][0][0]
0x00F  RDM            A = 7
0x010  CLC            C = 0
0x011  FIM R0, 0x01   carattere 1
0x013  SRC R0         seleziona [0][0][1]
0x014  ADM            A = 7 + 5 = 12 (0xC), C = 0
0x015  DAA            A = 2, C = 1   ← correzione BCD
       ── salva risultato ──
0x016  FIM R0, 0x02   carattere 2
0x018  SRC R0         seleziona [0][0][2]
0x019  WRM            RAM[0][0][2] = 2
0x01A  WMP            Port[0] = 2
       ── HALT ──
0x01B  JUN 0x01B
```

---

## Il momento chiave: ADM → DAA

Nel trace il passaggio da binario a BCD è visibile in due righe:

```
PC=014 OP=EB ADM   A=C  C=false   ← 7+5 = 0xC (12): risultato binario, NON è BCD
PC=015 OP=FB DAA   A=2  C=true    ← DAA corregge: cifra 2, e alza il riporto
```

`ADM` somma binario su 4 bit: `7 + 5 = 12 = 0xC`, un valore che un display
decimale non sa rappresentare. `DAA` vede `A > 9`, aggiunge 6
(`12 + 6 = 18 = 0x12`), tiene il nibble basso `2` e alza `C`. Risultato:
cifra `2` con riporto `1`, cioè "12". Vedi [docs/bcd.md](../../docs/bcd.md)
per il perché del "+6".

---

## ⚠️ ADM somma anche il carry → serve CLC

`ADM` non fa `A + RAM`, fa `A + RAM + C`. Se il carry fosse rimasto alto da
un'operazione precedente, si infilerebbe silenziosamente nella somma. Per
questo prima di `ADM` c'è `CLC`. Qui il carry è già 0, ma azzerarlo
esplicitamente è l'abitudine che diventa indispensabile nello Step 11, dove
il riporto fra cifre si propaga di proposito.

---

## Gli operandi vivono in RAM, non nei registri

Il firmware salva A e B in RAM, poi li **rilegge** (`RDM` per il primo, `ADM`
per il secondo) per fare il conto. È così che lavora una calcolatrice vera:
le cifre stanno in RAM (una per cella), e la CPU le tira dentro l'accumulatore
solo al momento di operarci. Cambiata la cella bersaglio, va sempre rifatto
`SRC` — il chip RAM non si sposta da solo.

---

## Risultato atteso

```
RAM[0][0][0] = 7   (operando A)
RAM[0][0][1] = 5   (operando B)
RAM[0][0][2] = 2   (cifra unità)
riporto (C)  = true (la decina)
Port[0]      = 2   (uscita display)
✓ 7 + 5 = 12 (cifra 2, riporto 1)
```

---

## Come eseguire

```
go run ./examples/somma-bcd
```
