package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock para a interface LocationProvider.
type MockLocationProvider struct {
	mock.Mock
}

func (m *MockLocationProvider) GetLocationByCEP(ctx context.Context, cep string) (string, error) {
	args := m.Called(ctx, cep)
	return args.String(0), args.Error(1)
}

// Mock para a interface WeatherProvider.
type MockWeatherProvider struct {
	mock.Mock
}

func (m *MockWeatherProvider) GetWeatherByCity(ctx context.Context, city string) (float64, error) {
	args := m.Called(ctx, city)
	return args.Get(0).(float64), args.Error(1)
}

func TestGetWeatherByCepUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("should return weather when a valid cep is provided", func(t *testing.T) {
		mockLocation := new(MockLocationProvider)
		mockWeather := new(MockWeatherProvider)
		uc := NewGetWeatherByCepUseCase(mockLocation, mockWeather)

		input := GetWeatherByCepInputDTO{CEP: "01001000"}
		expectedCity := "São Paulo"
		expectedTempC := 25.0

		mockLocation.On("GetLocationByCEP", ctx, input.CEP).Return(expectedCity, nil).Once()
		mockWeather.On("GetWeatherByCity", ctx, expectedCity).Return(expectedTempC, nil).Once()

		output, err := uc.Execute(ctx, input)

		assert.NoError(t, err)
		assert.NotNil(t, output)
		assert.Equal(t, 25.0, output.TempC)
		assert.Equal(t, 77.0, output.TempF)   // 25 * 1.8 + 32
		assert.Equal(t, 298.15, output.TempK) // 25 + 273.15

		//Verifica se os mocks foram chamados como esperado
		mockLocation.AssertExpectations(t)
		mockWeather.AssertExpectations(t)
	})

	t.Run("should return ErrInvalidZipCode when cep is not valid", func(t *testing.T) {
		mockLocation := new(MockLocationProvider)
		mockWeather := new(MockWeatherProvider)
		uc := NewGetWeatherByCepUseCase(mockLocation, mockWeather)
		input := GetWeatherByCepInputDTO{CEP: "12345"}

		output, err := uc.Execute(ctx, input)

		assert.Nil(t, output)
		assert.Error(t, err)
		assert.Equal(t, ErrInvalidZipCode, err)
		mockLocation.AssertNotCalled(t, "GetLocationByCEP", ctx, mock.Anything)
	})

	t.Run("should return ErrZipCodeNotFound when location provider fails to find cep", func(t *testing.T) {
		mockLocation := new(MockLocationProvider)
		mockWeather := new(MockWeatherProvider)
		uc := NewGetWeatherByCepUseCase(mockLocation, mockWeather)
		input := GetWeatherByCepInputDTO{CEP: "99999999"}

		//Configura o mock para retornar o erro específico
		mockLocation.On("GetLocationByCEP", ctx, input.CEP).Return("", ErrZipCodeNotFound).Once()

		output, err := uc.Execute(ctx, input)

		assert.Nil(t, output)
		assert.Error(t, err)
		assert.Equal(t, ErrZipCodeNotFound, err)
		mockLocation.AssertExpectations(t)
		mockWeather.AssertNotCalled(t, "GetWeatherByCity", ctx, mock.Anything)
	})

	t.Run("should return an error when weather provider fails", func(t *testing.T) {
		mockLocation := new(MockLocationProvider)
		mockWeather := new(MockWeatherProvider)
		uc := NewGetWeatherByCepUseCase(mockLocation, mockWeather)

		input := GetWeatherByCepInputDTO{CEP: "01001000"}
		expectedCity := "São Paulo"
		providerError := errors.New("weather service unavailable")

		mockLocation.On("GetLocationByCEP", ctx, input.CEP).Return(expectedCity, nil).Once()
		mockWeather.On("GetWeatherByCity", ctx, expectedCity).Return(0.0, providerError).Once()

		output, err := uc.Execute(ctx, input)

		assert.Nil(t, output)
		assert.Error(t, err)
		assert.Equal(t, providerError, err)
		mockLocation.AssertExpectations(t)
		mockWeather.AssertExpectations(t)
	})
}
