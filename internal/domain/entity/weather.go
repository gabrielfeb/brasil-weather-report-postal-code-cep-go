package entity

//Entidade de domínio que representa a saída final da temperatura.
type WeatherOutput struct {
	TempC float64 `json:"temp_C"` // Temperatura em Celsius
	TempF float64 `json:"temp_F"` // Temperatura em Fahrenheit
	TempK float64 `json:"temp_K"` // Temperatura em Kelvin
}

//Cria uma nova instância de WeatherOutput com as temperaturas convertidas.
func NewWeatherOutput(tempC float64) *WeatherOutput {
	return &WeatherOutput{
		TempC: tempC,
		TempF: tempC*1.8 + 32, // F = C * 1.8 + 32
		TempK: tempC + 273.15, // K = C + 273.15 (usando 273.15 para maior precisão)
	}
}
