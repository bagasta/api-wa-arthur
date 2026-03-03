# 📱 Cara Penggunaan WhatsApp Bot

## Perintah yang Tersedia

Kirim perintah ini melalui **personal chat** langsung ke nomor bot:

| Perintah | Fungsi |
|----------|--------|
| `/start` | Menampilkan pesan sambutan & panduan |
| `/connect <AGENT_ID>` | Menghubungkan chat ini ke AI Agent tertentu |
| `/disconnect` | Memutus koneksi dari agent aktif |
| `/status` | Melihat agent mana yang sedang terhubung |
| _(pesan biasa)_ | Diteruskan ke AI Agent yang sedang aktif |

---

## Cara Menghubungkan ke Agent

1. Kirim perintah `/connect` diikuti UUID Agent:
   ```
   /connect 7766af90-c71f-4f6d-8cea-897fc0d27797
   ```
2. Bot akan otomatis mengambil API Key dari database (**Auto-Fetch**).
3. Jika berhasil, bot membalas konfirmasi berisi nama agent yang terhubung.

> **Catatan:** API Key tidak perlu diketik manual. Sistem mengambilnya otomatis dari tabel `api_keys` di database berdasarkan `agent_id`.

---

## Cara Chat dengan Agent

Setelah `/connect` berhasil, cukup kirim pesan biasa:

```
Halo, siapa kamu?
Analisa laporan keuangan Q1 2026
Buatkan email profesional untuk vendor kami
```

Bot akan menampilkan indikator **"typing..."** (sedang mengetik) saat memproses, kemudian mengirim balasan dari agent.

---

## Memutus Koneksi

```
/disconnect
```

Setelah disconnect, pesan biasa tidak akan diproses sampai `/connect` ulang ke agent.

---

## Cek Status

```
/status
```

Bot membalas dengan Agent ID yang sedang aktif, atau memberi tahu jika belum terhubung.

---

## Batasan

- Bot **hanya merespons pesan dari personal/private chat**.
- Pesan dari **grup WhatsApp diabaikan** sepenuhnya (tidak ada balasan).
- Satu nomor WhatsApp hanya bisa terhubung ke **satu agent** dalam satu waktu.
