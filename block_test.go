package blockchain

import (
	"bytes"
	"testing"
)

func TestNewBlock(t *testing.T) {
	type args struct {
		previousHash []byte
		data         interface{}
	}
	tests := []struct {
		name string
		args args
		want *Block
	}{
		{name: "no data", args: args{
			previousHash: nil,
			data:         nil,
		}, want: &Block{
			Data:         nil,
			Hash:         nil,
			PreviousHash: nil,
			TimeStamp:    0,
			ProofOfWork:  0,
		}},
		{name: "any data", args: args{
			previousHash: []byte("abc"),
			data:         "abc",
		}, want: &Block{
			Data:         "abc",
			Hash:         nil,
			PreviousHash: []byte("abc"),
			TimeStamp:    0,
			ProofOfWork:  0,
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlock(tt.args.previousHash, tt.args.data); !bytes.Equal(got.PreviousHash, tt.args.previousHash) || got.ProofOfWork != 0 || got.Data != tt.args.data {
				t.Errorf("NewBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
