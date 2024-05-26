package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

// Block represents a block in the blockchain
type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

// Blockchain is a slice of blocks
var Blockchain []Block
var mutex = &sync.Mutex{}

func calculateHash(block Block) string {
	record := string(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func generateBlock(oldBlock Block, data string) Block {
	var newBlock Block

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = time.Now().String()
	newBlock.Data = data
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateHash(newBlock)

	return newBlock
}

func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

func replaceChain(newBlocks []Block) {
	if len(newBlocks) > len(Blockchain) {
		Blockchain = newBlocks
	}
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s\n", string(bytes))
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var m Block

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	mutex.Lock()
	if isBlockValid(m, Blockchain[len(Blockchain)-1]) {
		newBlockchain := append(Blockchain, m)
		replaceChain(newBlockchain)
		saveBlockchain()
		log.Println("Block added")
	}
	mutex.Unlock()

	respondWithJSON(w, r, http.StatusCreated, m)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func saveBlockchain() {
	mutex.Lock()
	bytes, err := json.MarshalIndent(Blockchain, "", "  ")
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = ioutil.WriteFile("blockchain.json", bytes, 0644)
	if err != nil {
		log.Println(err.Error())
	}
	mutex.Unlock()
}

func loadBlockchain() {
	mutex.Lock()
	defer mutex.Unlock()

	if _, err := os.Stat("blockchain.json"); os.IsNotExist(err) {
		genesisBlock := Block{0, time.Now().String(), "Genesis Block", "", ""}
		genesisBlock.Hash = calculateHash(genesisBlock)
		Blockchain = append(Blockchain, genesisBlock)
		saveBlockchain()
		return
	}

	bytes, err := ioutil.ReadFile("blockchain.json")
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = json.Unmarshal(bytes, &Blockchain)
	if err != nil {
		log.Println(err.Error())
	}
}

func main() {
	loadBlockchain()

	r := mux.NewRouter()
	r.HandleFunc("/blockchain", handleGetBlockchain).Methods("GET")
	r.HandleFunc("/write", handleWriteBlock).Methods("POST")

	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
