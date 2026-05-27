package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testResponse[T any] struct {
	Data  T      `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func TestCreateBooking_Success(t *testing.T) {
	store := NewStore()
	request := BookingRequest{
		FlightID:   "fl_8a2b3c4d-5e6f-7a8b-9c0d",
		Cabin:      "economy",
		FareClass:  "M",
		Seats:      []string{"14A"},
		Passengers: []Passenger{{FirstName: "Somchai", LastName: "Jaidee", DateOfBirth: "1990-05-15", Nationality: "THA", PassportNumber: "AA1234567", PassportExpiry: "2030-12-31"}},
		ContactEmail: "somchai@example.com",
		ContactPhone: "+66812345678",
	}
	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	rw := httptest.NewRecorder()

	store.createBookingHandler(rw, req)

	if rw.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rw.Code)
	}

	var resp testResponse[BookingData]
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Error != "" {
		t.Fatalf("unexpected error: %s", resp.Error)
	}
	if resp.Data.BookingID == "" || resp.Data.PNR == "" {
		t.Fatal("expected bookingId and pnr to be populated")
	}
}

func TestCreateBooking_SeatConflict(t *testing.T) {
	store := NewStore()
	reqBody := BookingRequest{
		FlightID:   "fl_8a2b3c4d-5e6f-7a8b-9c0d",
		Cabin:      "economy",
		FareClass:  "M",
		Seats:      []string{"14A"},
		Passengers: []Passenger{{FirstName: "Somchai", LastName: "Jaidee", DateOfBirth: "1990-05-15", Nationality: "THA", PassportNumber: "AA1234567", PassportExpiry: "2030-12-31"}},
		ContactEmail: "somchai@example.com",
		ContactPhone: "+66812345678",
	}
	body, _ := json.Marshal(reqBody)
	req1 := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	rec1 := httptest.NewRecorder()
	store.createBookingHandler(rec1, req1)
	if rec1.Code != http.StatusCreated {
		t.Fatalf("first request failed: %d", rec1.Code)
	}

	req2 := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	rec2 := httptest.NewRecorder()
	store.createBookingHandler(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Fatalf("expected conflict status, got %d", rec2.Code)
	}
}

func TestCreateBooking_ValidationError(t *testing.T) {
	store := NewStore()
	reqBody := BookingRequest{}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/bookings", bytes.NewBuffer(body))
	rec := httptest.NewRecorder()

	store.createBookingHandler(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request status, got %d", rec.Code)
	}
}
