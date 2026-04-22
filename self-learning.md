# Software Engineer Growth Roadmap

Welcome to the self-learning. Follow the checklist below to track your progress.

---

## 📋 Study Checklist

- [ ] [HTTP Protocol in Go](#http-protocol-in-go)
- [ ] [Linux Fundamentals](#linux-fundamentals)
- [ ] [Web Security (CORS & HTTPS/TLS)](#web-security-cors--httpstls)
- [ ] [TypeScript](#typescript)
- [ ] [Testing in Go](#testing-in-go)
- [ ] [Kubernetes Coordination](#kubernetes-coordination)

## Optional (Up to you)
- [ ] [Git & Version Control](#git--version-control)
- [ ] [Golang Fundamentals](#golang-fundamentals)
- [ ] [Go Standard Project Template](#go-standard-project-template)
- [ ] [Data Structures & Algorithms](#data-structures--algorithms)
- [ ] [Docker & Containers](#docker--containers)
- [ ] [System Monitoring (Prometheus & Grafana)](#system-monitoring-prometheus--grafana)
- [ ] [Performance Testing Fundamentals](#performance-testing-fundamentals)

---

## 📝 Assignment 1 — Individual Study Plan

Before you start learning, create a personal plan for how you'll complete the Study Checklist.

### What to include
- **Which topics** you'll study and in what order
- **When** you'll study them — at least 15–30 mins a day, specify the days and time slot (e.g. "Mon–Fri, 8:00–8:30 AM")
- **Your target date** to finish each topic

### Format
Write your plan in the `README.md` of your repository. There's no required template — structure it however is clearest to you. A simple table works well:

```markdown
| Topic | Scheduled Time | Target Date |
|-------|---------------|-------------|
| HTTP Protocol in Go | Mon–Fri 8:00–8:30 AM | 2026-05-10 |
| Linux Fundamentals  | Mon–Fri 8:00–8:30 AM | 2026-05-17 |
| ...   | ...           | ...         |
```

### Submit
Send your GitHub or GitLab repository link to the chat group. Your README should already contain your plan.

---

## 📝 Assignment 2 — Hands-on Project

> Pick a project idea you're interested in and apply what you learn to it.
> **No idea yet?** Use the guided project below — a real-world flight booking API.

### How it works
Each topic has a task. Complete it using **your own project**, or use the guided project spec if you're unsure what to build.

### ✅ Submit Criteria

When you're done, you should be able to show all of the following:

**API related**
| # | What to show |
|---|---|
| 1 | A running API with the 3 booking endpoints — demo via curl or Postman |
| 2 | A CORS demo — one request blocked, one allowed, explain why |
| 3 | A working TypeScript client or CLI that calls your API — no `any`, typed req/res |
| 4 | Test output showing green for happy path, 409, and 400 cases |
| 5 | Your API and frontend run as systemd services on a Linux VM, accessible over HTTPS — share the public URL, show `systemctl status` output, and include your systemd `.service` config files in your repo |

**Standalone**
| # | What to show |
|---|---|
| 6 | A shell script that reads `bookings.psv`, processes records, and writes `report.txt` — including overdue detection |

> You don't need to finish all topics at once. Submit each one as you complete it by sending your repository link to the chat group.

---

### Guided Project: Flight Booking Service (if you have no idea)

You will build a backend service that handles flight bookings — step by step across each topic.

#### REST API in Go
Build these 3 endpoints in Go:

**`POST /api/v1/bookings`** — Create a pending booking (locks seats 15 min in Redis, publishes `booking.created`)

<details>
<summary>📥 Request / 📤 Response / ❌ Errors</summary>

**Request Body**
```json
{
  "flightId":   "fl_8a2b3c4d-5e6f-7a8b-9c0d",
  "cabin":      "economy",
  "fareClass":  "M",
  "seats":      ["14A"],
  "passengers": [{
    "firstName":     "Somchai",
    "lastName":      "Jaidee",
    "dateOfBirth":   "1990-05-15",
    "nationality":   "THA",
    "passportNumber":"AA1234567",
    "passportExpiry":"2030-12-31"
  }],
  "contactEmail": "somchai@example.com",
  "contactPhone": "+66812345678"
}
```

**201 Created**
```json
{
  "data": {
    "bookingId":      "bk_7f8a9b0c-1d2e-3f4a-5b6c",
    "pnr":            "QL3XF7",
    "status":         "PENDING",
    "flightId":       "fl_8a2b3c4d-5e6f-7a8b-9c0d",
    "flightNumber":   "QL101",
    "origin":         "BKK",
    "destination":    "NRT",
    "departureAt":    "2026-06-01T08:30:00+07:00",
    "cabin":          "economy",
    "fareClass":      "M",
    "seats":          ["14A"],
    "totalAmount":    12500.00,
    "currency":       "THB",
    "paymentDeadline":"2026-04-07T10:15:00Z",
    "createdAt":      "2026-04-07T10:00:00Z"
  }
}
```

**Errors**
| Code | Reason |
|------|--------|
| 409  | Seat already locked by another session |
| 404  | `flightId` not found |
| 422  | Seats count ≠ passengers count |
| 400  | Validation error on required fields |

</details>

---

**`GET /api/v1/bookings/:bookingId`** — Fetch booking by internal UUID (used by Payment Service to validate amount)

<details>
<summary>📥 Request / 📤 Response / ❌ Errors</summary>

**Request**
```
GET /api/v1/bookings/bk_7f8a9b0c-1d2e-3f4a-5b6c
Authorization: Bearer <internal-service-jwt>
```

**200 OK**
```json
{
  "data": {
    "bookingId":      "bk_7f8a9b0c-1d2e-3f4a-5b6c",
    "pnr":            "QL3XF7",
    "status":         "PENDING",
    "totalAmount":    12500.00,
    "currency":       "THB",
    "passengerId":    "ps_abc123",
    "flightId":       "fl_8a2b3c4d",
    "seats":          ["14A"],
    "paymentDeadline":"2026-04-07T10:15:00Z"
  }
}
```

**Errors**
| Code | Reason |
|------|--------|
| 404  | `bookingId` not found |

</details>

---

**`GET /api/v1/bookings/pnr/:pnr`** — Fetch booking by 6-char PNR (used by check-in, airport staff, customer portal)

<details>
<summary>📥 Request / 📤 Response</summary>

**Request**
```
GET /api/v1/bookings/pnr/QL3XF7
```

**200 OK** — same shape as `GET /:bookingId`
```json
{
  "data": { "pnr": "QL3XF7", "status": "CONFIRMED", "..." : "..." }
}
```

</details>

---

#### Web Security CORS
Add CORS to your API so that only your frontend origin is allowed. Show the difference between a blocked and allowed request.

#### TypeScript with Bun
Build a **Web frontend page** or **CLI tool** using Bun that calls your booking API — typed request/response, no `any`. Then deploy it online so it's publicly accessible.

Pick one and follow the matching path below:

---

**Option A: Web**

An HTML page with a form that submits a booking and shows the response.

- Use `bun build` to bundle TypeScript for the browser
- Copy the output to your Debian VM and let nginx serve it as static files (this will be done together with the Deployment task)

> Once live, your frontend URL becomes your allowed CORS origin. Update your Go API's CORS config to allow it — this is where the CORS task becomes real.

---

**Option B: CLI**

A terminal tool (e.g. `bun run book.ts --flight fl_xxx --seat 14A`) that calls the API and prints the result.

- Run directly with `bun run` — no compile step needed
- Make it installable so others can run it on their machine. Options:
  - **Publish to npm** — `bun publish` then others run `npx your-tool`
  - **Compile to binary** — `bun build --compile book.ts --outfile booking-cli` produces a single executable, share it via GitHub Releases
  - **Install script** — write a `install.sh` that downloads the binary and puts it in `$PATH`

> Goal: someone else should be able to install and run your CLI without cloning your repo.

#### Testing
Write tests for your booking endpoints. You decide how many types of testing to apply — pick what makes sense for your implementation.

**Required cases to cover (regardless of testing type):**
- Happy path — booking created successfully (201)
- Seat conflict — same seat booked twice (409)
- Missing fields — required fields omitted (400)

**Testing types you can choose from:**
- **Unit test** — test handler logic in isolation with mocked dependencies
- **Integration test** — test the full request/response cycle against a real DB/Redis
- **Both** — unit for logic, integration for the full flow

> There's no single right answer. Think about what gives you the most confidence that your code works.

#### Deployment on a Linux VM

Deploy your API and frontend on a free cloud Linux VM, managed by systemd, served over HTTPS via nginx.

Any Linux distro that supports systemd works (suggestion: Debian) — pick based on what your cloud provider offers in the free tier.

**Get a free VM:**

| Provider | Free Tier |
|----------|-----------|
| [Oracle Cloud](https://www.oracle.com/cloud/free/) | Always-free AMD VM (no expiry) — best option |
| [Google Cloud](https://cloud.google.com/free) | e2-micro in US regions — always free |
| [AWS](https://aws.amazon.com/free/) | t2.micro — free for 12 months |
| [Linode/Akamai](https://www.linode.com/lp/free-credit/) | $100 credit for new accounts |
| [Civo](https://www.civo.com/) | $250 credit for new accounts |

**What to set up on the VM:**
1. Install your Linux distro of choice, Go, and nginx
2. Run your **Go API** as a systemd service — auto-starts on boot, auto-restarts on crash
3. Build your **frontend** with `bun build` and serve the static files via nginx (if Web), or skip if CLI
4. Configure **nginx** as a reverse proxy in front of the API (and static files for Web)
5. Enable **HTTPS** with Let's Encrypt via `certbot` — free SSL certificate

**Domain for HTTPS:**
- **Have your own domain?** Point it to your VM IP and use it — Let's Encrypt will issue a cert for it via `certbot`
- **No domain?** Use `sslip.io` — a free wildcard DNS that maps any IP to a domain, no registration needed. If your VM IP is `1.2.3.4`, your domain is `1-2-3-4.sslip.io` and certbot works the same way

> Your API and frontend must be reachable over HTTPS from the public internet when done.

---

#### Linux with shell script
Write a shell script that reads booking records from a file, processes them, and writes a summary report.

Your script should be able to:
- **Read** — load `bookings.psv` (pipe-separated, one booking per line)
- **Process** — count bookings by status (`PENDING`, `CONFIRMED`, `CANCELLED`), sum total revenue from confirmed bookings, and list any bookings past their `paymentDeadline` that are still `PENDING`
- **Write** — output a `report.txt` with the results in a clean, readable format

Example `bookings.psv` input (pipe-separated):
```
bookingId|pnr|status|totalAmount|currency|paymentDeadline
bk_001|QL3XF7|CONFIRMED|12500.00|THB|2026-04-07T10:15:00Z
bk_002|AB1234|PENDING|8900.00|THB|2026-03-01T10:00:00Z
bk_003|ZX9900|CANCELLED|5000.00|THB|2026-04-10T08:00:00Z
```

Example `report.txt` output:
```
=== Booking Report (generated: 2026-04-20 09:00:00) ===

Status Summary:
  CONFIRMED : 1
  PENDING   : 1
  CANCELLED : 1

Total Revenue (CONFIRMED): 12500.00 THB

Expired PENDING Bookings (payment overdue):
  bk_002  AB1234  deadline: 2026-03-01T10:00:00Z
```

**Hint:** Use `awk -F'|'` to parse fields, `date -u +%s` to convert timestamps to Unix epoch for comparison, and redirect output with `>` to write the report file. No external tools needed — pure shell.

---

## 📖 Learning Resources

### Linux Fundamentals
- **Resource**: [Course: Learn Linux](https://www.boot.dev/courses/learn-linux)
- **Objective**: Understand the OS layer, shell usage, and file systems.
- **Why**: Almost all services run on Linux.

### Git & Version Control
- **Resource**: [Course: Git for Software Engineers](https://www.sphere.academy/classroom/?class=git-for-software-engineers)
- **Objective**: Mastery of version control, branching strategies, and collaboration.
- **Why**: Git is the foundation of modern team-based development.

### HTTP Protocol in Go
- **Resource**: [Course: Learn the HTTP Protocol in GO](https://www.boot.dev/courses/learn-http-protocol-golang)
- **Objective**: Deep dive into headers, methods, status codes, and the request-response lifecycle.

### Web Security (CORS & HTTPS/TLS)
- **Resource 1**: [VDO: How HTTPS Works](https://youtu.be/UIcCwuYzxcE) (TLS Handshake)
- **Resource 2**: [VDO: CORS Explained](https://youtu.be/cGg7aRcIm8o) & [Mozilla CORS Docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- **Objective**: Understand how to secure data and handle cross-origin requests.

### Golang Fundamentals
- **Resource**: [Golang Fundamentals](https://www.sphere.academy/course/?name=golang-fundamentals)
- **Objective**: Learn the core syntax, types, and concurrency model of the Go language.

### Testing in Go
- **Resource**: [Tips: Testing by Example](https://youtu.be/1-o-iJlL4ak)
- **Objective**: Learn how to test your code in Go.

### Go Standard Project Template
- **Resource**: [Go Standard Template](https://www.sphere.academy/course/?name=go-standard-template-pallat)
- **Objective**: Learn standard project structures and patterns used in professional Go environments.

### Data Structures & Algorithms
- **Resource**: [Learn DSA in Python](https://www.boot.dev/courses/learn-data-structures-and-algorithms-python)
- **Objective**: Strengthen problem-solving skills and algorithmic efficiency.

### Functional Programming Concepts
- **Resource**: [Learn Functional Programming in Python](https://www.boot.dev/courses/learn-functional-programming-python)
- **Objective**: Understand immutability, pure functions, and high-order logic.

### Docker & Containers
- **Resource**: [Docker Guide for Beginners](https://youtu.be/pTFZFxd4hOI)
- **Objective**: Learn to package and run applications in isolated environments.

### Kubernetes Coordination
- **Resource**: [VDO: Kubernetes Is Not A Deployment Tool](https://youtu.be/iVj5vEnbFr0)
- **Resource**: [Kubernetes Tutorials](https://youtu.be/7bA0gTroJjw)
- **Objective**: Understand orchestration and the distinction between deployment and control planes.

### System Monitoring (Prometheus & Grafana)
- **Resources**: [Prometheus](https://prometheus.io/) & [Grafana](https://grafana.com/)
- **Objective**: Learn to collect metrics and build dashboards to visualize system health.

### Performance Testing Fundamentals
- **Resource**: [Course: Performance Testing Fundamentals](https://www.sphere.academy/classroom/?class=Performance-Testing-Fundamentals)
- **Objective**: Learn to stress test and find bottlenecks in your services.

### TypeScript
- **Topic**: Typed JavaScript for robust frontend development.
- **Resource**: [TypeScript (TH)](https://www.youtube.com/live/8TzWP4GgSwM?t=555)
- **Resource**: [TypeScript Official Docs](https://www.typescriptlang.org/docs/)

### Architecture Documentation (D2Lang)
- **Topic**: Declarative Diagramming (Text-to-Diagram).
- **Resource**: [D2Lang Documentation](https://d2lang.com/)
