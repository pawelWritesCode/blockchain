package blockchain

import "bytes"

//Blockchain represents structure that holds block in fixed order
type Blockchain struct {
	//Chain holds blocks in fixed order
	Chain []*Block
}

//NewBlockchain returns pointer to Blockchain with genesis block in it
func NewBlockchain() *Blockchain {
	genesisBlock, _ := NewBlock([]byte("0"), map[string]bool{"isGenesis": true})

	b := Blockchain{Chain: []*Block{genesisBlock}}

	return &b
}

//AddBlock adds block of data to blockchain
func (b *Blockchain) AddBlock(data interface{}) error {
	lastBlock := b.Chain[len(b.Chain) - 1]
	newBlock, err := NewBlock(lastBlock.Hash, data)
	if err != nil {
		return err
	}

	newBlock.Mine(1)
	b.Chain = append(b.Chain, newBlock)

	return nil
}

//IsValid checks whether blockchain is valid according to set of rules
func (b * Blockchain) IsValid() bool {
	for i := 1; i < len(b.Chain); i++ {
		currentBlock := b.Chain[i]
		previousBlock := b.Chain[i - 1]

		currentBlockComputedHash, _ := CalculateHash(*currentBlock)

		if !bytes.Equal(currentBlock.Hash, currentBlockComputedHash) {
			return false
		}

		if !bytes.Equal(currentBlock.PreviousHash, previousBlock.Hash) {
			return false
		}
	}

	return true
}