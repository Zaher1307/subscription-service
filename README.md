# Subscription Service API

A RESTful API for managing product subscriptions, built with Go, Gin, and PostgreSQL.

## Features

- User management
- Product catalog
- Subscription management
- Monthly billing system
- Bill payment handling

## Requirements

- Go 1.18+
- PostgreSQL 12+
- Docker
- Docker Compose

## Setup

### Running the Application

```sh
$ git clone https://github.com/Zaher1307/subscription-service.git
$ cd subscription-service
$ docker-compose up --build
```

The service will be available at http://localhost:8080.

## Testing

Run the tests with:

```bash
go test ./tests/...
```

## API Documentation

### Base URL

All API endpoints are relative to: `http://localhost:8080/api`

### Authentication

This version does not include authentication.

### Endpoints

#### Health Check

```
GET /health
```

Returns the health status of the service.

**Response Example:**

```json
{
  "status": "ok",
  "timestamp": "2025-03-10T12:00:00Z",
  "components": {
    "database": "ok"
  }
}
```

#### Users

##### Create User

```
POST /api/users
```

Create a new user account.

**Request Body:**

```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Response:**

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-03-10T12:00:00Z"
}
```

##### Get User

```
GET /api/users/:id
```

Retrieve a user by ID.

**Response:**

```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-03-10T12:00:00Z"
}
```

#### Products

##### List Products

```
GET /api/products
```

Get a list of all available products.

**Response:**

```json
[
  {
    "id": 1,
    "name": "Premium Coffee Subscription",
    "description": "Artisanal coffee delivered monthly",
    "price": 19.99,
    "created_at": "2025-03-10T12:00:00Z"
  },
  {
    "id": 2,
    "name": "Standard Coffee Subscription",
    "description": "Great quality coffee delivered monthly",
    "price": 14.99,
    "created_at": "2025-03-10T12:00:00Z"
  }
]
```

##### Get Product

```
GET /api/products/:id
```

Get details for a specific product.

**Response:**

```json
{
  "id": 1,
  "name": "Premium Coffee Subscription",
  "description": "Artisanal coffee delivered monthly",
  "price": 19.99,
  "created_at": "2025-03-10T12:00:00Z"
}
```

#### Subscriptions

##### Create Subscription

```
POST /api/subscriptions
```

Subscribe a user to a product.

**Request Body:**

```json
{
  "user_id": 1,
  "product_id": 1
}
```

**Response:**

```json
{
  "subscription": {
    "id": 1,
    "user_id": 1,
    "product_id": 1,
    "start_date": "2025-03-10T12:00:00Z",
    "next_billing_date": "2025-04-10T12:00:00Z",
    "status": "active",
    "created_at": "2025-03-10T12:00:00Z"
  },
  "initial_bill": {
    "id": 1,
    "subscription_id": 1,
    "amount": 19.99,
    "due_date": "2025-03-10T12:00:00Z",
    "status": "pending",
    "created_at": "2025-03-10T12:00:00Z"
  }
}
```

##### Get Subscription

```
GET /api/subscriptions/:id
```

Get details for a specific subscription.

**Response:**

```json
{
  "id": 1,
  "user_id": 1,
  "product_id": 1,
  "start_date": "2025-03-10T12:00:00Z",
  "next_billing_date": "2025-04-10T12:00:00Z",
  "status": "active",
  "created_at": "2025-03-10T12:00:00Z"
}
```

#### Bills

##### Get User Bills

```
GET /api/bills/user/:user_id
```

Get all bills for a specific user.

**Response:**

```json
[
  {
    "id": 1,
    "subscription_id": 1,
    "amount": 19.99,
    "due_date": "2025-03-10T12:00:00Z",
    "status": "pending",
    "created_at": "2025-03-10T12:00:00Z"
  }
]
```

##### Get Bill

```
GET /api/bills/:id
```

Get details for a specific bill.

**Response:**

```json
{
  "id": 1,
  "subscription_id": 1,
  "amount": 19.99,
  "due_date": "2025-03-10T12:00:00Z",
  "status": "pending",
  "created_at": "2025-03-10T12:00:00Z"
}
```

##### Pay Bill

```
POST /api/bills/:id/pay
```

Mark a bill as paid.

**Response:**

```json
{
  "message": "Bill paid successfully"
}
```

### Error Responses

All endpoints return appropriate HTTP status codes:

- `200 OK` - Request succeeded
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request parameters
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

Error response format:

```json
{
  "error": "Error message description"
}
```

## Diagrams

### Sequence diagram for paying a bill

![Mermaid Diagram](https://mermaid.ink/svg/pako:eJy9lFFPwjAQx79K0ydMQIjxqQkk4IzwgBKHb3sp7QUat262N3QhfHc7N92ESnhQ97Lt7n__-_V22Y6KVAJl1MJLDlpAoPja8CTSxF0ZN6iEyrhGchMr0Hgcn6g4nnItYzD-pNLrEMxWCfDnHyFLrcLUFMf5MF9ZYVSGKtVtXaWskHqjUYuBkcVDuCT9lQvZPlOyn_HauKWqaxo0V8aLMtRR8qKRN4K6ooFg5A5wUsyC7xWNoOdpUlJVWh7jx9tliBxzS4ZDEtGMKxnRSuBBODzprTGpu0W0NCLKOlMDXBbE5_N5cmdSzY2R6XK5INeDQSWE2MKPrY9PP-fmeWwXrlMzgPOGACX2iU7-j87IUyY5gpuXwcA9dKrxtcSzoEt0-tqC8Vv9PtQ9vGFdcgLNieapxs2_Ao4Fqm05t1bex_cHUIf7GuZCgLXn7eXV115qSbs0AZO4XaNsR3EDSfnLkm4D6X7_DjSdogI)

### Sequence diagram for generating bills

#### NOTE: The Cronjob starts everyday at midnight, if you need to test it just change the matching regex to every minute for example.

![Mermaid Diagram](https://mermaid.ink/svg/pako:eJyVk01uwyAQha-CWNlS0wOwyCJ1f9JVlGy7wTBNUG2GDrhVFOXuxcVSjWxZKUvee8MHM1y4Qg1ccA-fHVgFlZFHku2bZXE5ScEo46QN7IHQvmI9FTamaYw9HoC-jIKpfuhqr8i4YNDuwaE3Aek89e0IdafCkqU_aqwnx0C2Wq9zFMGewQLJAP2-L8pkz00xNQ_Yp0PVwRPSkCgsfg815iOrGQQ_cvoUbhAde0diINUpMyT9X5Av2OixVozr3W-r8q_m7dBAhLQAM2nV72NtztsqP37wZRST7ByAS6YFhHwURBwCiI0u6rhd5rGbrwpWz05ITA0jJuIbKgXes9i9xz7J73gL1EqjubjwcIK2_0xa0ge_Xn8AqSgxJg)

### Sequence diagram for creating a subscription

![Mermaid Diagram](https://mermaid.ink/svg/pako:eJydVE1vwjAM_SteTiCBxibtUmmVgE6D09BgNy6h8SBam3ZJilYh_vtaUtaUfoCWS6T4-fn52e2B-BFD4hCF3wkKHz1Ot5KGawHZianU3OcxFRqmAUeh6-_LZKN8yWPNIzGjggUou0FLlHvuYzfoHeNIcR3JtI77UCi74gsZscTXXZAJDwI7bhCmxaHrNvTkwOJtuYJ7ZYWUSWtAX3AULTswlUg12qFeknUz9wYQG9Vzr19nLfIvWMsGHHhFPfY13-Mkze0ZC1a4cBN_yTRsU44_XGkutlnQ5NNA249w9wyCBybWpr_N2hcpo-xak1xsxpy5xFLYUQVUAD31Bbbx8BlJ0Duuzl2tSXPh8ziywma4DsxWqwU8jUYmAQOF3Zpdt7pvJ68n6dwrnO2X6VVgq5V53rWatR0uy14OMj81eGvxs2FX6rftmVngnj2MfjPXDVpslmuCql_sn5BN9mwJqKJaC-dZ_9tUW_Kgg6dt8R5HD8XiCUYGJEQZUs6IcyB6h2H-I2ZUfpHj8RflVPz3)

### Entity Relationship Diagram for Subscription Service

![Mermaid Diagram](https://mermaid.ink/svg/pako:eJylU8FuwyAM_RXEef2BnKdddtkHVIpcsFKrQCJjpElp_n0gsmZjmbaqHBC8Z8x7xszajBZ1p5GfCQYGfwwqjxSRo5rrpgwKosiqt9cNisIUBhXA4w8QPZDbUCGPUcBPyjCCoO1BKrscQ11MPNpk5JFLLUbDNAmNYeMsGvLgcn4yeJeimE63fH_JKlCpWZ_xlwZfnbXUpiDPLL3NKvbYgO_Sn8i5bLEJWn3nOEnxLm8l3388fa1Ba-CzsuDHFOQBVd_ZCWhHcG3I6_VwGOfmZTp1htg00S-BFIxLuU32Hng9UivTqQEDcpYZ9ZP2yLmfre5mLWf05b9Y4Itelg-42Adp)
