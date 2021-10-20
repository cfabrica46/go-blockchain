package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"
)

type Blockchain struct {
	Chain       []Block
	Proof       float64
	PreviusHash string
}

type Block struct {
	Index       int
	TimeStamp   time.Time
	Proof       float64
	PreviusHash string
}

func getBlockchain() (bc Blockchain) {
	bc = Blockchain{Proof: 1, PreviusHash: "0"}
	return
}

func (bc *Blockchain) createBlock(proof float64, previusHash string) Block {
	b := Block{
		Index:       len(bc.Chain),
		TimeStamp:   time.Now(),
		Proof:       proof,
		PreviusHash: previusHash,
	}
	bc.Chain = append(bc.Chain, b)
	return b
}

func (bc Blockchain) getPreviusBlock() Block {
	if len(bc.Chain) != 0 {
		return bc.Chain[len(bc.Chain)-1]
	}
	return Block{}
}

func (bc Blockchain) proofOfWork() (newProof float64, err error) {
	var checkProof bool
	newProof = 1.0
	previusBlock := bc.getPreviusBlock()

	for !checkProof {
		h := sha256.New()

		_, err = h.Write([]byte(fmt.Sprintf("%f", math.Pow(newProof, 2)-math.Pow(previusBlock.Proof, 2))))
		if err != nil {
			return 0, err
		}

		if hex.EncodeToString(h.Sum(nil))[:4] == "0000" {
			checkProof = true
		} else {
			newProof += 1
		}
	}
	return newProof, nil
}

func (b Block) getHash() (string, error) {
	encodedBlock, err := json.Marshal(b)
	if err != nil {
		return "", nil
	}

	h := sha256.New()

	_, err = h.Write(encodedBlock)
	if err != nil {
		return "", nil
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func chanValid(chain []Block) bool {
	for i := range chain {
		if i != 0 {
			previusHash, err := chain[i-1].getHash()
			if err != nil {
				return false
			}

			if chain[i].PreviusHash != previusHash {
				proof := chain[i].Proof
				h := sha256.New()
				_, err = h.Write([]byte(fmt.Sprintf("%f", math.Pow(proof, 2)-math.Pow(chain[i-1].Proof, 2))))
				if err != nil {
					return false
				}
				if hex.EncodeToString(h.Sum(nil))[:4] == "0000" {
					return false
				}

			}
		}
	}
	return true
}

var blockChain = getBlockchain()

func main() {
	http.HandleFunc("/mine_block", mineBlock)
	http.HandleFunc("/get_chain", getChain)
	http.HandleFunc("/valid", valid)

	log.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func mineBlock(w http.ResponseWriter, r *http.Request) {
	previusBlock := blockChain.getPreviusBlock()
	// previusProof := previusBlock.proof
	proof, err := blockChain.proofOfWork()
	if err != nil {
		log.Println(err)
		return
	}

	previusHash, err := previusBlock.getHash()
	if err != nil {
		log.Println(err)
		return
	}

	block := blockChain.createBlock(proof, previusHash)

	data, err := json.MarshalIndent(block, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func getChain(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Chain  []Block
		Length int
	}{blockChain.Chain, len(blockChain.Chain)}

	data, err := json.MarshalIndent(response, "", " ")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func valid(w http.ResponseWriter, r *http.Request) {
	check := chanValid(blockChain.Chain)
	if check {
		_, err := w.Write([]byte("VALIDADA"))
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		_, err := w.Write([]byte("NO VALIDADA"))
		if err != nil {
			log.Println(err)
			return
		}
	}
}
