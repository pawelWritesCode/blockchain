package blockchain

import (
	"bytes"
	"testing"
)

func testGenesisBlock(bc *Blockchain) func(t *testing.T) {
	return func(t *testing.T) {
		if len(bc.Chain) != 1 {
			t.Errorf("initialized blockchain lenght should be 1, got: %d", len(bc.Chain))
		}

		genesisBlock := bc.Chain[0]
		if !bytes.Equal(genesisBlock.PreviousHash, []byte("0")) {
			t.Errorf("genesis block previous hash doesnt match expected: 0")
		}

		expectedData, ok := genesisBlock.Data.(map[string]bool)
		if !ok {
			t.Errorf("genesis block Data is not map[string]bool")
		}

		isGenesis, exists := expectedData["isGenesis"]
		if !exists {
			t.Errorf("genesis block Data does not contain key: isGenesis")
		}

		if !isGenesis {
			t.Errorf("genesis block Data key 'isGenesis' has value false")
		}
	}
}

func testRegularBlock(data map[string]string, bc *Blockchain) func(t *testing.T) {
	return func(t *testing.T) {
		if err := bc.AddBlock(data); err != nil {
			t.Errorf("could not add second block to blockchain, err: %v", err)
		}

		firstBlock := bc.Chain[len(bc.Chain) -1]
		genesisBlock := bc.Chain[len(bc.Chain) - 2]

		if !bytes.Equal(firstBlock.PreviousHash, genesisBlock.Hash) {
			t.Errorf("first block PreviousHash differs from genesis block Hash")
		}

		if ! (firstBlock.ProofOfWork > 0) {
			t.Errorf("first block proof of work should be greater than 0")
		}

		expectedData, ok := firstBlock.Data.(map[string]string)
		if !ok {
			t.Errorf("genesis block Data is not map[string]bool")
		}

		name, exists := expectedData["name"]
		if !exists {
			t.Errorf("block does not contain 'name' property")
		}

		if name != data["name"] {
			t.Errorf("block data name property: %s is not what expected: %s", name, data["name"])
		}
	}
}

func TestNewBlockchain(t *testing.T) {
	blockchain1 := NewBlockchain()
	testGensisBlockFunc := testGenesisBlock(blockchain1)
	t.Run("initialization blockchain with genesis block", testGensisBlockFunc)
}

func TestBlockchain_AddBlock(t *testing.T) {
	blockchain2 := NewBlockchain()
	testGensisBlockFunc2 := testGenesisBlock(blockchain2)

	t.Run("initialize blockchain with 2 blocks, one genesis, second other", func(t *testing.T) {
		testGensisBlockFunc2(t)
		addRegularBlockFunc := testRegularBlock(map[string]string{"name": "pawel"}, blockchain2)
		addRegularBlockFunc(t)
		addRegularBlockFunc2 := testRegularBlock(map[string]string{"name": "val"}, blockchain2)
		addRegularBlockFunc2(t)
	})
}

func TestBlockchain_IsValid(t *testing.T) {
	blockchain2 := NewBlockchain()
	testGensisBlockFunc2 := testGenesisBlock(blockchain2)

	t.Run("initialize blockchain with 2 blocks, one genesis, second other", func(t *testing.T) {
		testGensisBlockFunc2(t)
		addRegularBlockFunc := testRegularBlock(map[string]string{"name": "pawel"}, blockchain2)
		addRegularBlockFunc(t)
		addRegularBlockFunc2 := testRegularBlock(map[string]string{"name": "val"}, blockchain2)
		addRegularBlockFunc2(t)

		if !blockchain2.IsValid() {
			t.Errorf("blockchain is not valid")
		}
	})
}