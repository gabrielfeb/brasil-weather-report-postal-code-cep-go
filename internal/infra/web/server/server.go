package server

import (
	"clima-cep/internal/infra/web/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Configura e retorna o roteador HTTP com todas as rotas da aplicação.
func SetupRoutes(weatherHandler *handler.WeatherHandler) *chi.Mux {
	r := chi.NewRouter()

	// Middlewares
	r.Use(middleware.Logger)    // Loga as requisições
	r.Use(middleware.Recoverer) // Se a aplicação entrar em pânico, recupera e retorna 500

	// Rotas
	r.Get("/weather/{cep}", weatherHandler.GetWeather)

	return r
}
