package service

import (
	"github.com/el1ljah/cp_db/internal/models"
	"github.com/el1ljah/cp_db/pkg/logger"
	"github.com/pkg/errors"
)

type OrderRepo interface {
	Commit(int) error
	Get(int) (models.Order, error)
	GetUsersAll(int) ([]models.Order, error)
	GetAll() ([]models.Order, error)
	Update(models.Order) (models.Order, error)
}

type OrderService struct {
	OrderRepo OrderRepo
	Logger    logger.Logger
}

func (os OrderService) Get(id int) (models.Order, error) {
	order, err := os.OrderRepo.Get(id)
	if err != nil {
		return models.Order{}, errors.Wrap(err, "can`t get from repo")
	}

	return order, nil
}

func (bs OrderService) Commit(id int) error  {
	err := bs.OrderRepo.Commit(id)
	if err != nil {
		return errors.Wrap(err, "can`t commit in repo")
	}

	return  nil
}

func (os OrderService) GetUsersAll(user int) ([]models.Order, error) {
	orders, err := os.OrderRepo.GetUsersAll(user)
	if err != nil {
		return nil, errors.Wrap(err, "can`t get from repo")
	}

	return orders, nil
}

func (os OrderService) GetAll() ([]models.Order, error) {
	orders, err := os.OrderRepo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "can`t get from repo")
	}

	return orders, nil
}

func (os OrderService) Update(order models.Order) (models.Order, error) {
	order, err := os.OrderRepo.Update(order)
	if err != nil {
		return order, errors.Wrap(err, "can`t update repo")
	}

	return order, nil
}
