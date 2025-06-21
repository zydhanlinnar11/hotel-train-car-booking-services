# Two-Phase Commit Coordinator

Modul ini mengimplementasikan two-phase commit protocol untuk koordinasi transaksi di sistem booking hotel, kereta, dan mobil.

## Fitur

- **Two-Phase Commit Protocol**: Implementasi lengkap prepare-commit protocol
- **Transaction Logging**: Log transaksi disimpan di Firestore dengan prefix `twophase_`
- **Retry Mechanism**: Exponential backoff untuk retry otomatis
- **Timeout Handling**: Penanganan timeout transaksi
- **REST API**: Endpoint untuk koordinasi dan monitoring
- **Health Check**: Endpoint untuk monitoring kesehatan service

## Struktur Proyek

```
twophase/
├── cmd/
│   └── coordinator/
│       └── main.go
├── internal/
│   └── coordinator/
│       ├── model.go
│       ├── repository.go
│       ├── service.go
│       └── handler.go
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## API Endpoints

### Order Management

- `POST /api/orders` - Membuat order dengan two-phase commit
- `GET /api/transactions/:transactionID` - Melihat status transaksi

### Two-Phase Commit (untuk participants)

- `POST /api/twophase/prepare` - Prepare phase
- `POST /api/twophase/commit` - Commit phase
- `POST /api/twophase/abort` - Abort phase

### Health Check

- `GET /api/health` - Status kesehatan service

## Konfigurasi

### Environment Variables

```bash
# Required
GOOGLE_CLOUD_PROJECT=your-project-id

# Optional (dengan default)
PORT=8080
TRANSACTION_TIMEOUT=30s
MAX_RETRIES=3
RETRY_DELAY=2s

# Service URLs
HOTEL_SERVICE_URL=http://localhost:8081
CAR_SERVICE_URL=http://localhost:8082
TRAIN_SERVICE_URL=http://localhost:8083
```

## Cara Menjalankan

### Development

1. Install dependencies:

```bash
go mod download
```

2. Set environment variables:

```bash
export GOOGLE_CLOUD_PROJECT=your-project-id
```

3. Jalankan service:

```bash
go run cmd/coordinator/main.go
```

### Docker

1. Build dan jalankan dengan Docker Compose:

```bash
docker-compose up --build
```

2. Atau build image secara manual:

```bash
docker build -t twophase-coordinator .
docker run -p 8080:8080 -e GOOGLE_CLOUD_PROJECT=your-project-id twophase-coordinator
```

## Contoh Penggunaan

### Membuat Order

```bash
curl -X POST http://localhost:8080/api/orders \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "hotel_id": "hotel-abc",
    "room_id": "room-xyz",
    "car_id": "car-123",
    "train_id": "train-456",
    "seat_id": "seat-789",
    "check_in_date": "2024-02-01T14:00:00Z",
    "check_out_date": "2024-02-03T12:00:00Z",
    "travel_date": "2024-02-01T08:00:00Z",
    "total_price": 1500000
  }'
```

Response:

```json
{
  "order_id": "uuid-order-id",
  "transaction_id": "uuid-transaction-id",
  "status": "initiated",
  "message": "Transaction initiated successfully"
}
```

### Cek Status Transaksi

```bash
curl http://localhost:8080/api/transactions/uuid-transaction-id
```

Response:

```json
{
  "transaction_id": "uuid-transaction-id",
  "order_id": "uuid-order-id",
  "status": "committed",
  "participants": [
    {
      "service_name": "hotel",
      "service_url": "http://localhost:8081",
      "status": "committed",
      "error": "",
      "retry_count": 0
    },
    {
      "service_name": "car",
      "service_url": "http://localhost:8082",
      "status": "committed",
      "error": "",
      "retry_count": 0
    },
    {
      "service_name": "train",
      "service_url": "http://localhost:8083",
      "status": "committed",
      "error": "",
      "retry_count": 0
    }
  ],
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:05Z",
  "timeout_at": "2024-01-15T10:31:00Z",
  "retry_count": 0,
  "failure_reason": ""
}
```

## Status Transaksi

- `initiated` - Transaksi baru dibuat
- `prepared` - Semua participants siap untuk commit
- `committed` - Transaksi berhasil di-commit
- `aborted` - Transaksi di-abort
- `rolled_back` - Transaksi di-rollback
- `timed_out` - Transaksi timeout

## Integrasi dengan Service Lain

Service lain (hotel, car, train) harus mengimplementasikan endpoint two-phase commit:

### Prepare Endpoint

```bash
POST /api/twophase/prepare
{
  "transaction_id": "uuid",
  "order_id": "uuid",
  "service_name": "hotel"
}
```

### Commit Endpoint

```bash
POST /api/twophase/commit
{
  "transaction_id": "uuid",
  "order_id": "uuid",
  "service_name": "hotel"
}
```

### Abort Endpoint

```bash
POST /api/twophase/abort
{
  "transaction_id": "uuid",
  "order_id": "uuid",
  "service_name": "hotel",
  "reason": "Insufficient inventory"
}
```

## Monitoring dan Maintenance

### Cleanup Otomatis

Service secara otomatis membersihkan transaksi yang timeout setiap menit.

### Log Transaksi

Semua transaksi disimpan di Firestore collection `twophase_transactions` dengan struktur:

```json
{
  "id": "transaction-id",
  "order_id": "order-id",
  "status": "committed",
  "participants": [...],
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "timeout_at": "timestamp",
  "retry_count": 0,
  "max_retries": 3,
  "failure_reason": "",
  "commit_timestamp": "timestamp"
}
```

## Troubleshooting

### Transaksi Timeout

- Cek log untuk melihat alasan timeout
- Pastikan semua service participants berjalan
- Verifikasi konektivitas jaringan

### Retry Failures

- Cek status participants di response transaksi
- Verifikasi endpoint two-phase commit di service participants
- Cek log error di service participants

### Firestore Issues

- Pastikan credentials Google Cloud sudah benar
- Verifikasi project ID dan permissions
- Cek konektivitas ke Firestore
