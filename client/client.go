package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/faeemshaikh/train_booking_system/proto"
	"google.golang.org/grpc"
)

// This is sample client to demonstrate Purchase, Modify ticket, view receipt, and view seats
func main() {
	// Initialize gRPC client connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewTrainServiceClient(conn)

	// Example: Purchase Ticket
	purchaseResp, err := client.PurchaseTicket(context.Background(), &pb.PurchaseRequest{
		FirstName: "Faeem",
		LastName:  "Shaikh",
		Email:     "faeemshaikh123@gmail.com",
		From:      "London",
		To:        "France",
	})
	if err != nil {
		log.Fatalf("error purchasing ticket: %v", err)
	}
	fmt.Printf("Purchase Response: %s, Ticket ID: %s\n", purchaseResp.Message, purchaseResp.TicketId)

	// Example: Get Receipt
	receiptResp, err := client.GetReceipt(context.Background(), &pb.ReceiptRequest{
		TicketId: purchaseResp.TicketId,
	})
	if err != nil {
		log.Fatalf("error getting receipt: %v", err)
	}
	fmt.Printf("Receipt: User: %s, From: %s, To: %s, Price: %.2f, Section: %s, Seat: %d\n",
		receiptResp.User, receiptResp.From, receiptResp.To, receiptResp.Price, receiptResp.Section, receiptResp.Seat)

	// Example: Modify Seat
	modifyResp, err := client.ModifySeat(context.Background(), &pb.ModifySeatRequest{
		TicketId: purchaseResp.TicketId,
		Section:  "B",
		Seat:     5,
	})
	if err != nil {
		log.Fatalf("error modifying seat: %v", err)
	}
	fmt.Printf("Modify Seat Response: %s\nPrevious Section: %s, Previous Seat: %d\nUpdated Section: %s, Updated Seat: %d\n",
		modifyResp.Message, modifyResp.PreviousSection, modifyResp.PreviousSeat, modifyResp.UpdatedSection, modifyResp.UpdatedSeat)

	// View seats in both sections
	viewSeats(client, "A")
	viewSeats(client, "B")
}

// viewSeats retrieves and displays seat allocations for a specified section.
func viewSeats(client pb.TrainServiceClient, section string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.ViewSeatRequest{Section: section}
	res, err := client.ViewSeats(ctx, req)
	if err != nil {
		log.Fatalf("error viewing seats in section %s: %v", section, err)
	}

	if len(res.Seats) == 0 {
		fmt.Printf("Seats in section %s:\n", section)
		fmt.Printf("\tNone\n")
	} else {
		fmt.Printf("Seats in section %s:\n", section)
		for _, seatAllocation := range res.Seats {
			fmt.Printf("User: %s, Seat: %d\n", seatAllocation.User, seatAllocation.Seat)
		}
	}
}
