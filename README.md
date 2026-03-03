# 📱 WhatsApp Endpoint Service (Go + Whatsmeow)

Layanan **standalone** yang menjadi jembatan antara **WhatsApp** dan **Langchain API**. Polanya sama seperti versi Telegram, yaitu **"Single Bot, Multiple Agents"** dengan _Auto-Fetch API Key_. Ditulis dengan **Go** menggunakan library `whatsmeow` agar lebih ringan dan efisien via WebSocket.

---

## 📚 Dokumentasi Lengkap

Semua dokumentasi tersimpan di folder [`docs/`](./docs/):

| File | Isi |
|------|-----|
| [📖 01-setup.md](./docs/01-setup.md) | Instalasi, konfigurasi `.env`, cara pertama kali menjalankan & scan QR Code |
| [💬 02-usage.md](./docs/02-usage.md) | Perintah bot (`/connect`, `/disconnect`, `/status`), cara chat dengan agent |
| [🔧 03-maintenance.md](./docs/03-maintenance.md) | Mengganti nomor WA, reset sesi, restart, troubleshooting |

---

## ⚡ Quick Start

```bash
# 1. Masuk ke folder
cd /home/bagas/telegram-endpoint/whatsapp-endpoint

# 2. Salin & isi konfigurasi
cp .env.example .env

# 3. Jalankan
go run .
```

Scan QR Code yang muncul di terminal menggunakan WhatsApp → **Linked Devices → Link a Device**.

Lihat [docs/01-setup.md](./docs/01-setup.md) untuk panduan lengkap.
