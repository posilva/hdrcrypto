package hederalib

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/rs/zerolog/log"
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

func TokenAssociate(client *HDRClient, tokenId string, account string, key string) error {
	k, err := hedera.PrivateKeyFromString(key)
	if err != nil {
		return err
	}
	acc, err := hedera.AccountIDFromString(account)
	if err != nil {
		return err
	}
	token, err := hedera.TokenIDFromString(tokenId)
	if err != nil {
		return err
	}

	association, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(acc).
		SetTokenIDs(token).
		FreezeWith(client.Get())
	if err != nil {
		return err
	}
	signTx := association.Sign(k)

	associationTxSubmit, err := signTx.Execute(client.Get())
	if err != nil {
		return err
	}

	associationRx, err := associationTxSubmit.GetReceipt(client.Get())
	if err != nil {
		return err
	}

	if associationRx.Status != hedera.StatusSuccess {
		return fmt.Errorf("failed to associate: %v", associationRx.Status)
	}
	return nil
}

func TokenTransfer(client *HDRClient, tokenId string, toAccount string, amount float64) error {
	to, err := hedera.AccountIDFromString(toAccount)
	if err != nil {
		return err
	}
	token, err := hedera.TokenIDFromString(tokenId)
	if err != nil {
		return err
	}

	log.Logger.Info().Msgf("tranfer %v to %v", amount, to.Account)
	transaction, err := hedera.NewTransferTransaction().
		AddTokenTransfer(token, client.OperatorAccount(), int64(amount*-1)).
		AddTokenTransfer(token, to, int64(amount)).
		FreezeWith(client.Get())

	if err != nil {
		panic(err)
	}
	signTransferTx := transaction.Sign(client.OperatorPrivateKey())
	tokenTransferSubmit, err := signTransferTx.Execute(client.Get())
	if err != nil {
		panic(err)
	}
	receipt, err := tokenTransferSubmit.GetReceipt(client.Get())
	if err != nil {
		panic(err)
	}

	//Get the transaction consensus status
	transactionStatus := receipt.Status

	fmt.Printf("The transaction consensus status is %v\n", transactionStatus)

	return nil

}
func TokenHbarTransfer(client *HDRClient, fromAccount string, toAccount string, amount float64) error {
	from, err := hedera.AccountIDFromString(fromAccount)
	if err != nil {
		return err
	}
	to, err := hedera.AccountIDFromString(toAccount)
	if err != nil {
		return err
	}

	// Create a transaction to transfer 100 hbars
	transaction := hedera.NewTransferTransaction().
		AddHbarTransfer(from, hedera.NewHbar(amount*-1)).
		AddHbarTransfer(to, hedera.NewHbar(amount))

	//Submit the transaction to a Hedera network
	txResponse, err := transaction.Execute(client.Get())

	if err != nil {
		panic(err)
	}

	//Request the receipt of the transaction
	receipt, err := txResponse.GetReceipt(client.Get())

	if err != nil {
		return err
	}

	//Get the transaction consensus status
	transactionStatus := receipt.Status

	fmt.Printf("The transaction consensus status is %v\n", transactionStatus)

	return nil
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
