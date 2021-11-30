package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

//Block represents structure that holds user passed data and internal metadata.
type Block struct {
	//Data is user passed arbitrary data that should be serializable to JSON format.
	Data interface{}

	//Hash represents Block sha256 hash value.
	Hash []byte

	//PreviousHash holds value of previous Blockchain's Block sha256 hash value.
	PreviousHash []byte

	//TimeStamp represents UNIX timestamp of Block creation moment.
	TimeStamp int64

	//ProofOfWork represents work required to Mine this Block.
	ProofOfWork uint64
}

//NewBlock returns pointer to Block with following error from hashing func.
//Data should be serializable to JSON format.
func NewBlock(previousHash []byte, data interface{}) *Block {
	return &Block{
		Data:         data,
		PreviousHash: previousHash,
		TimeStamp:    time.Now().Unix(),
		ProofOfWork:  0,
	}
}

//CalculateHash is function that calculates sha256 hash from Block.
func CalculateHash(b *Block) ([]byte, error) {
	jsonData, err := json.Marshal(b.Data)
	if err != nil {
		return []byte(""), err
	}

	timeStampString := strconv.Itoa(int(b.TimeStamp))
	proofOfWorkString := strconv.Itoa(int(b.ProofOfWork))

	h := sha256.New()
	h.Write(append(append([]byte(timeStampString+proofOfWorkString), jsonData...), b.PreviousHash...))

	return h.Sum(nil), nil
}

//Mine is function that finds proper sha256 hash.
//Internally it increases Block's ProofOfWork property, up to time when proper Hash is computed.
//The higher is value of difficulty, the more time it cost to mine Block.
func (b *Block) Mine(difficulty uint8) error {
	val := "0"
	for !bytes.HasPrefix(b.Hash, []byte(strings.Repeat(val, int(difficulty)))) {
		b.ProofOfWork++
		newHash, err := CalculateHash(b)
		if err != nil {
			return err
		}

		b.Hash = newHash
	}

	return nil
}
