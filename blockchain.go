package blockchain

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

//ErrValidation occurs when validation does not pass.
var ErrValidation = errors.New("validation error")

//DataValidator represents entity that has ability to verify Block's data integrity and type.
type DataValidator interface {
	//Validate validates data against set of rules
	Validate(bc *Blockchain, data interface{}) bool
}

//ChainValidator represents entity that has ability to validate data integrity of Blockchain.
type ChainValidator interface {
	//Validate validates Blockchain
	Validate(bc *Blockchain) bool
}

//Blockchain represents structure that holds blocks in fixed order.
type Blockchain struct {
	//Chain holds Blocks in fixed order
	Chain []*Block

	//chainValidator validates data integrity of Chain
	chainValidator ChainValidator

	//dataValidator validates new Block's data while adding it to Blockchain
	dataValidator DataValidator

	//difficulty is level of mining difficulty while adding new Block to Blockchain,
	//should be value greater than 0,
	//the higher value, the more time needs to pass until new Block is added to Blockchain,
	//but the more secure Blockchain becomes itself.
	difficulty uint8
}

//NewBlockchain returns pointer to Blockchain with genesis block in it.
func NewBlockchain(genesisBlock *Block, cv ChainValidator, dv DataValidator, mineDifficulty uint8) (*Blockchain, error) {
	hash, err := CalculateHash(genesisBlock)
	if err != nil {
		return &Blockchain{}, err
	}

	genesisBlock.Hash = hash

	return &Blockchain{
		Chain:          []*Block{genesisBlock},
		difficulty:     mineDifficulty,
		chainValidator: cv,
		dataValidator:  dv,
	}, nil
}

//AddBlock adds Block containing provided data to Blockchain
func (b *Blockchain) AddBlock(data interface{}) error {
	//validating incoming data
	if !b.dataValidator.Validate(b, data) {
		return fmt.Errorf("%w: data does not pass validation", ErrValidation)
	}

	//creating new block with provided data
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := NewBlock(lastBlock.Hash, data)

	//mining block according to Blockchain difficulty
	if err := newBlock.Mine(b.difficulty); err != nil {
		return err
	}

	//adding new Block to Chain
	b.Chain = append(b.Chain, newBlock)

	return nil
}

//IsValid checks whether Blockchain is valid according to set of rules.
func (b *Blockchain) IsValid() bool {
	return b.chainValidator.Validate(b)
}

//String is function that prints Blockchain in pretty format.
func (b *Blockchain) String() string {
	s := ""
	s += fmt.Sprintf("Mining difficulty: %d\n", b.difficulty)
	s += fmt.Sprintf("Number of blocks:  %d\n", len(b.Chain))

	for i, block := range b.Chain {
		serializedData, _ := json.MarshalIndent(block.Data, "\t", "\t")
		unixTimeUTC := time.Unix(block.TimeStamp, 0)

		s += fmt.Sprintf(`%d: 
{
	"data": %s,
	"hash": %x,
	"previousHash": %x,
	"proofOfWork": %d,
	"createdAt": "%s"
}
`, i, string(serializedData), block.Hash, block.PreviousHash, block.ProofOfWork, unixTimeUTC.Format(time.RFC3339))
	}

	return s
}
