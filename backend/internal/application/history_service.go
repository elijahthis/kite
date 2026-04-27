package application

import (
	"context"

	"github.com/elijahthis/kite/internal/domain"
	"github.com/google/uuid"
)

type HistoryService struct {
	ledgerRepo domain.LedgerRepository
}

func NewHistoryService(lr domain.LedgerRepository) *HistoryService {
	return &HistoryService{ledgerRepo: lr}
}

func (s *HistoryService) GetUserHistory(ctx context.Context, userID uuid.UUID, page, limit int) ([]domain.TransactionHistory, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	return s.ledgerRepo.GetTransactionHistory(ctx, userID, limit, offset)
}
