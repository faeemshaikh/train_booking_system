package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/faeemshaikh/train_booking_system/proto"
	"github.com/google/uuid" // We are generating unique Id for ticket in server logic since persistance layer is not yet added
	"google.golang.org/grpc"
)

// Struct definition to store user and ticket information
type User struct {
	TicketID  string
	FirstName string
	LastName  string
	Email     string
	From      string
	To        string
	Price     float64
	Section   string
	Seat      int32
}

// Struct to handle seat limits and track occupied seats
type Section struct {
	CurrentSeat int32          // The next available seat in the section
	MaxSeats    int32          // Maximum number of seats in the section
	Occupied    map[int32]bool // Map to track occupied seats
}

// Struct to handle train sections dynamically and store users by ticket ID
type server struct {
	pb.UnimplementedTrainServiceServer
	mu       sync.Mutex
	users    map[string]*User    // Keyed by TicketID for unique identification
	sections map[string]*Section // Store section with available seat counter
}

// NewServer initializes the server with sections and max capacity per section
func NewServer() *server {
	return &server{
		users: make(map[string]*User),
		sections: map[string]*Section{
			"A": {CurrentSeat: 1, MaxSeats: 10, Occupied: make(map[int32]bool)}, // Section A with 10 seats
			"B": {CurrentSeat: 1, MaxSeats: 10, Occupied: make(map[int32]bool)}, // Section B with 10 seats
		},
	}
}

// PurchaseTicket allows a user to purchase a ticket and assigns a seat and section
func (s *server) PurchaseTicket(ctx context.Context, req *pb.PurchaseRequest) (*pb.PurchaseResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ticketID := uuid.New().String()

	// Try to find an available section and seat
	var sectionName string
	var seatNumber int32
	for sec, section := range s.sections {
		// Find the next available seat that is unoccupied
		for seat := section.CurrentSeat; seat <= section.MaxSeats; seat++ {
			if !section.Occupied[seat] { // Seat is unoccupied
				sectionName = sec
				seatNumber = seat
				section.Occupied[seat] = true // Mark seat as occupied
				section.CurrentSeat = seat + 1
				break
			}
		}
		if sectionName != "" {
			break
		}
	}

	// If no section with available seats, return an error message
	if sectionName == "" {
		return &pb.PurchaseResponse{Message: "No seats available in any section"}, fmt.Errorf("no seats available")
	}

	// Assign the user ticket information
	user := &User{
		TicketID:  ticketID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		From:      req.From,
		To:        req.To,
		Price:     20.0,
		Section:   sectionName,
		Seat:      seatNumber,
	}
	s.users[ticketID] = user
	log.Printf("ticket %s purchased for user %s %s in section %s, seat %d", ticketID, user.FirstName, user.LastName, sectionName, seatNumber)

	return &pb.PurchaseResponse{
		Message:  fmt.Sprintf("Ticket purchased successfully in section %s, seat %d!", sectionName, seatNumber),
		TicketId: ticketID,
	}, nil
}

// GetReceipt provides the receipt for a specific ticket ID
func (s *server) GetReceipt(ctx context.Context, req *pb.ReceiptRequest) (*pb.ReceiptResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, exists := s.users[req.TicketId]; exists {
		log.Printf("retrieving receipt for ticket %s", req.TicketId)
		return &pb.ReceiptResponse{
			TicketId: user.TicketID,
			From:     user.From,
			To:       user.To,
			User:     fmt.Sprintf("%s %s", user.FirstName, user.LastName),
			Price:    user.Price,
			Section:  user.Section,
			Seat:     user.Seat,
		}, nil
	}
	return nil, fmt.Errorf("ticket ID not found")
}

// ViewSeats shows all users and seats in a specified section
func (s *server) ViewSeats(ctx context.Context, req *pb.ViewSeatRequest) (*pb.ViewSeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var seats []*pb.SeatAllocation
	for _, user := range s.users {
		if user.Section == req.Section {
			seats = append(seats, &pb.SeatAllocation{
				User: fmt.Sprintf("%s %s", user.FirstName, user.LastName),
				Seat: user.Seat,
			})
		}
	}
	log.Printf("viewing seats in section %s", req.Section)
	return &pb.ViewSeatResponse{Seats: seats}, nil
}

// RemoveUser deletes a user based on ticket ID
func (s *server) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, exists := s.users[req.TicketId]; exists {
		delete(s.users, req.TicketId)
		// Free up the seat in the respective section
		s.sections[user.Section].Occupied[user.Seat] = false
		log.Printf("removed user with ticket ID %s", req.TicketId)
		return &pb.RemoveUserResponse{Message: "User removed successfully"}, nil
	}
	return &pb.RemoveUserResponse{Message: "Ticket ID not found"}, nil
}

// ModifySeat allows modification of the seat assignment for a user by ticket ID
func (s *server) ModifySeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.ModifySeatResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if user, exists := s.users[req.TicketId]; exists {
		section, ok := s.sections[req.Section]
		if !ok || req.Seat > section.MaxSeats || section.Occupied[req.Seat] {
			return &pb.ModifySeatResponse{Message: "Invalid or already occupied seat"}, nil
		}

		// Release the previous seat
		s.sections[user.Section].Occupied[user.Seat] = false

		// Update to the new seat and section
		previousSection := user.Section
		previousSeat := user.Seat
		user.Section = req.Section
		user.Seat = req.Seat
		section.Occupied[req.Seat] = true // Mark new seat as occupied

		log.Printf("modified seat for ticket %s: from section %s, seat %d to section %s, seat %d", req.TicketId, previousSection, previousSeat, user.Section, user.Seat)

		return &pb.ModifySeatResponse{
			Message:         "Seat modified successfully",
			PreviousSection: previousSection,
			PreviousSeat:    previousSeat,
			UpdatedSection:  user.Section,
			UpdatedSeat:     user.Seat,
		}, nil
	}
	return &pb.ModifySeatResponse{Message: "Ticket ID not found"}, nil
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterTrainServiceServer(grpcServer, NewServer())
	log.Println("server is running on port :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
