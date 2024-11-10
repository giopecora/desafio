package routes

import (
	"main/cache"
	"main/handlers"
	"main/services"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupRoutes(router *mux.Router, client *mongo.Client, redisClient cache.RedisClient) {
	assetService := services.NewAssetService(client)
	assetHandler := handlers.NewAssetHandler(assetService)

	debtService := services.NewDebtService(client)
	debtHandler := handlers.NewDebtHandler(debtService)

	userService := services.NewUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	authHandler := handlers.NewAuthHandler(client)

	scoreService := services.NewScoreService(assetService, debtService)
	scoreHandler := handlers.NewScoreHandler(scoreService, redisClient)

	router.HandleFunc("/assets", handlers.UserMiddleware(assetHandler.CreateAssetHandler)).Methods("POST")
	router.HandleFunc("/assets/{id}", handlers.UserMiddleware(assetHandler.UpdateAssetHandler)).Methods("PUT")
	router.HandleFunc("/assets/{id}", handlers.UserMiddleware(assetHandler.DeleteAssetHandler)).Methods("DELETE")
	router.HandleFunc("/assets", handlers.UserMiddleware(assetHandler.GetAssetsHandler)).Methods("GET")

	router.HandleFunc("/debts", handlers.AdminMiddleware(debtHandler.CreateDebtHandler)).Methods("POST")
	router.HandleFunc("/debts/{id}", handlers.AdminMiddleware(debtHandler.UpdateDebtHandler)).Methods("PUT")
	router.HandleFunc("/debts/{id}", handlers.AdminMiddleware(debtHandler.DeleteDebtHandler)).Methods("DELETE")
	router.HandleFunc("/debts/{user_id}", handlers.AdminMiddleware(debtHandler.GetDebtsHandler)).Methods("GET")

	router.HandleFunc("/users/{user_id}/score", func(w http.ResponseWriter, r *http.Request) {
		handlers.RateLimitMiddleware(
			handlers.AdminMiddleware(scoreHandler.GetUserScoreHandler),
		).ServeHTTP(w, r)
	}).Methods("GET")

	router.HandleFunc("/users/register", userHandler.RegisterUser).Methods("POST")
	router.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
}
