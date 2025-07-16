package main

import (
	"log"
	"net/http"

	"clima-cep/internal/infra/adapter"
	"clima-cep/internal/infra/web/handler"
	"clima-cep/internal/infra/web/server"
	"clima-cep/internal/usecase"

	"github.com/spf13/viper"
)

func main() {
	//Carrega configurações usando Viper
	viper.SetConfigName("configs")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./configs")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	apiKey := viper.GetString("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY must be set")
	}
	serverPort := viper.GetString("SERVER_PORT")

	//Injeção de Dependências
	locationProvider := adapter.NewViaCEPAdapter()
	weatherProvider := adapter.NewWeatherAPIAdapter(apiKey)
	getWeatherUseCase := usecase.NewGetWeatherByCepUseCase(locationProvider, weatherProvider)
	weatherHandler := handler.NewWeatherHandler(getWeatherUseCase)

	//Configura o servidor e as rotas
	router := server.SetupRoutes(weatherHandler)

	//Inicia o servidor
	log.Printf("Server running on port %s", serverPort)
	if err := http.ListenAndServe(serverPort, router); err != nil {
		log.Fatalf("could not start server: %s\n", err)
	}
}
