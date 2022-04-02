package hederalib

import (
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// HDRClient holds data to manage hedera client
type HDRClient struct {
	client   *hedera.Client
	operator AccountEntity
}

func NewClientForMainNet() *HDRClient {
	return &HDRClient{
		client: hedera.ClientForMainnet(),
	}
}

func NewClientForTestNet() *HDRClient {
	return &HDRClient{
		client: hedera.ClientForTestnet(),
	}
}

func (c *HDRClient) Operator(accountId string, privKey string) error {
	operatorAccountID, err := hedera.AccountIDFromString(accountId)
	if err != nil {
		return err
	}

	operatorPrivKey, err := hedera.PrivateKeyFromString(privKey)
	if err != nil {
		return err
	}

	c.client.SetOperator(operatorAccountID, operatorPrivKey)
	c.operator = AccountEntity{
		Account:    operatorAccountID,
		PrivateKey: operatorPrivKey,
	}
	return nil
}

func (c *HDRClient) Get() *hedera.Client {
	return c.client
}

func (c *HDRClient) OperatorAccount() hedera.AccountID {
	return c.operator.Account
}

func (c *HDRClient) OperatorPrivateKey() hedera.PrivateKey {
	return c.operator.PrivateKey
}
