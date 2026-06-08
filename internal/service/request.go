package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/niiilov/go-dog-trapping/internal/dto"
	validate "github.com/niiilov/go-dog-trapping/pkg/validator"
)

func (s *Service) CreateRequest(request *dto.CreateRequestDTO) error {
	err := validate.Validate(request)
	if err != nil {
		return ErrInvalidData
	}

	return s.repository.CreateRequest(request)
}

func (s *Service) GetAllRequests() ([]*dto.GetRequestsDTO, error) {
	requests, err := s.repository.GetAllRequests()
	if err != nil {
		return nil, err
	}
	s.validateDelay(requests)
	return requests, nil
}

func (s *Service) GetRequestsByTerOtdel(id string) ([]*dto.GetRequestsDTO, error) {
	requests, err := s.repository.GetRequestsByTerOtdel(id)
	if err != nil {
		return nil, err
	}
	s.validateDelay(requests)
	return requests, nil
}

func (s *Service) GetRequestsByIDs(ids []string) ([]*dto.GetRequestsDTO, error) {
	return s.repository.GetRequestsByIDs(ids)
}

func (s *Service) GetRequestsByTerOtdelIDs(terOtdelID string, ids []string) ([]*dto.GetRequestsDTO, error) {
	return s.repository.GetRequestsByTerOtdelIDs(terOtdelID, ids)
}

func (s *Service) ChangeStatusRequest(req *dto.ChangeStatusRequestDTO) error {
	return s.repository.ChangeStatusRequest(req)
}

func (s *Service) DeleteRequest(id string) error {
	return s.repository.DeleteRequest(id)
}

func (s *Service) AddActFile(reqID string, actFilename, actFilePath string) error {
	if err := s.storage.UploadFile(context.TODO(), actFilename, actFilePath); err != nil {
		fmt.Println("Error upload file to S3:", err)
		return err
	}
	os.Remove(actFilePath)
	var statusReq dto.ChangeStatusRequestDTO
	statusReq.ID = reqID
	statusReq.Status = "Выполнена"
	if err := s.ChangeStatusRequest(&statusReq); err != nil {
		return err
	}

	return s.repository.AddActFile(reqID, actFilename)
}

func (s *Service) validateDelay(requests []*dto.GetRequestsDTO) error {
	twoWeeks := 14 * 24 * time.Hour

	for _, req := range requests {
		if req == nil {
			continue
		}
		if req.Status == "Завершена" || req.Status == "Просрочена" {
			continue
		}
		if time.Since(req.CreatedAt) > twoWeeks {
			statusReq := &dto.ChangeStatusRequestDTO{
				ID:     req.ID,
				Status: "Просрочена",
			}

			if err := s.repository.ChangeStatusRequest(statusReq); err != nil {
				fmt.Println("failed to change status for request", req.ID, ":", err)
			}
		}
	}

	return nil
}
