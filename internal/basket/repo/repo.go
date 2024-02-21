package repo

import (
	"github.com/el1ljah/cp_db/internal/models"
	"github.com/el1ljah/cp_db/pkg/logger"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type PgBasketRepo struct {
	Logger logger.Logger
	DB     *sqlx.DB
}

func (pbr *PgBasketRepo) Get(id int) (models.Basket, error) {
	rows, err := pbr.DB.Queryx("select * from ItemsInUsersBasket($1)", id)
	if err != nil {
		return models.Basket{}, errors.Wrap(err, "can`t get from db, query: ")
	}

	basket := models.NewBasket()

	for rows.Next() {
		item := models.OrderItem{}

		err := rows.StructScan(&item)
		if err != nil {
			return models.Basket{}, errors.Wrap(err, "can`t scan struct from db query result")
		}

		basket.Items = append(basket.Items, item)
	}

	err = pbr.DB.Get(
		&basket.ID,
		"select o.id from Ordering o where o.user_id = $1 and current_status = 'корзина'",
		id)
	if err != nil {
		return models.Basket{}, errors.Wrap(err, "can`t get from db")
	}

	var price int

	err = pbr.DB.Get(
		&price,
		"select UserBasketPrice($1)",
		id)
	if err != nil {
		return models.Basket{}, errors.Wrap(err, "can`t get from db")
	}

	basket.Price = price

	return *basket, nil
}


func (pbr *PgBasketRepo) AddItem(itemID, userID int) error {
	var succ bool
	err := pbr.DB.Get(
		&succ,
		"SELECT AddItemUsersBasket($1, $2, $3)", itemID, userID, 1)
	if !succ {
		return errors.Wrap(err, "Item is not available")
	}
	if err != nil {
		return errors.Wrap(err, "can`t add item in basket in db")
	}

	return nil
}

func (pbr *PgBasketRepo) DecItem(itemID, userID int) error {
	_, err := pbr.DB.Exec(
		"SELECT DecItemUsersBasket($1, $2, $3)", itemID, userID, 1)
	if err != nil {
		return errors.Wrap(err, "can`t dec item from basket in db")
	}

	return nil
}
