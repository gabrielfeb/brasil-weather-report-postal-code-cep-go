package usecase

import (
	"context"
	"errors"
	"regexp"

	"clima-cep/internal/domain/entity"
)

var (
	ErrInvalidZipCode  = errors.New("invalid zipcode")
	ErrZipCodeNotFound = errors.New("can not find zipcode")
)

// Define a interface para um provedor de geolocalização por CEP.
type LocationProvider interface {
	GetLocationByCEP(ctx context.Context, cep string) (string, error) // Retorna o nome da cidade
}

// Define a interface para um provedor de clima.
type WeatherProvider interface {
	GetWeatherByCity(ctx context.Context, city string) (float64, error) // Retorna a temperatura em Celsius
}

// DTO (Data Transfer Object) de entrada para o caso de uso.
type GetWeatherByCepInputDTO struct {
	CEP string `json:"cep"`
}

// DTO de saída.
type GetWeatherByCepOutputDTO struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// Implementação do caso de uso.
type GetWeatherByCepUseCase struct {
	LocationProvider LocationProvider
	WeatherProvider  WeatherProvider
}

// Cria uma nova instância do caso de uso.
func NewGetWeatherByCepUseCase(locationProvider LocationProvider, weatherProvider WeatherProvider) *GetWeatherByCepUseCase {
	return &GetWeatherByCepUseCase{
		LocationProvider: locationProvider,
		WeatherProvider:  weatherProvider,
	}
}

// Orquestra a lógica principal da aplicação.
func (uc *GetWeatherByCepUseCase) Execute(ctx context.Context, input GetWeatherByCepInputDTO) (*GetWeatherByCepOutputDTO, error) {
	//Validação do CEP
	if !isValidCEP(input.CEP) {
		return nil, ErrInvalidZipCode
	}

	//Buscar localização usando o provedor injetado
	city, err := uc.LocationProvider.GetLocationByCEP(ctx, input.CEP)
	if err != nil {
		//O adaptador deve retornar ErrZipCodeNotFound se o CEP não for encontrado
		return nil, err
	}

	//Buscar clima usando o provedor injetado
	tempC, err := uc.WeatherProvider.GetWeatherByCity(ctx, city)
	if err != nil {
		return nil, err
	}

	//Criar a entidade de domínio com a lógica de conversão
	weather := entity.NewWeatherOutput(tempC)

	//Mapear para o DTO de saída
	output := &GetWeatherByCepOutputDTO{
		TempC: weather.TempC,
		TempF: weather.TempF,
		TempK: weather.TempK,
	}

	return output, nil
}

// Valida se o CEP tem 8 dígitos numéricos.
func isValidCEP(cep string) bool {
	re := regexp.MustCompile(`^\d{8}$`)
	return re.MatchString(cep)
}
