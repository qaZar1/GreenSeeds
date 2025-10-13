package service

import (
	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddReceipts(receipts models.Receipts) (models.Receipts, error) {
	if err := s.validate.Struct(receipts); err != nil {
		return models.Receipts{}, err
	}

	return s.repo.RptRepo.AddReceipts(receipts)
}

func (s *Service) GetReceipts() ([]models.Receipts, error) {
	return s.repo.RptRepo.GetReceipts()
}

func (s *Service) GetReceiptsByReceipt(receipt string) (models.Receipts, error) {
	return s.repo.RptRepo.GetReceiptsByReceipt(receipt)
}

func (s *Service) UpdateReceipts(receipts models.Receipts) (models.Receipts, error) {
	if err := s.validate.Struct(receipts); err != nil {
		return models.Receipts{}, err
	}

	return s.repo.RptRepo.UpdateReceipts(receipts)
}

func (s *Service) DeleteReceipts(receipt string) (bool, error) {
	return s.repo.RptRepo.DeleteReceipts(receipt)
}
