package hederalib

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

type AccountEntity struct {
	Name       string
	Account    hedera.AccountID
	PrivateKey hedera.PrivateKey
}

func NewAccount(client *HDRClient, name string, initialBalance float64) (*AccountEntity, error) {
	key, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		return nil, err
	}

	account, err := hedera.NewAccountCreateTransaction().
		SetKey(key.PublicKey()).
		SetInitialBalance(hedera.NewHbar(initialBalance)).
		Execute(client.Get())
	if err != nil {
		return nil, err
	}

	receipt, err := account.GetReceipt(client.Get())
	if err != nil {
		return nil, err
	}

	return &AccountEntity{
		Name:       name,
		Account:    *receipt.AccountID,
		PrivateKey: key,
	}, nil
}

func (a AccountEntity) String() string {
	return fmt.Sprintf("[%s] Account: %s, PrivateKey: %s",
		a.Name,
		a.Account.String(),
		a.PrivateKey.StringRaw())
}
