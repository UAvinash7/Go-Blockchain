# Go-Blockchain
A blockchain project in Golang

Install Gorilla Tool Kit
`go get -u github.com/gorilla/mux`

Install MongoDB Driver
`go get go.mongodb.org/mongo-driver/mongo`

To run the blockchain node, run `go run main.go` command.

We can then access the blockchain interface by navigating to `http://localhost:8080` and interact with blockchain node by sending `GET` and `POST` requests i.e.,
`http://localhost:8080/blockchain` and `http://localhost:8080/write`

Until now, we have created a blockchain with the ability to add new blocks, mine blocks, and view the full blockchain, then we have created user interface rendering html templates with a basic web interface, then we used a file-based approach to store the blockchain data where the blockchain data is stored in a json file, then we replaced the file-based approach with a database by persisting the blockchain state in MongoDB tben added the logic for persisting the blockchain state in multiple server instances.

Things yet to done
1. state sync mechanism within the multiple server instances with a consensus algorithm.
2. wallet feature and token feature
3. decentrailzed storage and few more aspects.



