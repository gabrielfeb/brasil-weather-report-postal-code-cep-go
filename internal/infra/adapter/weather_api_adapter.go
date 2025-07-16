package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Estrutura para decodificar a resposta da WeatherAPI.
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// Adaptador que implementa a interface WeatherProvider.
type WeatherAPIAdapter struct {
	Client *http.Client
	APIKey string
}

// Cria um novo adaptador para a WeatherAPI.
func NewWeatherAPIAdapter(apiKey string) *WeatherAPIAdapter {
	return &WeatherAPIAdapter{
		Client: &http.Client{},
		APIKey: apiKey,
	}
}

// Busca a temperatura em Celsius para uma cidade.
func (a *WeatherAPIAdapter) GetWeatherByCity(ctx context.Context, city string) (float64, error) {
	// Codifica a cidade para garantir que caracteres especiais sejam tratados corretamente.
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", a.APIKey, encodedCity)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Retorna um erro se a resposta não for OK.
		return 0, fmt.Errorf("weatherapi returned status: %s", resp.Status)
	}

	var weatherResponse WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&weatherResponse); err != nil {
		return 0, err
	}

	return weatherResponse.Current.TempC, nil
}
