type Passenger = {
  firstName: string
  lastName: string
  dateOfBirth: string
  nationality: string
  passportNumber: string
  passportExpiry: string
}

type BookingRequest = {
  flightId: string
  cabin: string
  fareClass: string
  seats: string[]
  passengers: Passenger[]
  contactEmail: string
  contactPhone: string
}

type BookingData = {
  bookingId: string
  pnr: string
  status: string
  flightId: string
  flightNumber: string
  origin: string
  destination: string
  departureAt: string
  cabin: string
  fareClass: string
  seats: string[]
  totalAmount: number
  currency: string
  paymentDeadline: string
  createdAt: string
}

type ApiResponse<T> = {
  data?: T
  error?: string
}

const apiUrl = process.env.API_URL ?? "http://localhost:8080/api/v1/bookings"

function parseArgs(argv: string[]) {
  const args = new Map<string, string>()
  let currentKey: string | null = null
  for (const token of argv) {
    if (token.startsWith("--")) {
      currentKey = token.slice(2)
      args.set(currentKey, "")
    } else if (currentKey) {
      args.set(currentKey, token)
      currentKey = null
    }
  }
  return args
}

function requiredArg(args: Map<string, string>, key: string) {
  const value = args.get(key) ?? ""
  if (!value) {
    throw new Error(`Missing required option --${key}`)
  }
  return value
}

async function run() {
  const args = parseArgs(process.argv.slice(2))
  const flightId = requiredArg(args, "flight")
  const seat = requiredArg(args, "seat")
  const name = requiredArg(args, "name")
  const email = requiredArg(args, "email")
  const phone = requiredArg(args, "phone")

  const [firstName, lastName] = name.split(" ", 2)
  if (!firstName || !lastName) {
    throw new Error("Please provide passenger name as two words: first last")
  }

  const request: BookingRequest = {
    flightId,
    cabin: "economy",
    fareClass: "M",
    seats: [seat],
    passengers: [
      {
        firstName,
        lastName,
        dateOfBirth: "1990-05-15",
        nationality: "THA",
        passportNumber: "AA1234567",
        passportExpiry: "2030-12-31",
      },
    ],
    contactEmail: email,
    contactPhone: phone,
  }

  const response = await fetch(apiUrl, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(request),
  })

  const payload = (await response.json()) as ApiResponse<BookingData>
  if (!response.ok || payload.error) {
    console.error("Booking failed:", payload.error ?? response.statusText)
    process.exit(1)
  }

  console.log("Booking created successfully:")
  console.log(JSON.stringify(payload.data, null, 2))
}

run().catch((error) => {
  console.error(error.message)
  process.exit(1)
})
