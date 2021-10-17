package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"time"
)

type Blockchain struct {
	chain       []Block
	proof       float64
	previusHash string
}

type Block struct {
	index       int
	timeStamp   time.Time
	proof       float64
	previusHash string
}

func getBlockchain() (bc Blockchain) {
	bc = Blockchain{proof: 1, previusHash: "0"}
	return
}

func (bc *Blockchain) createBlock(b Block) Block {
	b = Block{
		index:       len(bc.chain),
		timeStamp:   time.Now(),
		proof:       bc.proof,
		previusHash: bc.previusHash,
	}
	bc.chain = append(bc.chain, b)
	return b
}

func (bc Blockchain) getPreviusBlock() Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc Blockchain) proofOfWord() (newProof float64, err error) {
	var checkProof bool
	newProof = 1.0
	previusBlock := bc.getPreviusBlock()

	for !checkProof {
		h := sha256.New()
		_, err = h.Write([]byte(fmt.Sprintf("%f", math.Pow(newProof, 2)-math.Pow(float64(previusBlock.proof), 2))))
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

func hash() {}

func main() {

}
