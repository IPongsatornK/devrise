package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type BookingRequest struct {
	FlightID     string      `json:"flightId"`
	Cabin        string      `json:"cabin"`
	FareClass    string      `json:"fareClass"`
	Seats        []string    `json:"seats"`
	Passengers   []Passenger `json:"passengers"`
	ContactEmail string      `json:"contactEmail"`
	ContactPhone string      `json:"contactPhone"`
}

type Passenger struct {
	FirstName      string `json:"firstName"`
	LastName       string `json:"lastName"`
	DateOfBirth    string `json:"dateOfBirth"`
	Nationality    string `json:"nationality"`
	PassportNumber string `json:"passportNumber"`
	PassportExpiry string `json:"passportExpiry"`
}

type BookingData struct {
	BookingID       string   `json:"bookingId"`
	PNR             string   `json:"pnr"`
	Status          string   `json:"status"`
	FlightID        string   `json:"flightId"`
	FlightNumber    string   `json:"flightNumber"`
	Origin          string   `json:"origin"`
	Destination     string   `json:"destination"`
	DepartureAt     string   `json:"departureAt"`
	Cabin           string   `json:"cabin"`
	FareClass       string   `json:"fareClass"`
	Seats           []string `json:"seats"`
	TotalAmount     float64  `json:"totalAmount"`
	Currency        string   `json:"currency"`
	PaymentDeadline string   `json:"paymentDeadline"`
	CreatedAt       string   `json:"createdAt"`
}

type Booking struct {
	BookingData
}

type Flight struct {
	ID           string
	FlightNumber string
	Origin       string
	Destination  string
	DepartureAt  time.Time
	BaseAmount   float64
	Currency     string
}

type apiResponse struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

type Store struct {
	mu            sync.Mutex
	bookings      map[string]*Booking
	bookingsByPNR map[string]*Booking
	seatLocks     map[string]map[string]time.Time
	flights       map[string]Flight
}

var defaultAllowedOrigins = []string{
	"http://localhost:5173",
	"http://127.0.0.1:5173",
}

func main() {
	rand.Seed(time.Now().UnixNano())
	store := NewStore()
	mux := http.NewServeMux()
	mux.Handle("/api/v1/bookings", corsMiddleware(http.HandlerFunc(store.createBookingHandler)))
	mux.Handle("/api/v1/bookings/", corsMiddleware(http.HandlerFunc(store.bookingHandler)))

	port := ":8080"
	log.Printf("Booking API is running on %s", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func NewStore() *Store {
	return &Store{
		bookings:      make(map[string]*Booking),
		bookingsByPNR: make(map[string]*Booking),
		seatLocks:     make(map[string]map[string]time.Time),
		flights: map[string]Flight{
			"fl_8a2b3c4d-5e6f-7a8b-9c0d": {
				ID:           "fl_8a2b3c4d-5e6f-7a8b-9c0d",
				FlightNumber: "QL101",
				Origin:       "BKK",
				Destination:  "NRT",
				DepartureAt:  time.Date(2026, 6, 1, 8, 30, 0, 0, time.FixedZone("ICT", 7*3600)),
				BaseAmount:   12500.00,
				Currency:     "THB",
			},
		},
	}
}

func (r BookingRequest) Validate() error {
	if strings.TrimSpace(r.FlightID) == "" {
		return fmt.Errorf("flightId is required")
	}
	if strings.TrimSpace(r.Cabin) == "" {
		return fmt.Errorf("cabin is required")
	}
	if strings.TrimSpace(r.FareClass) == "" {
		return fmt.Errorf("fareClass is required")
	}
	if len(r.Seats) == 0 {
		return fmt.Errorf("at least one seat is required")
	}
	if len(r.Passengers) == 0 {
		return fmt.Errorf("at least one passenger is required")
	}
	if strings.TrimSpace(r.ContactEmail) == "" {
		return fmt.Errorf("contactEmail is required")
	}
	if strings.TrimSpace(r.ContactPhone) == "" {
		return fmt.Errorf("contactPhone is required")
	}
	for i, passenger := range r.Passengers {
		if strings.TrimSpace(passenger.FirstName) == "" {
			return fmt.Errorf("passenger %d firstName is required", i+1)
		}
		if strings.TrimSpace(passenger.LastName) == "" {
			return fmt.Errorf("passenger %d lastName is required", i+1)
		}
		if strings.TrimSpace(passenger.DateOfBirth) == "" {
			return fmt.Errorf("passenger %d dateOfBirth is required", i+1)
		}
		if strings.TrimSpace(passenger.Nationality) == "" {
			return fmt.Errorf("passenger %d nationality is required", i+1)
		}
		if strings.TrimSpace(passenger.PassportNumber) == "" {
			return fmt.Errorf("passenger %d passportNumber is required", i+1)
		}
		if strings.TrimSpace(passenger.PassportExpiry) == "" {
			return fmt.Errorf("passenger %d passportExpiry is required", i+1)
		}
	}
	return nil
}

func (s *Store) createBookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJSON(w, http.StatusMethodNotAllowed, apiResponse{Error: "method not allowed"})
		return
	}

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid JSON payload"})
		return
	}

	if err := req.Validate(); err != nil {
		sendJSON(w, http.StatusBadRequest, apiResponse{Error: err.Error()})
		return
	}

	flight, ok := s.flights[req.FlightID]
	if !ok {
		sendJSON(w, http.StatusNotFound, apiResponse{Error: "flightId not found"})
		return
	}

	if len(req.Seats) != len(req.Passengers) {
		sendJSON(w, http.StatusUnprocessableEntity, apiResponse{Error: "seat count must equal passenger count"})
		return
	}

	if err := s.reserveSeats(req.FlightID, req.Seats); err != nil {
		sendJSON(w, http.StatusConflict, apiResponse{Error: err.Error()})
		return
	}

	now := time.Now().UTC()
	booking := &Booking{
		BookingData: BookingData{
			BookingID:       fmt.Sprintf("bk_%s", randomID()),
			PNR:             s.generateUniquePNR(),
			Status:          "PENDING",
			FlightID:        flight.ID,
			FlightNumber:    flight.FlightNumber,
			Origin:          flight.Origin,
			Destination:     flight.Destination,
			DepartureAt:     flight.DepartureAt.Format(time.RFC3339),
			Cabin:           req.Cabin,
			FareClass:       req.FareClass,
			Seats:           req.Seats,
			TotalAmount:     flight.BaseAmount * float64(len(req.Seats)),
			Currency:        flight.Currency,
			PaymentDeadline: now.Add(15 * time.Minute).UTC().Format(time.RFC3339),
			CreatedAt:       now.Format(time.RFC3339),
		},
	}

	s.mu.Lock()
	s.bookings[booking.BookingID] = booking
	s.bookingsByPNR[booking.PNR] = booking
	s.mu.Unlock()

	sendJSON(w, http.StatusCreated, apiResponse{Data: booking.BookingData})
}

func (s *Store) bookingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendJSON(w, http.StatusMethodNotAllowed, apiResponse{Error: "method not allowed"})
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/v1/bookings/")
	if path == "" || path == "/" {
		sendJSON(w, http.StatusNotFound, apiResponse{Error: "booking not found"})
		return
	}

	if strings.HasPrefix(path, "pnr/") {
		pnr := strings.TrimPrefix(path, "pnr/")
		booking := s.getBookingByPNR(pnr)
		if booking == nil {
			sendJSON(w, http.StatusNotFound, apiResponse{Error: "booking not found"})
			return
		}
		sendJSON(w, http.StatusOK, apiResponse{Data: booking.BookingData})
		return
	}

	booking := s.getBookingByID(path)
	if booking == nil {
		sendJSON(w, http.StatusNotFound, apiResponse{Error: "booking not found"})
		return
	}

	sendJSON(w, http.StatusOK, apiResponse{Data: booking.BookingData})
}

func (s *Store) getBookingByID(id string) *Booking {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bookings[id]
}

func (s *Store) getBookingByPNR(pnr string) *Booking {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.bookingsByPNR[strings.ToUpper(pnr)]
}

func (s *Store) reserveSeats(flightID string, seats []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanupExpiredLocksLocked()

	if s.seatLocks[flightID] == nil {
		s.seatLocks[flightID] = make(map[string]time.Time)
	}

	now := time.Now()
	for _, seat := range seats {
		seat = strings.ToUpper(strings.TrimSpace(seat))
		if expiry, locked := s.seatLocks[flightID][seat]; locked && now.Before(expiry) {
			return fmt.Errorf("seat %s already locked by another session", seat)
		}
	}

	expiry := now.Add(15 * time.Minute)
	for _, seat := range seats {
		seat = strings.ToUpper(strings.TrimSpace(seat))
		s.seatLocks[flightID][seat] = expiry
	}

	return nil
}

func (s *Store) cleanupExpiredLocksLocked() {
	now := time.Now()
	for flightID, seats := range s.seatLocks {
		for seat, expiry := range seats {
			if now.After(expiry) {
				delete(seats, seat)
			}
		}
		if len(seats) == 0 {
			delete(s.seatLocks, flightID)
		}
	}
}

func (s *Store) generateUniquePNR() string {
	for {
		pnr := randomPNR()
		s.mu.Lock()
		_, exists := s.bookingsByPNR[pnr]
		s.mu.Unlock()
		if !exists {
			return pnr
		}
	}
}

func randomPNR() string {
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, 6)
	for i := 0; i < 6; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func randomID() string {
	chars := "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, 12)
	for i := 0; i < len(result); i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func sendJSON(w http.ResponseWriter, status int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" && isAllowedOrigin(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func isAllowedOrigin(origin string) bool {
	allowed := defaultAllowedOrigins
	if env := strings.TrimSpace(getEnv("ALLOWED_ORIGIN", "")); env != "" {
		allowed = append(allowed, env)
	}
	for _, candidate := range allowed {
		if strings.EqualFold(candidate, origin) {
			return true
		}
	}
	return false
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return strings.TrimSpace(value)
}
