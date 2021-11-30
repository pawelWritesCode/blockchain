package blockchain

import (
	"bytes"
	"testing"
)

//User is struct that will be passed as data to blockchain
type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserDataValidator struct{}

//Validate checks whether passed data is of type User, nothing more
func (u UserDataValidator) Validate(bc *Blockchain, data interface{}) bool {
	_, ok := data.(User)
	return ok
}

type UserBlockchainValidator struct{}

//Validate checks whether first block is genesis block and whether following block Hash and PreviousHash matches
func (u UserBlockchainValidator) Validate(bc *Blockchain) bool {
	//validates first block (only genesis block)
	genesisBlock := bc.Chain[0]
	data, ok := genesisBlock.Data.(map[string]bool)
	if !ok {
		return false
	}

	isGenesis, exists := data["isGenesis"]
	if !exists {
		return false
	}

	if !isGenesis {
		return false
	}

	//validates rest of blockchain (except genesis block)
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		previousBlock := bc.Chain[i-1]

		currentBlockComputedHash, _ := CalculateHash(currentBlock)

		if !bytes.Equal(currentBlock.Hash, currentBlockComputedHash) {
			return false
		}

		if !bytes.Equal(currentBlock.PreviousHash, previousBlock.Hash) {
			return false
		}
	}

	return true
}

func TestNewBlockchain(t *testing.T) {
	udv := UserDataValidator{}
	ubv := UserBlockchainValidator{}

	genesisBlock := NewBlock([]byte("0"), map[string]bool{"isGenesis": true})
	difficulty := uint8(2)

	blockchain1, err := NewBlockchain(genesisBlock, ubv, udv, difficulty)
	if err != nil {
		t.Errorf("%v", err)
	}

	if blockchain1.difficulty != difficulty {
		t.Errorf("difficulty is not %d", difficulty)
	}

	if len(blockchain1.Chain) != 1 {
		t.Errorf("invalid number of Block in blockchain, expected 1, got: %d", len(blockchain1.Chain))
	}
}

func TestBlockchain_AddBlock(t *testing.T) {
	udv := UserDataValidator{}
	ubv := UserBlockchainValidator{}

	genesisBlock := NewBlock([]byte("0"), map[string]bool{"isGenesis": true})

	blockchain1, err := NewBlockchain(genesisBlock, ubv, udv, 2)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = blockchain1.AddBlock(User{Name: "Iwo", Age: 20})
	if err != nil {
		t.Errorf("%v", err)
	}

	err = blockchain1.AddBlock(User{Name: "Agness", Age: 22})
	if err != nil {
		t.Errorf("%v", err)
	}

	err = blockchain1.AddBlock(User{Name: "Ty", Age: 58})
	if err != nil {
		t.Errorf("%v", err)
	}

	lastBlock := blockchain1.Chain[len(blockchain1.Chain)-1]

	if lastBlock.ProofOfWork == 0 {
		t.Errorf("last block is not minned")
	}

	if bytes.Equal(lastBlock.Hash, []byte("")) || bytes.Equal(lastBlock.PreviousHash, []byte("")) {
		t.Errorf("something wrong with hashes in last block")
	}
}

func TestBlockchain_IsValid(t *testing.T) {
	udv := UserDataValidator{}
	ubv := UserBlockchainValidator{}

	genesisBlock := NewBlock([]byte("0"), map[string]bool{"isGenesis": true})

	blockchain1, err := NewBlockchain(genesisBlock, ubv, udv, 2)
	if err != nil {
		t.Errorf("%v", err)
	}

	err = blockchain1.AddBlock(User{Name: "Iwo", Age: 20})
	if err != nil {
		t.Errorf("%v", err)
	}

	err = blockchain1.AddBlock(User{Name: "Agness", Age: 22})
	if err != nil {
		t.Errorf("%v", err)
	}

	if !blockchain1.IsValid() {
		t.Errorf("something wrong with IsValid func")
	}

	blockchain1.Chain[1].Data = map[string]string{"corrupted": "data"}

	if blockchain1.IsValid() {
		t.Errorf("validation does not work")
	}
}
