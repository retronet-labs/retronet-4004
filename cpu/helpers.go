package cpu

// Queste funzioni semplificano la creazione di istruzioni restituendo il byte
// opcode corretto per ciascuna istruzione del 4004.

// NOP non esegue alcuna operazione: avanza solo il program counter.
func NOP() byte { return OP_NOP }

// LDM carica il valore immediato v (0–15) nell'accumulatore A. Non tocca il carry.
func LDM(v byte) byte { return OP_LDM | nibble(v) }

// XCH scambia il contenuto dell'accumulatore A con il registro Rr. Non tocca il carry.
func XCH(r byte) byte { return OP_XCH | nibble(r) }

// INC incrementa di 1 il registro Rr (wrap a 4 bit: 0xF→0x0). Non tocca il carry.
func INC(r byte) byte { return OP_INC | nibble(r) }

// ADD somma il registro Rr e il carry all'accumulatore: A = A + Rr + C.
// Imposta C se il risultato supera 0xF (riporto in uscita).
func ADD(r byte) byte { return OP_ADD | nibble(r) }

// LD carica il registro Rr nell'accumulatore A. Non tocca il carry.
func LD(r byte) byte { return OP_LD | nibble(r) }

// SUB sottrae il registro Rr dall'accumulatore con prestito: A = A + ~Rr + C.
// In ingresso C=1 significa "nessun borrow"; in uscita C=1 significa "nessun borrow generato".
func SUB(r byte) byte { return OP_SUB | nibble(r) }

// IAC incrementa l'accumulatore di 1 (Increment Accumulator). Imposta C se A passa da 0xF a 0x0.
func IAC() byte { return OP_IAC }

// DAC decrementa l'accumulatore di 1 (Decrement Accumulator).
// C=1 se non c'è prestito, C=0 se A era 0 (underflow).
func DAC() byte { return OP_DAC }

// CMA complementa l'accumulatore (complemento a 1: inverte tutti i 4 bit). Non tocca il carry.
func CMA() byte { return OP_CMA }

// CLB azzera sia l'accumulatore sia il carry (Clear Both): A=0, C=0.
func CLB() byte { return OP_CLB }

// CLC azzera il carry (Clear Carry): C=0. Non tocca l'accumulatore.
func CLC() byte { return OP_CLC }

// STC imposta il carry (Set Carry): C=1. Non tocca l'accumulatore.
func STC() byte { return OP_STC }

// CMC complementa il carry (Complement Carry): C=!C. Non tocca l'accumulatore.
func CMC() byte { return OP_CMC }

// RAL ruota l'accumulatore a sinistra attraverso il carry: il bit 3 finisce in C, il vecchio C entra nel bit 0.
func RAL() byte { return OP_RAL }

// RAR ruota l'accumulatore a destra attraverso il carry: il bit 0 finisce in C, il vecchio C entra nel bit 3.
func RAR() byte { return OP_RAR }

// TCC trasferisce il carry nell'accumulatore e lo azzera (Transfer Carry and Clear):
// A=1 se C era vero, altrimenti A=0; poi C=0. Utile per trasformare un riporto in cifra.
func TCC() byte { return OP_TCC }

// TCS trasferisce il carry per la sottrazione BCD (Transfer Carry Subtract):
// A=10 se C era vero, altrimenti A=9; poi C=0.
func TCS() byte { return OP_TCS }

// DAA corregge l'accumulatore in formato BCD dopo un'addizione (Decimal Adjust Accumulator):
// se C è vero oppure A>9, aggiunge 6 ad A e imposta C (riporto verso la cifra successiva).
func DAA() byte { return OP_DAA }

// KBP converte un valore one-hot dell'accumulatore nella posizione del bit attivo (Keyboard Process):
// 0000→0, 0001→1, 0010→2, 0100→3, 1000→4; qualsiasi altro valore (più bit a 1) → 0xF (errore).
func KBP() byte { return OP_KBP }

// DCL designa il banco RAM attivo (Designate Command Line): CL = A & 0x7. Non tocca A né il carry.
func DCL() byte { return OP_DCL }

// BBL ritorna da una subroutine (Branch Back and Load): ripristina il PC dallo stack
// e carica il valore immediato v (0–15) nell'accumulatore. Non tocca il carry.
func BBL(v byte) byte { return OP_BBL | nibble(v) }

// Istruzioni di salto — restituiscono il primo byte; il secondo byte (indirizzo/dato) va aggiunto separatamente nella ROM.

// JUN salta incondizionatamente a un indirizzo a 12 bit (Jump Unconditional).
// addrHigh sono i 4 bit alti dell'indirizzo; il secondo byte ne contiene gli 8 bassi.
func JUN(addrHigh byte) byte { return OP_JUN | (addrHigh & 0x0F) }

// JMS chiama una subroutine a un indirizzo a 12 bit (Jump to Subroutine): salva il PC sullo stack, poi salta.
// addrHigh sono i 4 bit alti dell'indirizzo; il secondo byte ne contiene gli 8 bassi.
func JMS(addrHigh byte) byte { return OP_JMS | (addrHigh & 0x0F) }

// JCN salta se la condizione cond è verificata (Jump Conditional), restando nella pagina del PC.
// Nibble cond: bit3=inverti (NOT), bit2=salta se A==0, bit1=salta se C==1, bit0=salta se TEST==0 (non emulato).
// Il secondo byte è l'offset (8 bit) nella pagina corrente.
func JCN(cond byte) byte { return OP_JCN | (cond & 0x0F) }

// ISZ incrementa il registro Rr e salta se il risultato è diverso da zero (Increment and Skip if Zero).
// Il secondo byte è l'offset nella pagina corrente. Tipico contatore di loop: inizializzare Rr a 16-N.
func ISZ(r byte) byte { return OP_ISZ | nibble(r) }

// FIM carica i due nibble del secondo byte nella coppia di registri Rr/Rr+1 (Fetch Immediate).
// rp deve essere pari (R0, R2, R4, ...): Rr riceve il nibble alto, Rr+1 il nibble basso.
func FIM(rp byte) byte { return OP_FIM | (nibble(rp) &^ 1) } // rp pari: coppia Rr/Rr+1

// SRC invia l'indirizzo RAM/ROM sul bus esterno (Send Register Control): indirizzo = (Rr<<4)|Rr+1.
// rp deve essere pari. Va richiamato ogni volta che cambia la cella RAM target prima di WRM/RDM/ecc.
func SRC(rp byte) byte { return OP_SRC | (nibble(rp) &^ 1) }

// FIN legge un byte dalla ROM indirizzato da R0:R1 nella pagina corrente (Fetch Indirect)
// e lo carica nella coppia Rr/Rr+1. rp deve essere pari.
func FIN(rp byte) byte { return OP_FIN | (nibble(rp) &^ 1) }

// JIN salta all'indirizzo contenuto nella coppia Rr/Rr+1, nella pagina corrente (Jump Indirect).
// rp deve essere pari.
func JIN(rp byte) byte { return OP_JIN | (nibble(rp) &^ 1) }

// istruzioni per la ram

// WRM scrive l'accumulatore nella cella dati RAM selezionata da DCL (banco) e SRC (registro+carattere). Non tocca il carry.
func WRM() byte { return OP_WRM }

// RDM legge la cella dati RAM selezionata nell'accumulatore (Read RAM data). Non tocca il carry.
func RDM() byte { return OP_RDM }

// ADM somma la cella dati RAM selezionata e il carry all'accumulatore: A = A + RAM + C (Add from Memory).
// Imposta C se il risultato supera 0xF.
func ADM() byte { return OP_ADM }

// SBM sottrae la cella dati RAM selezionata dall'accumulatore con prestito: A = A + ~RAM + C (Subtract from Memory).
// C=1 in uscita significa "nessun borrow generato".
func SBM() byte { return OP_SBM }

// WMP scrive l'accumulatore sulla porta di output del banco RAM corrente (Write Memory Port). Non tocca il carry.
func WMP() byte { return OP_WMP }

// WR0 scrive l'accumulatore nel nibble di stato 0 del registro RAM selezionato. Non tocca il carry.
func WR0() byte { return OP_WR0 }

// WR1 scrive l'accumulatore nel nibble di stato 1 del registro RAM selezionato. Non tocca il carry.
func WR1() byte { return OP_WR1 }

// WR2 scrive l'accumulatore nel nibble di stato 2 del registro RAM selezionato. Non tocca il carry.
func WR2() byte { return OP_WR2 }

// WR3 scrive l'accumulatore nel nibble di stato 3 del registro RAM selezionato. Non tocca il carry.
func WR3() byte { return OP_WR3 }

// RD0 legge il nibble di stato 0 del registro RAM selezionato nell'accumulatore. Non tocca il carry.
func RD0() byte { return OP_RD0 }

// RD1 legge il nibble di stato 1 del registro RAM selezionato nell'accumulatore. Non tocca il carry.
func RD1() byte { return OP_RD1 }

// RD2 legge il nibble di stato 2 del registro RAM selezionato nell'accumulatore. Non tocca il carry.
func RD2() byte { return OP_RD2 }

// RD3 legge il nibble di stato 3 del registro RAM selezionato nell'accumulatore. Non tocca il carry.
func RD3() byte { return OP_RD3 }

// WRR scrive l'accumulatore sulla porta di output del chip ROM (Write ROM port). Non tocca il carry.
func WRR() byte { return OP_WRR }

// WPM scrive l'accumulatore nella program memory (Write Program Memory).
// Nell'emulatore con ROM statica è un no-op. Non tocca il carry.
func WPM() byte { return OP_WPM }

// RDR legge la porta di input del chip ROM nell'accumulatore (Read ROM port). Non tocca il carry.
func RDR() byte { return OP_RDR }
