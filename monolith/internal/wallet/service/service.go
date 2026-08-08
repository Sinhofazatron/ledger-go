package service

import (
	"context"
	"net/http"

	customError "perfect_api/internal/errors"
	"perfect_api/internal/wallet/model"
	"perfect_api/internal/wallet/repository"
)

type WalletService interface {
	GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(repo repository.WalletRepository) WalletService {
	return &walletService{repo: repo}
}

func (s *walletService) GetWalletByUserID(ctx context.Context, userID string) (*model.Wallet, error) {
	w, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, customError.NewAppError(http.StatusNotFound, "WALLET_NOT_FOUND", "wallet not found")
	}

	return w, nil
}
