# PK Trade

[![GoDoc](https://godoc.org/github.com/pohsi/pktrade?status.png)](http://godoc.org/github.com/pohsi/pktrade)
[![Build Status](https://github.com/pohsi/pktrade/workflows/build/badge.svg)](https://github.com/pohsi/pktrade/actions?query=workflow%3Abuild)
[![Code Coverage](https://codecov.io/gh/pohsi/pktrade/branch/master/graph/badge.svg)](https://codecov.io/gh/pohsi/pktrade)
[![Go Report](https://goreportcard.com/badge/github.com/pohsi/pktrade)](https://goreportcard.com/report/github.com/pohsi/pktrade)

## Purpose
- This online trading platform trades 4 kinds of cards only: Pikachu, Bulbasaur, Charmander, and Squirtle.
- The price of cards is between 1.00 USD and 10.00 USD.
- Users on this platform are called traders.
- There are 10K traders.
- Traders own unlimited USD and cards.
- Traders can send orders to the platform when they want to buy or sell cards at certain prices.
- A trader can only buy or sell 1 card in 1 order.
- Traders can only buy cards using USD or sell cards for USD.
- Orders are first come first serve.
- There are 2 situations to make a trade:
    - When a buy order is sent to the platform, there exists an uncompleted sell order, whose price is the lowest one among all uncompleted sell orders and less than or equal to the price of the buy order. Then, a trade is made at the price of the sell order. Both buy and sell orders are completed. Otherwise, the buy order is uncompleted.
    - When a sell order is sent to the platform, there exists an uncompleted buy order, whose price is the highest one among all uncompleted buy orders and greater than or equal to the price of the sell order. Then, a trade is made at the price of the buy order. Both buy and sell orders are completed. Otherwise, the sell order is uncompleted.
- Traders can view the status of their latest 50 orders.
- Traders can view the latest 50 trades on each kind of cards.
- If the sequence of orders is fixed, the results must be the same no matter how many times you execute the sequence.
## Basic Requirements:
- RESTful API
- Relational database (PostgreSQL, MySQL, ...)
- Containerize
- Testing
- Gracefully shutdown
## Advanced Requirements:
- Multithreading
- Maximize performance of finishing 1M orders
- OpenAPI (Swagger)
- Set up configs using environment variables
- View logs on visualization dashboard (Kibana, Grafana, ...)
- Microservice
- Message queue (Apache Kafka, Apache Pulsar, ...)
- gRPC
- GraphQL
- Docker Compose
- Kubernetes
- Cloud computing platforms (AWS, Azure, GCP, ...) 
- NoSQL
- CI/CD
- User authentication and authorization
- High availability

## Getting Started

This demo requires
- **Go 1.16 or above**
- **Docker 17.05 or higher** for the multi-stage build support

After installing Go and Docker, run the following commands to start demo:

```shell
# download from repo
git clone https://github.com/pohsi/pktrade.git

cd pktrade

# start a PostgreSQL database server in a Docker container
make start-db

# seed the database with some test data
make testdata

# run the RESTful API server
make run

```

At this time, you have a RESTful API server running at `http://127.0.0.1:8001`. It provides the following endpoints:

* `POST /v1/login`: authenticates a user and generates a JWT(please use user name: user1 ~ user10000, password: pass as login account)
* `GET /v1/trades/records/:cardtype`: returns latest 50 trades on each kind of cards(1~4)
* `PUT /v1/trades/status/:type`: returns latest 50 for  each kind of orders(purchase: 1, sell: 2, completed: 3)

* `POST /v1/trades`: create new trades for either purchase or sell card

You may try `cURL` or some API client tools (e.g. [Postman](https://www.getpostman.com/)), you may try the following 
more complex scenarios:

```shell
# authenticate the user via: POST /v1/login
curl -X POST -H "Content-Type: application/json" -d '{"username": "user123", "password": "pass"}' http://localhost:8001/v1/login
# should return a JWT token like: {"token":"...JWT token here..."}

# with the above JWT token, try to list trades record by card type, such as: GET /v1/trades/records/1
curl -X GET -H "Authorization: Bearer ...JWT token here..." http://localhost:8001/v1/trades/records/1
# should return a list of album records in the JSON format
```