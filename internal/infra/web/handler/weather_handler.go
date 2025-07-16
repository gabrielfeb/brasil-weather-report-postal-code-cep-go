package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"clima-cep/internal/usecase"

	"github.com/go-chi/chi/v5"
)

// Lida com as requisições HTTP para o clima.
type WeatherHandler struct {
	GetWeatherUseCase *usecase.GetWeatherByCepUseCase
}

// Cria um novo handler de clima.
func NewWeatherHandler(getWeatherUseCase *usecase.GetWeatherByCepUseCase) *WeatherHandler {
	return &WeatherHandler{
		GetWeatherUseCase: getWeatherUseCase,
	}
}

// Método do handler que efetivamente processa a requisição.
// @Summary      Get weather by CEP
// @Description  Get the current temperature in C, F, and K for a given Brazilian CEP
// @Tags         Weather
// @Accept       json
// @Produce      json
// @Param        cep   path      string  true  "CEP (8 digits)"
// @Success      200   {object}  usecase.GetWeatherByCepOutputDTO
// @Failure      404   {string}  string "can not find zipcode"
// @Failure      422   {string}  string "invalid zipcode"
// @Failure      500   {string}  string "internal server error"
// @Router       /weather/{cep} [get]
func (h *WeatherHandler) GetWeather(w http.ResponseWriter, r *http.Request) {
	cep := chi.URLParam(r, "cep")
	if cep == "" {
		http.Error(w, "zipcode is required", http.StatusBadRequest)
		return
	}

	input := usecase.GetWeatherByCepInputDTO{CEP: cep}
	output, err := h.GetWeatherUseCase.Execute(r.Context(), input)

	if err != nil {
		//Mapeando erros do caso de uso para status codes HTTP
		if errors.Is(err, usecase.ErrInvalidZipCode) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity) // 422
			return
		}
		if errors.Is(err, usecase.ErrZipCodeNotFound) {
			http.Error(w, "can not find zipcode", http.StatusNotFound) // 404
			return
		}
		//Para qualquer outro erro, retornamos um erro interno do servidor
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError) // 500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
