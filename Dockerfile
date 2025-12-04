# --- TAHAP 1: Build ---
# Kita pakai image Go resmi untuk meng-compile kode
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy file dependensi dulu (agar cache optimal)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh kode sumber
COPY . .

# Build aplikasi menjadi binary bernama 'main'
RUN go build -o main main.go

# --- TAHAP 2: Running ---
# Kita pindahkan hasil masakan ke piring kecil (Alpine Linux)
FROM alpine:latest

WORKDIR /root/

# Copy binary hasil build dari Tahap 1
COPY --from=builder /app/main .

# PENTING: Copy folder docs agar Swagger bisa jalan!
COPY --from=builder /app/docs ./docs

# Expose port (hanya untuk dokumentasi, Railway otomatis detect)
EXPOSE 8080

# Jalankan aplikasi
CMD ["./main"]