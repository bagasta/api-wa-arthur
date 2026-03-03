# 🔧 Manajemen Sesi & Maintenance

## Mengganti Nomor WhatsApp

Jika kamu ingin mengganti nomor WhatsApp yang dipakai bot (misal: pindah ke nomor baru atau nomor bisnis), cukup **hapus file `store.db`** yang menyimpan sesi aktif:

```bash
# 1. Hentikan aplikasi terlebih dahulu (Ctrl+C di terminal)

# 2. Hapus file sesi lama
rm store.db

# 3. Jalankan ulang — QR Code baru akan muncul
go run .
```

Scan QR Code yang tampil dengan nomor WhatsApp **baru** yang kamu inginkan.

> ⚠️ **Perhatian:** Data `/connect` di tabel `whatsapp_sessions` (PostgreSQL) **tidak ikut terhapus**. Pengguna lama yang sebelumnya sudah `/connect` perlu melakukan `/connect` ulang dari nomor baru.

---

## Reset Sesi Login (Logout)

Sama seperti ganti nomor, cukup hapus `store.db`:

```bash
rm store.db
```

---

## Restart Aplikasi

Untuk menerapkan perubahan kode atau konfigurasi:

```bash
# Hentikan proses yang berjalan
Ctrl+C

# Jalankan ulang
go run .
```

Jika sesi WhatsApp masih valid (`store.db` ada), aplikasi akan langsung terhubung tanpa perlu scan QR ulang.

---

## Lokasi File Penting

| File | Fungsi |
|------|--------|
| `store.db` | Menyimpan sesi login WhatsApp (kunci enkripsi, device ID) |
| `.env` | Konfigurasi environment (DATABASE_URL, API URL, dll) |
| `go.mod` | Daftar dependency Go |

---

## Cek Log

Saat berjalan, log akan tampil di terminal dengan format:

```
2026/03/03 14:20:51 Received WA Message from 6288xxx@s.whatsapp.net: halo
2026/03/03 14:20:51 WARN: N8N_WEBHOOK_URL is not set. Directly processing internally...
```

- `@s.whatsapp.net` → pesan dari personal chat (diproses)
- `@g.us` → pesan dari grup (diabaikan otomatis)

---

## Troubleshooting

### QR Code tidak muncul / langsung error
Coba hapus `store.db` dan jalankan ulang:
```bash
rm store.db && go run .
```

### "address already in use" saat start
Port `8101` sudah dipakai proses lain. Matikan proses lama:
```bash
kill -9 $(lsof -t -i:8101)
go run .
```

### Bot tidak merespons pesan
1. Pastikan sudah `/connect <AGENT_ID>` terlebih dahulu.
2. Pastikan `LANGCHAIN_API_URL` di `.env` benar dan dapat diakses.
3. Cek log terminal untuk error dari Langchain API.

### "Got device removed stream error"
Sesi WhatsApp dihapus dari HP (misalnya di *Linked Devices → Hapus perangkat*). Hapus `store.db` dan scan QR ulang:
```bash
rm store.db && go run .
```
