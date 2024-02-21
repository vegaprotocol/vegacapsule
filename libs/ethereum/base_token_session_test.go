package ethereum

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestBaseTokenSession_prepareRawMintData(t *testing.T) {
	want := "6000603460be8239805160601c9060145160183384609f565b908181116043575b82806044818088602f3082609f565b5063a9059cbb60e01b600052336004525af1005b63de5f72fd60e01b600052600492918280858180895af15060633086609f565b63de5f72fd60e01b60005291030491815b8381105a61ea6010161560925760019083808481808a5af150016074565b5090915060449050816020565b6024600081926044946370a0823160e01b83526004525afa5060245190563773a5c7aff77e014cbf067dd31801b4c6dc4136000000000000000000000000000000000000000000000000000006faf9a60e00"

	ts := &BaseTokenSession{
		address: common.HexToAddress("0x3773A5c7aFF77e014cBF067dd31801b4C6dc4136"),
	}

	targetBalance := big.NewInt(7675000000000)

	got, err := ts.prepareMintBytecode(targetBalance)
	if err != nil {
		t.Errorf("prepareMintBytecode() error = %s", err)
		return
	}

	if !reflect.DeepEqual(fmt.Sprintf("%x", got), want) {
		t.Errorf("prepareMintBytecode() got = %x, want %s", got, want)
	}
}
