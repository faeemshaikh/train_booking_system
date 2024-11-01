Train Booking System
This is a gRPC-based train booking system implemented in Golang. The system allows users to purchase train tickets, view receipts, allocate seats, and manage users and seats in real-time. All data is stored in memory, making it lightweight and suitable for local testing. 

There are two sections "A" and "B" in the train with max capacity of 10 each. The ticket price is hard coded to $20. These values should be fetched from persitance layer, but as a part of assignment I have hard coded them.

Features
Ticket Purchase: Users can purchase a ticket with their details and are automatically allocated a seat in either section A or B. The ticket Id is randomly generated
Receipt Viewing: Users can view their ticket details, including route, section, seat, and price.
Seat Management: View all seat allocations within a section, modify seat allocations, and remove users.
In-memory Storage: All data is stored in memory and no persistance layer added yet.

API Overview
1. PurchaseTicket
Purchases a ticket, assigning the user to a seat in either section A or B, based on availability.
Parameters: User details (first name, last name, email), route information (from, to, price etc.).
Response: Confirmation message, ticket ID, assigned section, and seat number.
2. GetReceipt
Retrieves ticket details for a given ticket ID. (ticket id is randomaly generated when purchase happens)
Parameters: ticket_id.
Response: Ticket details, including route, price, section, and seat.
3. ViewSeats
Shows seat allocations for a specified section (either A or B).
Parameters: Section identifier.
Response: List of seats and assigned users in the specified section.
4. RemoveUser
Removes a user based on their ticket ID.
Parameters: ticket_id.
Response: Confirmation message indicating successful removal.
5. ModifySeat
Modifies the seat assignment for an existing ticket only if requested seat is available
Parameters: ticket_id, section, seat.
Response: Confirmation message, previous seat details, and updated seat details.


Project Structure
.
├── README.md                   # Project description and setup
├── client/
│   └── client.go               # gRPC client for testing functionality
├── go.mod                      # Go module dependencies
├── go.sum                      # Go module dependencies checksums
├── proto/
│   ├── train.proto             # Protobuf definitions
│   ├── train.pb.go             # Generated Protobuf Go code
│   └── train_grpc.pb.go        # Generated gRPC Go code
└── server/
    └── server.go               # gRPC server with train booking logic
    └── server_test.go          # Unit tests for server functionality

Prerequisites
Go
protoc
gRPC and Protobuf Go plugins

Setup and Running
1. If any changes are made to proto/train.proto, regenerate the Go code:
protoc --go_out=. --go-grpc_out=. proto/train.proto
2. Run the Server
Start the gRPC server to listen on port 50051: go run server/server.go
3. Run the Client
go run client/client.go


Example Commands in client.go
Add commands in client/client.go to test various APIs. Examples:

Purchase Ticket: Buy a ticket and allocate a seat.
View Receipt: Fetch receipt details by ticket_id.
View Seats: Check seat allocations in sections A and B.
Remove User: Remove a user from the train by ticket_id.
Modify Seat: Change a user’s seat allocation.

Testing
Running Tests
The server_test.go file contains tests for all primary functionality, including ticket purchase, seat viewing, seat modification, and user removal. Run the tests with:
go test ./server

Test Coverage : 82%
