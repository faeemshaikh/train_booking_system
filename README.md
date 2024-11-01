# Train Ticket Booking System

This is a gRPC-based train booking system implemented in Golang. The system allows users to purchase train tickets, view receipts, allocate seats, and manage users and seats in real-time. All data is stored in memory, making it lightweight and suitable for local testing.

There are two sections, **"A"** and **"B"**, in the train, each with a maximum capacity of 10 seats. The ticket price is hard-coded to $20. These values should ideally be fetched from a persistence layer, but for this assignment, they are hard-coded.

## Features

- **Ticket Purchase**: Users can purchase a ticket with their details and are automatically allocated a seat in either section A or B. The ticket ID is randomly generated.
- **Receipt Viewing**: Users can view their ticket details, including route, section, seat, and price.
- **Seat Management**: View all seat allocations within a section, modify seat allocations, and remove users.
- **In-memory Storage**: All data is stored in memory, no persistence layer has been added yet.

## API Overview

1. **PurchaseTicket**
   - **Description**: Purchases a ticket, assigning the user to a seat in either section A or B, based on availability.
   - **Parameters**: User details (first name, last name, email), route information (from, to, price, etc.).
   - **Response**: Confirmation message, ticket ID, assigned section, and seat number.

2. **GetReceipt**
   - **Description**: Retrieves ticket details for a given ticket ID. (Ticket ID is randomly generated when purchase happens.)
   - **Parameters**: `ticket_id`.
   - **Response**: Ticket details, including route, price, section, and seat.

3. **ViewSeats**
   - **Description**: Shows seat allocations for a specified section (either A or B).
   - **Parameters**: Section identifier.
   - **Response**: List of seats and assigned users in the specified section.

4. **RemoveUser**
   - **Description**: Removes a user based on their ticket ID.
   - **Parameters**: `ticket_id`.
   - **Response**: Confirmation message indicating successful removal.

5. **ModifySeat**
   - **Description**: Modifies the seat assignment for an existing ticket only if the requested seat is available.
   - **Parameters**: `ticket_id`, `section`, `seat`.
   - **Response**: Confirmation message, previous seat details, and updated seat details.

## Project Structure

```plaintext
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
    ├── server.go               # gRPC server with train booking logic
    └── server_test.go          # Unit tests for server functionality
```

## Prerequisites
1. Go
2. protoc
3. gRPC and Protobuf Go plugins


## Setup and Running
1. Generate gRPC Code
    If any changes are made to proto/train.proto, regenerate the Go code with:
    ```
    protoc --go_out=. --go-grpc_out=. proto/train.proto
    ```
2. Run the Server
    ```
    go run server/server.go
    ```
3. Run the Client
    ```
    go run client/client.go
    ```

## Example Commands in client.go
Add commands in client/client.go to test various APIs. Examples:
1. Purchase Ticket: Buy a ticket and allocate a seat.
2. View Receipt: Fetch receipt details by ticket_id.
3. View Seats: Check seat allocations in sections A and B.
4. Remove User: Remove a user from the train by ticket_id.
5. Modify Seat: Change a user’s seat allocation.

## Testing
Running Tests
```
go test ./server
```

## Test Coverage
The latest code has a test coverage of approximately 82%.
