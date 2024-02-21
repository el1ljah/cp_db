package delivery

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/el1ljah/cp_db/internal/models"
	"github.com/el1ljah/cp_db/pkg/logger"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type OrderService interface {
	Get(int) (models.Order, error)
	GetUsersAll(int) ([]models.Order, error)
	Commit(int) error
	Update(models.Order) (models.Order, error)
}

type ContextManager interface {
	UserIDFromContext(ctx context.Context) (int, error)
}

type OrderHandler struct {
	OrderService   OrderService
	ContextManager ContextManager
	Logger         logger.Logger
}

// @Summary      Commit purchase
// @Tags         orders
// @Accept       json
// @Produce      json
// @Success      200
// @Failure      401
// @Failure      404
// @Failure      500
// @Security ApiKeyAuth
// @Router       /orders [post]
func (bh *OrderHandler) Commit(w http.ResponseWriter, r *http.Request) {
	userID, err := bh.ContextManager.UserIDFromContext(r.Context())
	if err != nil {
		bh.Logger.Errorw("fail to get id from context",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = bh.OrderService.Commit(userID)
	if err != nil {
		bh.Logger.Infow("can`t commit order",
			"err:", err.Error())
		http.Error(w, "can`t commit order", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte("success"))
	if err != nil {
		bh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// @Summary      Get an information about order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        ORDER_ID    path	integer  true  "ORDER_ID"
// @Success      200  {object}  models.Order
// @Failure      404
// @Failure      500
// @Security ApiKeyAuth
// @Router       /orders/{ORDER_ID} [get]
func (oh *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIdString, ok := vars["ORDER_ID"]
	if !ok {
		oh.Logger.Errorw("no ORDER_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	orderId, err := strconv.Atoi(orderIdString)
	if err != nil {
		oh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	order, err := oh.OrderService.Get(orderId)
	if err != nil {
		oh.Logger.Infow("can`t get order",
			"err:", err.Error())
		http.Error(w, "can`t get order", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(order)

	if err != nil {
		oh.Logger.Errorw("can`t marshal order",
			"err:", err.Error())
		http.Error(w, "can`t make order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		oh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// @Summary      Get all/user`s/my orders
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        USER_ID    query	integer  true  "ORDER_ID"
// @Success      200  {array}  models.Order
// @Failure      403
// @Failure      404
// @Failure      500
// @Security ApiKeyAuth
// @Router       /orders [get]
func (oh *OrderHandler) GetAllMy(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		oh.Logger.Errorw("can`t parse form",
			"err:", err.Error())
		http.Error(w, "can`t parse form", http.StatusBadRequest)
		return
	}

	orderUser := new(models.orderUser)
	err = schema.NewDecoder().Decode(orderUser, r.Form)
	if err != nil {
		oh.Logger.Infow("can`t decode form to struct",
			"err:", err.Error())
		http.Error(w, "can`t decode form to struct", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(orderUser)
	if err != nil {
		oh.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "can`t validate form", http.StatusBadRequest)
		return
	}

	//userID, err := oh.ContextManager.UserIDFromContext(r.Context())
	// if err != nil {
	// 	oh.Logger.Errorw("fail to get id from context",
	// 		"err:", err.Error())
	// 	http.Error(w, "unknown error", http.StatusInternalServerError)
	// 	return
	// }
	userID := orderUser.User_ID

	orders, err := oh.OrderService.GetUsersAll(userID)
	if err != nil {
		oh.Logger.Infow("can`t get order",
			"err:", err.Error())
		http.Error(w, "can`t get order", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(orders)

	if err != nil {
		oh.Logger.Errorw("can`t marshal orders",
			"err:", err.Error())
		http.Error(w, "can`t get orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		oh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// @Summary      Update order
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        ORDER_ID    path	integer  true  "ID of order"
// @Param 		 order_model body models.Order true "updated order"
// @Success      200
// @Failure      400
// @Failure      401
// @Failure      404
// @Failure      500
// @Security ApiKeyAuth
// @Router       /orders/{ORDER_ID} [post]
func (oh *OrderHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderIdString, ok := vars["ORDER_ID"]
	if !ok {
		oh.Logger.Errorw("no ORDER_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	orderId, err := strconv.Atoi(orderIdString)
	if err != nil {
		oh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	order := &models.Order{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		oh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, order)
	if err != nil {
		oh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(order)
	if err != nil {
		oh.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	order.ID = orderId
	*order, err = oh.OrderService.Update(*order)
	if err != nil {
		oh.Logger.Infow("can`t update order",
			"err:", err.Error())
		http.Error(w, "can`t update order", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(order)

	if err != nil {
		oh.Logger.Errorw("can`t marshal order",
			"err:", err.Error())
		http.Error(w, "can`t make order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		oh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}
