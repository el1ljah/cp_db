package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/el1ljah/cp_db/internal/models"
	"github.com/el1ljah/cp_db/pkg/logger"

	"github.com/gorilla/mux"
)

type BasketService interface {
	Get(int) (models.Basket, error)
	AddItem(int, int) error
	DecItem(int, int) error
}

type ContextManager interface {
	UserIDFromContext(ctx context.Context) (int, error)
}

type BasketHandler struct {
	BasketService  BasketService
	Logger         logger.Logger
	ContextManager ContextManager
}

// @Summary      Get all items in basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Basket
// @Failure      400
// @Failure      401
// @Failure      403  
// @Failure      500  
// @Security ApiKeyAuth
// @Router       /basket [get]
func (bh *BasketHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, err := bh.ContextManager.UserIDFromContext(r.Context())
	if err != nil {
		bh.Logger.Errorw("fail to get id from context",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	basket, err := bh.BasketService.Get(userID)
	if err != nil {
		bh.Logger.Infow("can`t get basket",
			"err:", err.Error())
		http.Error(w, "can`t get basket", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(basket)

	if err != nil {
		bh.Logger.Errorw("can`t marshal basket",
			"err:", err.Error())
		http.Error(w, "can`t get basket", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		bh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// @Summary      Add item to basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        ITEM_ID    path	integer  true  "ID of item adding to basket"
// @Success      200  
// @Failure      401  
// @Failure      404  
// @Failure      500  
// @Security ApiKeyAuth
// @Router       /basket/{ITEM_ID} [post]
func (bh *BasketHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdString, ok := vars["ITEM_ID"]
	if !ok {
		bh.Logger.Errorw("no ITEM_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		bh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	userID, err := bh.ContextManager.UserIDFromContext(r.Context())
	if err != nil {
		bh.Logger.Errorw("fail to get id from context",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = bh.BasketService.AddItem(itemId, userID)
	if err != nil {
		bh.Logger.Infow("can`t add item to basket (item is not available)",
			"err:", err.Error())
		http.Error(w, "can`t add item to basket (item is not available)", http.StatusBadRequest)
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

// @Summary      Remove item from basket
// @Tags         basket
// @Accept       json
// @Produce      json
// @Param        ITEM_ID    path	integer  true  "ID of item removing from basket"
// @Success      200  
// @Failure      401  
// @Failure      404  
// @Failure      500  
// @Security ApiKeyAuth
// @Router       /basket/{ITEM_ID} [delete]
func (bh *BasketHandler) DecItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemIdString, ok := vars["ITEM_ID"]
	if !ok {
		bh.Logger.Errorw("no ITEM_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	itemId, err := strconv.Atoi(itemIdString)
	if err != nil {
		bh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	userID, err := bh.ContextManager.UserIDFromContext(r.Context())
	if err != nil {
		bh.Logger.Errorw("fail to get id from context",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = bh.BasketService.DecItem(itemId, userID)
	if err != nil {
		bh.Logger.Infow("can`t dec item from basket",
			"err:", err.Error())
		http.Error(w, "can`t dec item from basket", http.StatusBadRequest)
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
