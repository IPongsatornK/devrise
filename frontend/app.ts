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

const apiUrl = window.location.hostname === "localhost"
  ? "http://localhost:8080/api/v1/bookings"
  : (window as any).API_URL ?? "http://localhost:8080/api/v1/bookings"

const form = document.querySelector<HTMLFormElement>("#booking-form")
const result = document.querySelector<HTMLDivElement>('#result')

function showMessage(message: string, isError = false) {
  if (!result) return
  result.textContent = message
  result.className = isError ? 'error' : ''
}

async function handleSubmit(event: Event) {
  event.preventDefault()
  if (!form) return

  const formData = new FormData(form)
  const seat = formData.get('seat')?.toString().trim() ?? ''
  const name = formData.get('name')?.toString().trim() ?? ''
  const email = formData.get('email')?.toString().trim() ?? ''
  const phone = formData.get('phone')?.toString().trim() ?? ''

  if (!seat || !name || !email || !phone) {
    showMessage('Please fill every field.', true)
    return
  }

  const [firstName, lastName] = name.split(' ', 2)
  if (!firstName || !lastName) {
    showMessage('Enter your name as first and last name.', true)
    return
  }

  const request: BookingRequest = {
    flightId: 'fl_8a2b3c4d-5e6f-7a8b-9c0d',
    cabin: 'economy',
    fareClass: 'M',
    seats: [seat.toUpperCase()],
    passengers: [
      {
        firstName,
        lastName,
        dateOfBirth: '1990-05-15',
        nationality: 'THA',
        passportNumber: 'AA1234567',
        passportExpiry: '2030-12-31',
      },
    ],
    contactEmail: email,
    contactPhone: phone,
  }

  showMessage('Sending booking request...')

  try {
    const response = await fetch(apiUrl, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(request),
    })

    const payload = (await response.json()) as ApiResponse<BookingData>
    if (!response.ok || payload.error) {
      showMessage(`Booking failed: ${payload.error ?? response.statusText}`, true)
      return
    }

    showMessage(`Booking created successfully. PNR: ${payload.data?.pnr ?? 'unknown'}`)
  } catch (error) {
    showMessage(`Request failed: ${(error as Error).message}`, true)
  }
}

if (form) {
  form.addEventListener('submit', handleSubmit)
}
