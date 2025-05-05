package actionlogman

import (
	"log"

	"gorm.io/gorm"
)

type Service struct {
	infoLog  *log.Logger
	errorLog *log.Logger
	db       *gorm.DB
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{
		infoLog:  infoLog,
		errorLog: errorLog,
		db:       db,
	}
}

func (s *Service) parseFilter(filter *Filter) *gorm.DB {
	query := s.db

	if filter == nil {
		return query
	}

	if filter.Action != "" {
		query = query.Where("action=?", filter.Action)
	}
	if filter.RefID > 0 {
		query = query.Where("ref_id=?", filter.RefID)
	}
	if filter.CustomerID > 0 {
		query = query.Where("user_id=?", filter.CustomerID)
	}
	if len(filter.Actions) > 0 {
		query = query.Where("action IN ?", filter.Actions)
	}
	if len(filter.RefIDs) > 0 {
		query = query.Where("ref_id IN ?", filter.RefIDs)
	}
	if len(filter.CustomerIDs) > 0 {
		query = query.Where("user_id IN ?", filter.CustomerIDs)
	}

	return query
}

func (s *Service) Count(filter *Filter) (int, error) {
	query := s.parseFilter(filter)

	var count int64
	if err := query.Model(&ActionLog{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *Service) GetAll(filter *AdminActionFilter, page, size int, preloads ...string) ([]*ActionLog, int, error) {
	query := s.db

	if len(filter.Actions) > 0 {
		query = query.Where("action IN ?", filter.Actions)
	}

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	var count int64
	if err := query.Model(&ActionLog{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if size > 0 {
		query = query.Limit(size)
		if page > 0 {
			query = query.Offset((page - 1) * size)
		}
	}

	var actionLogs []*ActionLog
	if err := query.Order("created_at desc").Find(&actionLogs).Error; err != nil {
		return nil, 0, err
	}

	return actionLogs, int(count), nil
}

func (s *Service) Save(data *ActionLog) (*ActionLog, error) {
	if err := s.db.Save(data).Error; err != nil {
		return nil, err
	}
	return data, nil
}
