# Dokumentasi

## Spesifikasi API

Baik EC maupun 2PC, punya endpoint `POST /orders` dengan payload berikut:

```json
{
  "hotel_room_id": "four-seasons-jakarta-101",
  "hotel_room_start_date": "YYYY-MM-DD",
  "hotel_room_end_date": "YYYY-MM-DD",
  "car_id": "toyota-avanza-001",
  "car_start_date": "YYYY-MM-DD",
  "car_end_date": "YYYY-MM-DD",
  "train_seat_id": "argo-bromo-anggrek-1",
  "user_id": "1"
}
```

Semua field wajib diisi. Autentikasi tidak diikutsertakan. Validasi isian tidak dicek oleh server, melainkan data uji sudah dipastikan valid.

## Metodologi

1. Implementasi 2PC dan EC, masing-masing terdiri dari 4 services: orders, car, hotel, dan train service
2. Pengujian
3. Pengumpulan data

## Pengujian

Langkah-langkah umum dalam pengujian:

1. Seeding data menggunakan script
2. Jalankan subtes
3. Catat hasil pengukuran
4. Hapus data pada database

### Subtes A: Load Testing (Mendapatkan staleness time, troughput, latency, dan komponen yang mengakibatkan latency)

1. Pakai JSR223 buat bikin dynamic request (start/end date dibuat konstan 2025-06-28)
2. hotel_room_id, car_id, dan train_seat_id diambil dari DB, sedangkan user_id akan diiterasi dari 1 hingga N.
3. Hit endpoint `POST /orders` menggunakan parameter dari poin 1 dan 2, pastikan untuk setiap request, kombinasi id yang ada unique.
4. Throughput/latency didapatkan langsung dari JMeter, staleness time didapatkan dari selisih waktu antara created_at dan done_at pada tabel orders (waktu untuk mencapai konsistensi atau berapa lama transaksi tersebut diproses)
5. Komponen yang mengakibatkan latency dapat diukur dari selisih antara created_at dan car_done_at, hotel_done_at, dan train_done_at pada tabel order

> Seluruh step di atas akan dilakukan dengan 100, 500, 1000, dan 5000 concurrent requests

### Subtes B: Consistency testing

1. Pakai JSR223 buat bikin dynamic request (start/end date dibuat konstan 2025-06-28)
2. Setiap kombinasi hotel_room_id, car_id, dan train_seat_id akan disimulasikan diorder oleh 10 user yang berbeda pada tanggal yang sama
3. Hit endpoint `POST /orders` menggunakan parameter dari poin 1 dan 2
4. Dilakukan query ke tabel transaksi pada masing-masing services, pastikan hotel_room_id, car_id, dan train_seat_id tidak ada yang memiliki status reserved lebih dari 1 user (hanya boleh 1 user yang bookingnya berhasil)
