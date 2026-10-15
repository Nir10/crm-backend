# CRMBackEnd

CRMBackEnd is a simple REST API built with Go for managing customer records. It provides endpoints to create, read, update, and delete customers using JSON over HTTP.

## Features

- Get all customers
- Get a single customer by ID
- Add a new customer
- Update an existing customer
- Delete a customer
- Simple landing page with available routes

## Tech Stack

- Go 1.25+
- Gorilla Mux router

## Project Structure 
```
CRMBackEnd/ 
├── intro 
├── go.mod 
├── main.go 
└── README.md
```

## Installation

### Prerequisites

Make sure you have Go installed.

You can verify your Go installation with:
```
bash go version
```

### Clone the Repository
```
bash git clone <repository-url> cd CRMBackEnd
```
### Install Dependencies
```
bash go mod download

``` 
## Launch

Run the application with:
```

bash go run main.go
```

The server will start on:
```
text http://localhost:3000
``` 

You should see the following message in the terminal:
```
Server is running on port 3000```
```

## Usage

The API works with customer data in JSON format.

### Customer Object
```
json { "id": 1, "name": "John", "role": "Customer", "email": "john@abc.com", "phone": "9234562341", "contacted": false } 
```

## API Endpoints

### Home Page
```

http GET /
```

Displays a simple page listing the available customer routes.

### Get All Customers
```
http GET /customers
``` 

Example:
```
bash curl http://localhost:3000/customers
```

### Get a Customer by ID
```
http GET /customers/{id}
``` 

Example:
```

bash curl http://localhost:3000/customers/1
```

### Create a Customer
```
http POST /customers
``` 

Example:
```
bash curl -X POST http://localhost:3000/customers
-H "Content-Type: application/json"
-d '{ "id": 4, "name": "Sam", "role": "Customer", "email": "sam@example.com", "phone": "9234562341", "contacted": false }'
```

### Update a Customer
```
http PUT /customers/{id}
```

Example:
```
bash curl -X PUT http://localhost:3000/customers/1
-H "Content-Type: application/json"
-d '{ "id": 1, "name": "John Smith", "role": "Customer", "email": "john.smith@example.com", "phone": "9234562341", "contacted": true }'
```

### Delete a Customer
```
http DELETE /customers/{id} 
```
Example:
```
bash curl -X DELETE http://localhost:3000/customers/1
```

## Notes

- Customer data is stored in memory.
- Data will reset when the server restarts.
- The API runs locally on port `3000` by default.

