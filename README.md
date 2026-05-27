# 🚀 DevRise - Individual Study Plan

Welcome! This is my personal learning roadmap for the Software Engineer Growth Program.

---

## 📅 Program Overview

**Study Frequency:** 
- **Weekdays (Mon–Fri):** 8:00–8:30 AM (30 minutes daily)
- **Weekends (Sat–Sun):** 9:00–10:00 AM (1 hour daily)

**Program Timeline:** April 24, 2026 – May 25, 2026 (5 weeks)  
**Total Study Days:** 30 days (20 weekdays + 10 weekends)  
**Total Study Hours:** ~35 hours

---

## 📋 Detailed Study Plan 

### Topics

| # | Topic | Duration | Scheduled Time | Start Date | End Date | Status |
|---|-------|----------|---|---|---|---|
| 1 | Linux Fundamentals | 4 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | Apr 24 | Apr 27 | ⏳ |
| 2 | Golang Fundamentals | 5 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | Apr 28 | May 2 | ⏳ |
| 3 | HTTP Protocol in Go | 4 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | May 3 | May 6 | ⏳ |
| 4 | Web Security (CORS & HTTPS/TLS) | 5 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | May 7 | May 10 | ⏳ |
| 5 | Kubernetes Coordination | 4 days | Mon–Fri 8:00–8:30 AM | May 11 | May 14 | ⏳ |
| 6 | Testing in Go | 4 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | May 15 | May 18 | ⏳ |
| 7 | Docker & Containers | 5 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | May 19 | May 23 | ⏳ |
| 8 | TypeScript | 2 days | Mon–Fri 8:00–8:30 AM, Sat–Sun 9:00–10:00 AM | May 24 | May 25 | ⏳ |


---

## ✅ Progress Tracking

**Mark completion with:**
- ⏳ = Not started
- 🔄 = In progress  
- ✅ = Completed
- ❌ = Skipped/Deferred

**Weekly Check-ins:** Review progress every Sunday evening

---

## 🔗 Quick Resource Links

**Core Topics:**
- Linux: [boot.dev/linux](https://www.boot.dev/courses/learn-linux)
- Golang: [sphere.academy/golang](https://www.sphere.academy/course/?name=golang-fundamentals)
- HTTP: [boot.dev/http](https://www.boot.dev/courses/learn-http-protocol-golang)
- Security: [MDN CORS](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- TypeScript: [typescriptlang.org](https://www.typescriptlang.org/docs/)
- Testing: [YouTube Testing](https://youtu.be/1-o-iJlL4ak)
- Docker: [YouTube Docker](https://youtu.be/pTFZFxd4hOI)
- Kubernetes: [YouTube K8s](https://youtu.be/7bA0gTroJjw)

---

## 🛫 Assignment 2 — Flight Booking Service

This repository now includes a working Go booking API, a TypeScript CLI client, request validation, CORS support, tests, and deployment templates.

### What is included

- `main.go` — Go HTTP API with endpoints:
  - `POST /api/v1/bookings`
  - `GET /api/v1/bookings/:bookingId`
  - `GET /api/v1/bookings/pnr/:pnr`
- `bookings_test.go` — Go tests for happy path, seat conflict, and validation error
- `booking-cli.ts` — Typed TypeScript CLI for creating bookings
- `process_bookings.sh` — Shell script for `bookings.psv` report generation
- `bookings.psv` — sample booking data for report generation
- `systemd/booking-api.service` and `systemd/booking-frontend.service` — Linux service templates

### Run the API locally

```sh
cd c:\Users\fookky\devrise-repo\backend
go run .
```

### Run Go tests

```sh
cd c:\Users\fookky\devrise-repo\backend
go test ./...
```

### Example booking request

```sh
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Content-Type: application/json" \
  -d '{
    "flightId":"fl_8a2b3c4d-5e6f-7a8b-9c0d",
    "cabin":"economy",
    "fareClass":"M",
    "seats":["14A"],
    "passengers":[{
      "firstName":"Somchai",
      "lastName":"Jaidee",
      "dateOfBirth":"1990-05-15",
      "nationality":"THA",
      "passportNumber":"AA1234567",
      "passportExpiry":"2030-12-31"
    }],
    "contactEmail":"somchai@example.com",
    "contactPhone":"+66812345678"
  }'
```

### Run the TypeScript CLI

```sh
bun run ./booking-cli.ts --flight fl_8a2b3c4d-5e6f-7a8b-9c0d --seat 14A --name "Somchai Jaidee" --email somchai@example.com --phone "+66812345678"
```

### Run the Web Frontend

1. Build the frontend bundle:

```sh
bun run build:frontend
```

2. Serve the `frontend/` directory on a static host such as `http-server`, `python3 -m http.server`, or your Linux VM web server.

3. Open `frontend/index.html` in the browser.

> If you use a live origin, add that origin to the API `ALLOWED_ORIGIN` environment variable.

### Generate the booking report

```sh
sh process_bookings.sh bookings.psv report.txt
```

### Notes

- The API includes CORS middleware and allows `http://localhost:5173` by default.
- Requests from any origin not on the allow list are blocked by CORS.
- Seat locks are held for 15 minutes to simulate Redis-style reservation behavior.
- Additional deployment templates are available under `systemd/`.
