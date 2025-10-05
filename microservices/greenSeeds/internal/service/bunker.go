package service

import (
	"net/http"
	"strconv"

	"github.com/qaZar1/GreenSeeds/microservices/greenSeeds/internal/models"
)

func (s *Service) AddBunker(bunker models.Bunkers) (int, error) {
	if err := s.validate.Struct(bunker); err != nil {
		return http.StatusBadRequest, ErrValidateStruct
	}

	ok, err := s.repo.BunkRepo.AddBunkers(bunker)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusInternalServerError, ErrInvalidAddBunker
	}

	return http.StatusNoContent, nil
}

func (s *Service) GetBunkers() ([]models.Bunkers, error) {
	return s.repo.BunkRepo.GetBunkers()
}

func (s *Service) UpdateBunker(bunker models.Bunkers) (int, error) {
	if err := s.validate.Struct(bunker); err != nil {
		return http.StatusBadRequest, ErrValidateStruct
	}

	ok, err := s.repo.BunkRepo.UpdateBunkers(bunker)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !ok {
		return http.StatusInternalServerError, ErrInvalidUpdateData
	}

	return http.StatusNoContent, nil
}

func (s *Service) DeleteBunker(bunkerId string) error {
	bunkerIdInt, err := strconv.Atoi(bunkerId)
	if err != nil {
		return err
	}

	ok, err := s.repo.BunkRepo.DeleteBunkers(bunkerIdInt)
	if err != nil {
		return err
	}

	if !ok {
		return ErrInvalidDeleteData
	}

	return nil
}
