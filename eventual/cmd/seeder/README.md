# Database Seeder

Seeder untuk mengisi data awal ke Firestore database dengan strategi yang telah ditentukan. Seeder menggunakan model yang sudah ada di `internal/*/model.go`.

## Struktur Data

### Car Data

- **Collection**: `cars`
- **Model**: `internal/car/model.go` - `Car`
- **ID**: Slug dari nama mobil
- **Nama**: Format `${brand} ${model} - ${unit_number}`
- **Brand & Model**: 10 brand terkenal, masing-masing 5 model
- **Unit**: 100 unit per model
- **Total**: 5,000 mobil

**Contoh ID**: `toyota-avanza-001`, `honda-crv-050`

### Hotel Room Data

- **Collection**: `hotel_rooms`
- **Model**: `internal/hotel/model.go` - `HotelRoom`
- **ID**: Slug dari nama hotel + nama kamar
- **Hotel Name**: Brand hotel terkenal
- **Nama Kamar**: Format `${floor}${2_digit_unit_number}`
- **Floors**: 5 lantai per hotel
- **Units**: 20 kamar per lantai
- **Total**: 1,500 kamar

**Contoh ID**: `marriott-jakarta-101`, `ritz-carlton-jakarta-520`

### Train Data

- **Collection**: `train_seats`
- **Model**: `internal/train/model.go` - `TrainSeat`
- **ID**: Slug dari nama kereta + seat ID
- **Seat ID**: Auto increment integer (1-500)
- **Train Name**: Nama kereta Indonesia
- **Seats**: 500 kursi per kereta
- **Total**: 5,000 kursi

**Contoh ID**: `argo-bromo-anggrek-1`, `bima-250`

## Cara Menjalankan

1. Pastikan environment variables sudah diset:

   ```bash
   export GOOGLE_PROJECT_ID="your-project-id"
   ```

2. Jalankan seeder:
   ```bash
   go run cmd/seeder/main.go
   ```

## Output

Seeder akan menampilkan progress untuk setiap jenis data:

```
Starting database seeder...
Starting car seeder...
Seeding Toyota Avanza...
Seeding Toyota Innova...
...
Car seeder completed. Total cars: 5000
Starting hotel room seeder...
Seeding Marriott Jakarta...
...
Hotel room seeder completed. Total rooms: 1500
Starting train seeder...
Seeding Argo Bromo Anggrek...
...
Train seeder completed. Total seats: 5000
Database seeding completed successfully!
```

## Data yang Dihasilkan

### Car Brands & Models

- Toyota: Avanza, Innova, Fortuner, Camry, Corolla
- Honda: Brio, Jazz, HR-V, CR-V, Civic
- Suzuki: Ertiga, XL7, Ignis, Baleno, Swift
- Daihatsu: Ayla, Calya, Xenia, Terios, Rocky
- Mitsubishi: Xpander, Pajero, L300, Colt, Mirage
- Nissan: Livina, Grand Livina, X-Trail, Serena, March
- Hyundai: Brio, Creta, Santa Fe, Stargazer, Palisade
- Kia: Picanto, Rio, Seltos, Sportage, Carnival
- Wuling: Almaz, Cortez, Confero, Air ev, Alvez
- MG: ZS, HS, RX5, 5, 3

### Hotel Brands

- Marriott Jakarta, Bandung, Surabaya, Medan
- Ritz-Carlton Jakarta, Bandung, Surabaya
- Mandarin Oriental Jakarta
- Four Seasons Jakarta
- Grand Hyatt Jakarta
- InterContinental Jakarta
- Sheraton Jakarta
- Pullman Jakarta
- Novotel Jakarta
- Ibis Jakarta

### Train Names

- Argo Bromo Anggrek, Argo Lawu, Argo Parahyangan
- Bima, Gajayana, Harina, Kertajaya
- Lodaya, Malabar, Matarmaja

## Model yang Digunakan

Seeder menggunakan model yang sudah ada di internal packages:

- **Car**: `internal/car/model.go` - `Car{ID, Name}`
- **HotelRoom**: `internal/hotel/model.go` - `HotelRoom{ID, HotelName, RoomName}`
- **TrainSeat**: `internal/train/model.go` - `TrainSeat{ID, SeatID, TrainName}`
