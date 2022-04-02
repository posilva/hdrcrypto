package hederalib

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
)

type TokenConfig struct {
	Name                  string
	Symbol                string
	InitialSupply         uint64
	MaxTransactionFeeHbar float64
	Decimals              uint
}

type Token struct {
	Id     hedera.TokenID
	Config TokenConfig
}

func NewToken(client *HDRClient, config TokenConfig) (*Token, error) {
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
		Id:     tokenId,
		Config: config,
	}, nil
}

func (t *Token) String() string {
	return fmt.Sprintf("[%s] Id: %s ",
		t.Config.Symbol,
		t.Id.String())
}

func NewTokenFromInfo(client *HDRClient, id string) (*Token, error) {
	tokenID, err := hedera.TokenIDFromString(id)
	if err != nil {
		return nil, err
	}
	query := hedera.NewTokenInfoQuery()
	info, err := query.SetTokenID(tokenID).
		Execute(client.Get())
	if err != nil {
		return nil, err
	}
	return fromTokenInfo(info), nil
}

func fromTokenInfo(info hedera.TokenInfo) *Token {
	return &Token{
		Id: info.TokenID,
		Config: TokenConfig{
			Name:          info.Name,
			Symbol:        info.Symbol,
			Decimals:      uint(info.Decimals),
			InitialSupply: info.TotalSupply,
		},
	}
}
