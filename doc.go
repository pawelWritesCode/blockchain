/*
	Package blockchain implements basic functions for custom Blockchain creation.

	Creating Blockchain can be done through

		func NewBlockchain(genesisBlock *Block, cv ChainValidator, dv DataValidator, mineDifficulty uint8) (*Blockchain, error)

		* genesis Block - first Block of custom Blockchain,
		* chain validator - entity that is able to validate whole Blockchain whether is not corrupted,
		* data validator - entity that is able to validate Block's data in context of whole Blockchain,
		* difficulty - level of mining difficulty, the greater, the more Blockchain is secure, but mining cost more time.

	Any Block's data should be serializable to JSON format.

	Checking whether Blockchain is valid or corrupted may be done thorough:
		func (b *Blockchain) IsValid() bool

	Adding new block to Blockchain can be done through:
		func (b *Blockchain) AddBlock(data interface{}) error

	and will cost time as specified in difficulty property in mining process.
*/
package blockchain
