# 🛠 Setup & Instalasi WhatsApp Endpoint

## Prasyarat
- Go 1.21+
- PostgreSQL (shared dengan Langchain API)
- `sqlite3` library (sudah diinclude via `go-sqlite3`)

---

## 1. Clone & Masuk ke Folder

```bash
cd /home/bagas/telegram-endpoint/whatsapp-endpoint
```

## 2. Konfigurasi Environment

```bash
cp .env.example .env
nano .env
```

Isi variabel berikut:

| Key | Keterangan | Wajib? |
|-----|------------|--------|
| `DATABASE_URL` | PostgreSQL connection string (shared dengan Langchain API) | ✅ |
| `LANGCHAIN_API_URL` | Base URL Langchain API | ✅ |
| `N8N_WEBHOOK_URL` | URL webhook N8N (jika pakai N8N) | ❌ |
| `PORT` | Port HTTP server (default: `8101`) | ❌ |

Contoh `.env`:
```env
DATABASE_URL=postgres://user:pass@localhost:5432/langchain_db
LANGCHAIN_API_URL=https://api.chiefaiofficer.id
N8N_WEBHOOK_URL=
PORT=8101
```

## 3. Download Dependencies

```bash
go mod tidy
```

## 4. Jalankan Aplikasi

```bash
go run .
```

---

## 5. Scan QR Code WhatsApp

Saat **pertama kali** dijalankan, QR Code WhatsApp akan tampil langsung di terminal:

```
====================
SCAN QR CODE THIS IN WHATSAPP:
█████████████████████████████
█ ▄▄▄▄▄ █▀█ █▄ ▄█ ▄▄▄▄▄ █
...
====================
```

**Langkah scan:**
1. Buka WhatsApp di HP kamu
2. Masuk ke **Settings → Linked Devices → Link a Device**  
3. Arahkan kamera ke QR Code di terminal
4. Tunggu konfirmasi "Successfully paired..."

Setelah berhasil, sesi akan tersimpan di file `store.db` dan aplikasi tidak akan minta QR Code lagi saat restart.

---

## 6. Verifikasi Berjalan

```bash
curl http://localhost:8101/api/v1/whatsapp/sessions/test
```

Jika server aktif, akan mengembalikan response JSON.
