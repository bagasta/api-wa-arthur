# 🐳 Deployment dengan Docker

## Struktur File Docker

```
whatsapp-endpoint/
├── Dockerfile           ← Multi-stage build (Go + sqlite3)
├── docker-compose.yml   ← Orchestration + volume
└── .env                 ← Konfigurasi (buat dari .env.example)
```

---

## 1. Persiapan di Server

```bash
# Clone repo
git clone https://github.com/bagasta/api-wa-arthur.git
cd api-wa-arthur

# Buat file .env dari template
cp .env.example .env
nano .env
```

Isi `.env` dengan nilai yang sebenarnya:
```env
DATABASE_URL=postgres://user:pass@host:5432/langchain_db
LANGCHAIN_API_URL=https://api.chiefaiofficer.id
N8N_WEBHOOK_URL=
```

---

## 2. Build & Jalankan

```bash
docker compose up --build -d
```

> ⚠️ Build pertama membutuhkan beberapa menit karena mengompilasi CGO (sqlite3).

---

## 3. Scan QR Code (Pertama Kali)

Karena WhatsApp perlu QR Code scan saat pertama kali, **jalankan tanpa `-d`** dulu:

```bash
docker compose up --build
```

QR Code akan muncul di terminal. Scan menggunakan WhatsApp di HP → **Linked Devices → Link a Device**.

Setelah berhasil (`Successfully paired...`), tekan `Ctrl+C`, lalu jalankan ulang sebagai background:

```bash
docker compose up -d
```

Sesi tersimpan di Docker volume bernama `whatsapp_session_data` dan tetap ada meski container di-restart.

---

## 4. Perintah Docker Berguna

```bash
# Lihat log real-time
docker compose logs -f

# Restart container
docker compose restart

# Stop
docker compose down

# Update (setelah git pull)
docker compose up --build -d
```

---

## 5. Mengganti Nomor WhatsApp di Docker

```bash
# 1. Stop container
docker compose down

# 2. Hapus volume (menghapus store.db)
docker volume rm whatsapp_session_data

# 3. Jalankan ulang & scan QR baru
docker compose up
```

---

## 6. Cek Status

```bash
# Cek container berjalan
docker compose ps

# Test endpoint
curl http://localhost:8101/api/v1/whatsapp/sessions/test
```
