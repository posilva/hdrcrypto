package hedera

import (
	"fmt"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

type AccountEntity struct {
	Account    hedera.AccountID
	PrivateKey hedera.PrivateKey
}

type TokenConfig struct {
	Name                  string
	Symbol                string
	InitialSupply         uint64
	MaxTransactionFeeHbar float64
	Decimals              uint
}

type Token struct {
	id     hedera.TokenID
	config TokenConfig
}

func CreateToken(client *HDRClient, config TokenConfig) (*Token, error) {

	tokenCreateTransaction, err := hedera.NewTokenCreateTransaction().
		SetTokenName(config.Name).
		SetTokenSymbol(config.Symbol).
		SetTreasuryAccountID(client.operator.Account).
		SetInitialSupply(config.InitialSupply).
		SetAdminKey(client.operator.PrivateKey.PublicKey()).
		SetSupplyType(hedera.TokenSupplyTypeInfinite).
		SetSupplyKey(client.operator.PrivateKey.PublicKey()).
		SetMaxTransactionFee(hedera.NewHbar(config.MaxTransactionFeeHbar)).
		SetDecimals(config.Decimals).
		FreezeWith(client.Get())

	if err != nil {
		panic(err)
	}

	txResponse, err := tokenCreateTransaction.
		Sign(client.operator.PrivateKey).
		Sign(client.operator.PrivateKey).
		Execute(client.Get())

	if err != nil {
		panic(err)
	}

	receipt, err := txResponse.GetReceipt(client.Get())
	if err != nil {
		panic(err)
	}

	tokenId := *receipt.TokenID
	hedera.NewCustomFixedFee().
		SetAmount(1).
		SetDenominatingTokenID(tokenId)

	return &Token{
		id:     tokenId,
		config: config,
	}, nil
}

func (t *Token) Display() {
	fmt.Printf("[%s] Id: %s ",
		t.config.Symbol,
		t.id.String())
}

func CreateAccountEntity(client *HDRClient, initialBalance float64) (*AccountEntity, error) {
	treasuryKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		return nil, err
	}

	treasuryAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(treasuryKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(initialBalance)).
		Execute(client.Get())
	if err != nil {
		return nil, err
	}

	receipt, err := treasuryAccount.GetReceipt(client.Get())
	if err != nil {
		return nil, err
	}

	return &AccountEntity{
		Account:    *receipt.AccountID,
		PrivateKey: treasuryKey,
	}, nil
}

func (a AccountEntity) Display(name string) {
	fmt.Printf("%s Account: %s, PrivateKey: %s \n",
		name,
		a.Account.String(),
		a.PrivateKey.StringRaw())
}
