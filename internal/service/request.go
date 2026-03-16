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

func (s *Service) GetRequestsByTerOtdel(id string) ([]*dto.GetRequestsDTO, error) {

	requests, err := s.repository.GetRequestsByTerOtdel(id)
	if err != nil {
		return nil, err
	}
	s.validateDelay(requests)
	return requests, nil
}

func (s *Service) GetRequestsByDistrictID(id string) ([]*dto.GetRequestsDTO, error) {
	return s.repository.GetRequestsByDistrictID(id)
}

func (s *Service) GetRequestsByDistrictIDs(district_id string, ids []string) ([]*dto.GetRequestsDTO, error) {
	return s.repository.GetRequestsByDistrictIDs(district_id, ids)
}

func (s *Service) ChangeStatusRequest(req *dto.ChangeStatusRequestDTO) error {
	return s.repository.ChangeStatusRequest(req)
}

func (s *Service) DeleteRequest(id string) error {
	return s.repository.DeleteRequest(id)
}

func (s *Service) AddActFile(reqID string, actFilename, actFilePath string) error {
	if err := s.storage.UploadFile(context.TODO(), actFilename, actFilePath); err != nil {
		//лог
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
		// если прошло больше двух недель с момента создания — меняем статус
		if time.Since(req.CreatedAt) > twoWeeks {

			statusReq := &dto.ChangeStatusRequestDTO{
				ID:     req.ID,
				Status: "Просрочена", // <-- поменяйте на нужный вам статус
			}

			if err := s.repository.ChangeStatusRequest(statusReq); err != nil {
				fmt.Println("failed to change status for request", req.ID, ":", err)
				// продолжаем обработку остальных заявок
			}
		}
	}

	return nil

}
