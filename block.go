package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

//Block represents structure that holds user passed data and internal metadata
type Block struct {
	//Data is user passed arbitrary data
	Data interface{}

	//Hash represents block hash
	Hash []byte

	//PreviousHash holds value of previous blockchain's block hash value
	PreviousHash []byte

	//TimeStamp represents UNIX data timestamp of block creation
	TimeStamp int64

	//ProofOfWork represents number of work required to Mine this block
	ProofOfWork uint64
}

//NewBlock returns pointer to Block with following error from hashing func
func NewBlock(previousHash []byte, data interface{}) (*Block, error) {
	b := Block{
		Data: data,
		PreviousHash: previousHash,
		TimeStamp: time.Now().Unix(),
		ProofOfWork: uint64(0),
	}

	hash, err := CalculateHash(b)
	if err != nil {
		return &Block{}, err
	}

	b.Hash = hash

	return &b, nil
}

//Mine is function that calculates hash according to difficulty level
//internally it increases ProofOfWork property of Block each time algorithm enters internal loop,
//up to time when proper Hash is computed
func (b *Block) Mine(difficulty uint8) {
	val := "0"
	for !bytes.HasPrefix(b.Hash, []byte(strings.Repeat(val, int(difficulty)))) {
		b.ProofOfWork++
		newHash, _ := CalculateHash(*b)
		b.Hash = newHash
	}
}

//CalculateHash is function that calculates hash from internal Block's data
func CalculateHash(b Block) ([]byte, error) {
	jsonData, err := json.Marshal(b.Data)
	if err != nil {
		return []byte(""), err
	}

	timeStampString := strconv.Itoa(int(b.TimeStamp))
	proofOfWorkString := strconv.Itoa(int(b.ProofOfWork))

	h := sha256.New()
	h.Write(append(append([]byte(timeStampString + proofOfWorkString), jsonData...,), b.PreviousHash...))

	return h.Sum(nil), nil
}