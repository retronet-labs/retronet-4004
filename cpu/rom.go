package cpu

// ROM rappresenta la memoria di programma del sistema Intel 4004.
// Contiene le istruzioni che la CPU legge ed esegue sequenzialmente.
//
// Il chip ROM reale (Intel 4001) ha anche una porta I/O a 4 bit
// usata per pilotare le righe della tastiera (WRR) e leggere
// le colonne premute (RDR). Il campo Port emula questa porta.
type ROM struct {
	Data []byte
	Port uint8 // porta I/O del chip ROM (Intel 4001) — usata da WRR e RDR
}

// NewROM crea una ROM a partire da uno slice di byte (il programma).
func NewROM(data []byte) *ROM {
	return &ROM{Data: data}
}
