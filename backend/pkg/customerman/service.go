package customerman

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

// #region Base
type Service struct {
	DB       *gorm.DB
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewService(db *gorm.DB, infoLog, errorLog *log.Logger) *Service {
	return &Service{
		DB:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
}

func (s *Service) parseFilter(filter *Filter) *gorm.DB {
	query := s.DB

	if filter == nil {
		return query
	}

	if filter.Keyword != "" {
		query = query.Where("(name ILIKE ? || '%%' OR email ILIKE ? || '%%')", filter.Keyword)
	}

	if len(filter.CustomerIDs) > 0 {
		query = query.Where("id IN (?)", filter.CustomerIDs)
	}

	return query
}

// #endregion Base

func (s *Service) GetList(filter *Filter, offset, limit int, preloads ...string) ([]*Customer, error) {
	var customers []*Customer

	query := s.parseFilter(filter)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&customers).Error; err != nil {
		return nil, err
	}

	return customers, nil
}

func (s *Service) Count(filter *Filter) (int, error) {
	var count int64

	query := s.parseFilter(filter)

	if err := query.Model(new(Customer)).Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func (s *Service) Get(id int) (*Customer, error) {
	var customer *Customer

	if err := s.DB.First(&customer, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return customer, nil
}

func (s *Service) GetWithEmail(email string) (*Customer, error) {
	var customer *Customer

	if err := s.DB.First(&customer, "email=?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	}

	return customer, nil
}

func (s *Service) Save(p *Customer) (*Customer, error) {
	if err := s.DB.Save(&p).Error; err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Delete(id int) error {
	if err := s.DB.Delete(new(Customer), id).Error; err != nil {
		return err
	}
	return nil
}
