# Real-Time Order & Job Tracking System

Aplikasi untuk mengelola dan memantau status pesanan atau pekerjaan secara real-time menggunakan Go, React, PostgreSQL, dan WebSocket.

## 1. Tech Stack

### Backend
- Go
- Gin Framework
- GORM
- PostgreSQL
- Gorilla WebSocket

### Frontend
- React
- Vite

### Infrastructure
- Docker
- Docker Compose

## 2. Project Structure

```text
study-case-1/
├── backend/
│   ├── main.go
│   ├── orders.go
│   ├── server
│   └── ...
├── frontend/
│   ├── src/
│   ├── package.json
│   └── ...
├── docker-compose.yaml
├── README.md
└── System-Design-Answers.pdf
```

## 3. Requirements

Pastikan sudah menginstal:
- Docker Desktop
- Git

## 4. How to Run

Clone repository:

```bash
git clone https://github.com/Herazor/real-time-order-tracking.git
cd real-time-order-tracking
```

Jalankan seluruh aplikasi menggunakan satu perintah:

```bash
docker compose up --build
```

Docker Compose akan menjalankan:
- PostgreSQL
- Go Backend API
- React Frontend

Akses aplikasi:

| Service | URL |
|---|---|
| Frontend | http://localhost:5173 |
| Backend | http://localhost:8080 |
| PostgreSQL | Internal Docker network |

Untuk menghentikan aplikasi:

```bash
docker compose down
```

Untuk menghentikan aplikasi sekaligus menghapus data database:

```bash
docker compose down -v
```

**Catatan:** Perintah kedua akan menghapus volume database.

## 5. Main Features

- Create new orders.
- Retrieve all orders.
- Retrieve order details.
- Update order status.
- Cancel orders.
- Review cancellation requests.
- Real-time order status updates using WebSocket.
- Order event history.
- Event sequence tracking.
- Optimistic concurrency control using version checking.

## 6. API Endpoints

| Method | Endpoint | Description |
|---|---|---|
| POST | `/api/orders` | Create a new order |
| GET | `/api/orders` | Retrieve all orders |
| GET | `/api/orders/:id` | Retrieve order details |
| GET | `/api/orders/:id/events` | Retrieve order event history |
| PATCH | `/api/orders/:id/status` | Update order status |
| POST | `/api/orders/:id/cancel` | Submit cancellation request |
| PATCH | `/api/cancellations/:id/review` | Review cancellation request |
| GET | `/ws/orders/:id` | Connect to real-time order updates |

## 7. Real-Time Communication

The application uses WebSocket to broadcast order status changes to connected clients.

When an order status changes:
1. The backend processes the status update.
2. The database stores the updated order.
3. An order event is recorded.
4. Connected WebSocket clients receive the update.

The frontend can receive updates without refreshing the page.

## 8. Implementation Status

| Feature | Status |
|---|---|
| REST API | Implemented |
| PostgreSQL integration | Implemented |
| Order event history | Implemented |
| WebSocket updates | Implemented |
| Event sequence tracking | Implemented |
| Optimistic concurrency | Implemented |
| Docker Compose | Implemented |
| Authentication and authorization | Not implemented |
| Idempotency request deduplication | Not implemented |
| Automatic WebSocket reconnection and replay | Not implemented |
| Redis Pub/Sub for multi-instance broadcasting | Not implemented |
| Transactional Outbox | Not implemented |

## 9. Known Limitations

- Authentication and role-based authorization are not fully implemented.
- Demo user IDs are used by the frontend.
- Idempotency keys are not yet used for request deduplication.
- WebSocket broadcasting uses an in-memory connection registry, so broadcasting is limited to a single backend instance.
- Automatic WebSocket reconnection and event replay are not implemented.
- Redis Pub/Sub is not integrated.
- Transactional Outbox is not implemented.
- The backend Docker image uses a precompiled Linux binary included in the repository.

## 10. Concurrency Control

The system uses optimistic concurrency control.

Each order has a version number. A status update includes the expected version. The backend checks the current version before updating the order.

This helps prevent conflicting updates from overwriting each other.

## 11. Event History and Auditability

Order status changes are recorded in the order event history.

Each event contains information such as:
- Order ID
- Event type
- Previous status
- New status
- Sequence number
- Actor ID
- Timestamp

This allows the system to maintain a history of order changes.

## 12. Additional Documentation

See `System-Design-Answers.pdf` for explanations covering:
- Real-time communication strategy
- Handling missed events
- Idempotency and event ordering
- Concurrency control
- Scaling strategy
- Auditability

## 13. Notes

This project is a study case implementation. Some advanced production features are documented as future improvements and are not yet implemented.
