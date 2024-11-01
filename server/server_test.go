package main

import (
	"context"
	"fmt"
	"strings"
	"testing"

	pb "github.com/faeemshaikh/train_booking_system/proto"
)

// Test for PurchaseTicket API to verify ticket purchasing and seat allocation
func TestPurchaseTicket(t *testing.T) {
	srv := NewServer()
	req := &pb.PurchaseRequest{
		FirstName: "Ramesh",
		LastName:  "Rao",
		Email:     "ramesh@example.com",
		From:      "CityA",
		To:        "CityB",
	}

	res, err := srv.PurchaseTicket(context.Background(), req)
	expectedMessage := "Ticket purchased successfully"
	if err != nil || !strings.Contains(res.Message, expectedMessage) {
		t.Errorf("Expected success message '%s', got '%v', error: %v", expectedMessage, res.Message, err)
	}
}

// Test for GetReceipt API to confirm receipt retrieval functionality
func TestGetReceipt(t *testing.T) {
	srv := NewServer()
	// First, purchase a ticket to ensure a user exists
	ticketReq := &pb.PurchaseRequest{
		FirstName: "Fname",
		LastName:  "Sname",
		Email:     "fname@example.com",
		From:      "Mumbai",
		To:        "Delhi",
	}
	purchaseRes, _ := srv.PurchaseTicket(context.Background(), ticketReq)

	// Now, test the GetReceipt API using the generated ticket ID
	req := &pb.ReceiptRequest{TicketId: purchaseRes.TicketId}
	res, err := srv.GetReceipt(context.Background(), req)
	if err != nil || res.User != "Fname Sname" {
		t.Errorf("Expected user Fname Sname in receipt, got '%v', error: %v", res.User, err)
	}
}

// Test for ViewSeats API to ensure correct viewing of allocated seats in a section
func TestViewSeats(t *testing.T) {
	srv := NewServer()
	srv.PurchaseTicket(context.Background(), &pb.PurchaseRequest{
		FirstName: "Bob",
		LastName:  "Brown",
		Email:     "bb@example.com",
		From:      "Patna",
		To:        "Gaya",
	})

	// Test for seats in section A
	req := &pb.ViewSeatRequest{Section: "A"}
	res, err := srv.ViewSeats(context.Background(), req)
	if err != nil || len(res.Seats) == 0 {
		t.Errorf("Expected seats in section A, got %v, error: %v", res.Seats, err)
	}
}

// Test for RemoveUser API to validate user removal based on Ticket ID
func TestRemoveUser(t *testing.T) {
	srv := NewServer()
	ticketReq := &pb.PurchaseRequest{
		FirstName: "Chetan",
		LastName:  "M",
		Email:     "mchetan@yahoo.com",
		From:      "Pune",
		To:        "Indore",
	}
	purchaseRes, _ := srv.PurchaseTicket(context.Background(), ticketReq)

	req := &pb.RemoveUserRequest{TicketId: purchaseRes.TicketId}
	res, err := srv.RemoveUser(context.Background(), req)
	if err != nil || res.Message != "User removed successfully" {
		t.Errorf("Expected success message 'User removed successfully', got '%v', error: %v", res.Message, err)
	}
}

// Test for ModifySeat API to validate seat modification functionality
func TestModifySeat(t *testing.T) {
	srv := NewServer()
	ticketReq := &pb.PurchaseRequest{
		FirstName: "Dev",
		LastName:  "Anand",
		Email:     "d.anand@gmail.com",
		From:      "Jaipur",
		To:        "Chennai",
	}
	purchaseRes, _ := srv.PurchaseTicket(context.Background(), ticketReq)

	// Capture the initial section and seat from the purchase response
	initialSection := srv.users[purchaseRes.TicketId].Section
	initialSeat := srv.users[purchaseRes.TicketId].Seat

	// Modify the seat to section B, seat 5
	req := &pb.ModifySeatRequest{
		TicketId: purchaseRes.TicketId,
		Section:  "B",
		Seat:     5,
	}
	res, err := srv.ModifySeat(context.Background(), req)

	expectedModifyMessage := "Seat modified successfully"
	if err != nil || res.Message != expectedModifyMessage || res.PreviousSection != initialSection || res.PreviousSeat != initialSeat || res.UpdatedSection != "B" || res.UpdatedSeat != 5 {
		t.Errorf("Expected modification message '%s' with previous section '%v' and seat %d, and updated section 'B' and seat 5. Got message '%v', PreviousSection '%v', PreviousSeat '%v', UpdatedSection '%v', UpdatedSeat '%v'",
			expectedModifyMessage, initialSection, initialSeat, res.Message, res.PreviousSection, res.PreviousSeat, res.UpdatedSection, res.UpdatedSeat)
	}
}

// Test for handling full capacity in PurchaseTicket API
func TestPurchaseTicketFullCapacity(t *testing.T) {
	srv := NewServer()
	// Fill section A and B to capacity
	for i := 0; i < 20; i++ { // Total seats = 10 in section A + 10 in section B
		srv.PurchaseTicket(context.Background(), &pb.PurchaseRequest{
			FirstName: fmt.Sprintf("User%d", i),
			LastName:  "LastName",
			Email:     fmt.Sprintf("user%d@example.com", i),
			From:      "CityA",
			To:        "CityB",
		})
	}

	// Attempt to purchase another ticket after full capacity
	req := &pb.PurchaseRequest{
		FirstName: "Overflow",
		LastName:  "User",
		Email:     "overflow.user@example.com",
		From:      "CityC",
		To:        "CityD",
	}
	res, _ := srv.PurchaseTicket(context.Background(), req)
	expectedErrorMessage := "No seats available in any section"
	if res.Message != expectedErrorMessage {
		t.Errorf("Expected full capacity message '%s', got '%v'", expectedErrorMessage, res.Message)
	}
}
