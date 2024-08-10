# Personal Finance Management API

This is a simple personal finance management API built with Go. The API allows users to manage their finances by creating, updating, and deleting transactions and categories. User authentication is also implemented with password hashing.

## Features

- **User Management**: Users can be created with secure password hashing, and user details can be updated or deleted.
- **Category Management**: Users can create, update, list, and delete categories for their transactions.
- **Transaction Management**: Users can create, update, list, and delete financial transactions associated with specific categories.

## Technologies Used

- **Go**: The core language used for developing the API.
- **Gorilla Mux**: Router for handling HTTP requests.
- **MySQL**: Database for storing users, categories, and transactions.
- **bcrypt**: For securely hashing user passwords.

## Endpoints

### User Endpoints

- **POST /users**: Create a new user with password confirmation and secure password hashing.
- **GET /users**: List all users.
- **GET /users/{id}**: Get details of a specific user by ID.
- **PUT /users/{id}**: Update user details.
- **DELETE /users/{id}**: Delete a user by ID.

### Category Endpoints

- **POST /categories**: Create a new category.
- **GET /categories**: List all categories.
- **GET /categories/{id}**: Get details of a specific category by ID.
- **PUT /categories/{id}**: Update category details.
- **DELETE /categories/{id}**: Delete a category by ID.

### Transaction Endpoints

- **POST /transactions**: Create a new transaction.
- **GET /transactions**: List all transactions.
- **GET /transactions/{id}**: Get details of a specific transaction by ID.
- **PUT /transactions/{id}**: Update transaction details.
- **DELETE /transactions/{id}**: Delete a transaction by ID.

## Setup

### Prerequisites

- Go 1.16 or later
- MySQL

### Installation

1. **Clone the repository**:
    ```bash
    git clone https://github.com/halilomergurgan/go-financeApp.git
    cd go-financeApp
    ```

2. **Install dependencies**:
    ```bash
    go mod tidy
    ```

3. **Set up MySQL**:
    - Create a database named `financeApp`.
    - Update the database credentials in `main.go` (DSN string).

4. **Run the application**:
    ```bash
    go run main.go
    ```

5. **Test with Postman**:
    - Use the provided endpoints to interact with the API.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Feel free to submit issues and pull requests.
