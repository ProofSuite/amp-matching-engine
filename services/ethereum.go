package services

type EthereumService struct {
	EthereumClient *ethclient.Client
}

func NewEthereumService(e *ethclient.Client) {
	return &EthereumService{
		EthereumClient: e
	}
}

func (s *EthereumService) WaitMined(tx *types.Transaction) {
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, s.EthereumClient, tx)
	if err != nil {
		return nil, err
	}

	return receipt, nil
}

func (s *EthereumService) GetPendingBalanceAt(a Common.Address) (*big.Int, error){
	ctx := context.Background()
	balance, err := s.EthereumClient.PendingBalanceAt(ctx, a)
	if err != nil {
		return nil, err
	}

	return balance, nil
}




