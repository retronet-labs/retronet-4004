# --- Build stage ---
FROM golang:1.26-alpine AS builder
WORKDIR /app

# Copia i sorgenti e compila la CLI in un binario statico.
COPY . .
RUN CGO_ENABLED=0 go build -o retronet-4004 ./cmd/retronet-4004

# --- Runtime stage ---
FROM alpine:latest
WORKDIR /app

# Solo il binario e le ROM di esempio: immagine minima.
COPY --from=builder /app/retronet-4004 .
COPY --from=builder /app/testdata ./testdata

# Di default esegue la demo BCD col trace. Override degli argomenti:
#   docker run <img> -dump-ram testdata/somma-multicifra.rom
ENTRYPOINT ["./retronet-4004"]
CMD ["-trace", "testdata/bcd-add.rom"]
