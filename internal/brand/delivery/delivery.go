package delivery

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/el1ljah/cp_db/internal/models"
	"github.com/el1ljah/cp_db/pkg/logger"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type BrandService interface {
	Create(models.Brand) (int, error)
	Get(int) (models.Brand, error)
	Update(models.Brand) (models.Brand, error)
	Delete(int) error
}

type BrandHandler struct {
	BrandService BrandService
	Logger       logger.Logger
}

// @Summary      Get an information about one brand
// @Tags         brands
// @Accept       json
// @Produce      json
// @Param        BRAND_ID    path	integer  true  "ID of brand"
// @Success      200  {object}  models.Brand
// @Failure      404  
// @Failure      500  
// @Router       /brands/{BRAND_ID} [get]
func (bh *BrandHandler) Get(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brandIdString, ok := vars["BRAND_ID"]
	if !ok {
		bh.Logger.Errorw("no BRAND_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	brandId, err := strconv.Atoi(brandIdString)
	if err != nil {
		bh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	brand, err := bh.BrandService.Get(brandId)
	if err != nil {
		bh.Logger.Infow("can`t get brand",
			"err:", err.Error())
		http.Error(w, "can`t get brand", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(brand)

	if err != nil {
		bh.Logger.Errorw("can`t marshal brand",
			"err:", err.Error())
		http.Error(w, "can`t make brand", http.StatusInternalServerError)
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

// @Summary      Add new brand
// @Tags         brands
// @Accept       json
// @Produce      json
// @Param data body models.Brand true "new brand"
// @Success      200  
// @Failure      400
// @Failure      401
// @Failure      404  
// @Failure      500  
// @Security ApiKeyAuth
// @Router       /brands [put]
func (bh *BrandHandler) Create(w http.ResponseWriter, r *http.Request) {
	brand := &models.Brand{}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		bh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, brand)
	if err != nil {
		bh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(brand)
	if err != nil {
		bh.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	brand.ID, err = bh.BrandService.Create(*brand)
	if err != nil {
		bh.Logger.Infow("can`t create brand",
			"err:", err.Error())
		http.Error(w, "can`t create brand", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(brand)

	if err != nil {
		bh.Logger.Errorw("can`t marshal brand",
			"err:", err.Error())
		http.Error(w, "can`t make brand", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	_, err = w.Write(resp)
	if err != nil {
		bh.Logger.Errorw("can`t write response",
			"err:", err.Error())
		http.Error(w, "can`t write response", http.StatusInternalServerError)
		return
	}
}

// @Summary      Update brand
// @Tags         brands
// @Accept       json
// @Produce      json
// @Param        BRAND_ID    path	integer  true  "ID of update brand"
// @Param 		 data body models.Brand true "updated brand"
// @Success      200  
// @Failure      400
// @Failure      401
// @Failure      404  
// @Failure      500  
// @Security ApiKeyAuth
// @Router       /brands/{BRAND_ID} [post]
func (bh *BrandHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brandIdString, ok := vars["BRAND_ID"]
	if !ok {
		bh.Logger.Errorw("no BRAND_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	brandId, err := strconv.Atoi(brandIdString)
	if err != nil {
		bh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	brand := &models.Brand{}
	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		bh.Logger.Errorw("can`t read body of request",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, brand)
	if err != nil {
		bh.Logger.Infow("can`t unmarshal form",
			"err:", err.Error())
		http.Error(w, "bad  data", http.StatusBadRequest)
		return
	}

	_, err = govalidator.ValidateStruct(brand)
	if err != nil {
		bh.Logger.Infow("can`t validate form",
			"err:", err.Error())
		http.Error(w, "bad data", http.StatusBadRequest)
		return
	}

	brand.ID = brandId
	*brand, err = bh.BrandService.Update(*brand)
	if err != nil {
		bh.Logger.Infow("can`t update brand",
			"err:", err.Error())
		http.Error(w, "can`t update brand", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(brand)

	if err != nil {
		bh.Logger.Errorw("can`t marshal brand",
			"err:", err.Error())
		http.Error(w, "can`t make brand", http.StatusInternalServerError)
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

// @Summary      Delete brand
// @Tags         brands
// @Accept       json
// @Produce      json
// @Param        BRAND_ID    path	integer  true  "ID of deleted brand"
// @Success      200 
// @Failure      401
// @Failure      404  
// @Failure      500  
// @Security ApiKeyAuth
// @Router       /brands/{BRAND_ID} [delete]
func (bh *BrandHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	brandIdString, ok := vars["BRAND_ID"]
	if !ok {
		bh.Logger.Errorw("no BRAND_ID var")
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	brandId, err := strconv.Atoi(brandIdString)
	if err != nil {
		bh.Logger.Errorw("fail to convert id to int",
			"err:", err.Error())
		http.Error(w, "unknown error", http.StatusInternalServerError)
		return
	}

	err = bh.BrandService.Delete(brandId)
	if err != nil {
		bh.Logger.Infow("can`t delete brand",
			"err:", err.Error())
		http.Error(w, "can`t delete brand", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}
