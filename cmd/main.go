package main

import (
	"fmt"
	"net/http"

	basketDel "github.com/el1ljah/cp_db/internal/basket/delivery"
	basketRepo "github.com/el1ljah/cp_db/internal/basket/repo"
	basketServ "github.com/el1ljah/cp_db/internal/basket/service"
	brandDel "github.com/el1ljah/cp_db/internal/brand/delivery"
	brandRepo "github.com/el1ljah/cp_db/internal/brand/repo"
	brandServ "github.com/el1ljah/cp_db/internal/brand/service"
	itemDel "github.com/el1ljah/cp_db/internal/item/delivery"
	itemRepo "github.com/el1ljah/cp_db/internal/item/repo"
	itemServ "github.com/el1ljah/cp_db/internal/item/service"
	orderDel "github.com/el1ljah/cp_db/internal/order/delivery"
	orderRepo "github.com/el1ljah/cp_db/internal/order/repo"
	orderServ "github.com/el1ljah/cp_db/internal/order/service"
	userDel "github.com/el1ljah/cp_db/internal/user/delivery"
	userRepo "github.com/el1ljah/cp_db/internal/user/repo"
	userServ "github.com/el1ljah/cp_db/internal/user/service"
	"github.com/el1ljah/cp_db/pkg/context"
	"github.com/el1ljah/cp_db/pkg/middleware"
	"github.com/el1ljah/cp_db/pkg/session"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "github.com/el1ljah/cp_db/docs" 
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

const port = ":8080"

// @title           Clothes store 👚
// @version         1.337
// @description     WEB лабы. Как я хочу спать.....

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

// @tag.name authentication
// @tag.name items
// @tag.name brands
// @tag.name basket
func main() {
	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	params := "user=postgres dbname=clothshop password=postgres host=localhost port=5432 sslmode=disable"
	db, err := sqlx.Connect("postgres", params)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	sessionManager := session.JWTSessionsManager{}
	contextManager := context.ContextManager{}

	authManager := middleware.AuthManager{
		SessionManager: sessionManager,
		Logger:         logger,
		ContextManager: contextManager,
	}

	userHandler := userDel.UserHandler{
		Logger:   logger,
		Sessions: sessionManager,
		UserService: userServ.UserService{
			UserRepo: &userRepo.PgUserRepo{
				Logger: logger,
				DB:     db,
			},
			Logger: logger,
		},
	}

	brandHandler := brandDel.BrandHandler{
		Logger: logger,
		BrandService: brandServ.BrandService{
			BrandRepo: &brandRepo.PgBrandRepo{
				Logger: logger,
				DB:     db,
			},
			Logger: logger,
		},
	}

	itemHandler := itemDel.ItemHandler{
		Logger: logger,
		ItemService: itemServ.ItemService{
			ItemRepo: &itemRepo.PgItemRepo{
				Logger: logger,
				DB:     db,
			},
			Logger: logger,
		},
	}

	basketHandler := basketDel.BasketHandler{
		ContextManager: &contextManager,
		Logger:         logger,
		BasketService: basketServ.BasketService{
			BasketRepo: &basketRepo.PgBasketRepo{
				Logger: logger,
				DB:     db,
			},
			Logger: logger,
		},
	}

	orderHandler := orderDel.OrderHandler{
		ContextManager: &contextManager,
		Logger:         logger,
		OrderService: orderServ.OrderService{
			OrderRepo: &orderRepo.PgOrderRepo{
				Logger: logger,
				DB:     db,
			},
			Logger: logger,
		},
	}

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	r.HandleFunc("/register", userHandler.Register).Methods("POST")
	r.HandleFunc("/login", userHandler.Login).Methods("POST")

	r.HandleFunc("/brands/{BRAND_ID:[0-9]+}", http.HandlerFunc(brandHandler.Get)).Methods("GET")
	r.Handle("/brands", authManager.Auth(http.HandlerFunc(brandHandler.Create), "admin")).Methods("PUT")
	r.Handle("/brands/{BRAND_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(brandHandler.Update), "admin")).Methods("POST")
	r.Handle("/brands/{BRAND_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(brandHandler.Delete), "admin")).Methods("DELETE")

	r.HandleFunc("/items/{ITEM_ID:[0-9]+}", http.HandlerFunc(itemHandler.Get)).Methods("GET")
	r.Handle("/items", authManager.Auth(http.HandlerFunc(itemHandler.Create), "admin")).Methods("PUT")
	r.Handle("/items/{ITEM_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(itemHandler.Update), "admin")).Methods("POST")
	r.Handle("/items/{ITEM_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(itemHandler.Patch), "admin")).Methods("PATCH")
	r.Handle("/items/{ITEM_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(itemHandler.Delete), "admin")).Methods("DELETE")
	r.HandleFunc("/items", http.HandlerFunc(itemHandler.GetAll)).Methods("GET")

	r.Handle("/basket", authManager.Auth(http.HandlerFunc(basketHandler.Get), "user", "admin")).Methods("GET")
	r.Handle("/basket/{ITEM_ID}", authManager.Auth(http.HandlerFunc(basketHandler.AddItem), "user", "admin")).Methods("POST")
	r.Handle("/basket/{ITEM_ID}", authManager.Auth(http.HandlerFunc(basketHandler.DecItem), "user", "admin")).Methods("DELETE")

	r.Handle("/orders", authManager.Auth(http.HandlerFunc(orderHandler.Commit), "user", "admin")).Methods("POST")
	r.Handle("/orders/{ORDER_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(orderHandler.Get), "admin")).Methods("GET")
	r.Handle("/orders/{ORDER_ID:[0-9]+}", authManager.Auth(http.HandlerFunc(orderHandler.Update), "admin")).Methods("POST")
	//r.Handle("/orders/my", authManager.Auth(http.HandlerFunc(orderHandler.GetAllMy), "user", "admin")).Methods("GET")
	r.Handle("/orders", authManager.Auth(http.HandlerFunc(orderHandler.GetAll), "admin")).Methods("GET")


	mux := middleware.AccessLog(logger, r)
	mux = middleware.Panic(logger, mux)

	logger.Infow("starting server",
		"type", "START",
		"port", port,
	)

	logger.Errorln(http.ListenAndServe(port, mux))

	err = zapLogger.Sync()
	if err != nil {
		fmt.Println(err)
	}
}
