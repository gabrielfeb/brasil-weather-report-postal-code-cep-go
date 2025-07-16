package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"clima-cep/internal/usecase"
)

// Estrutura para decodificar a resposta da API ViaCEP.
type ViaCEPResponse struct {
	Localidade string `json:"localidade"`
	Erro       bool   `json:"erro"`
}

// Adaptador que implementa a interface LocationProvider.
type ViaCEPAdapter struct {
	Client *http.Client
}

// Cria um novo adaptador para o ViaCEP.
func NewViaCEPAdapter() *ViaCEPAdapter {
	return &ViaCEPAdapter{Client: &http.Client{}}
}

// Busca a cidade a partir de um CEP.
func (a *ViaCEPAdapter) GetLocationByCEP(ctx context.Context, cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("viacep api returned status: %s", resp.Status)
	}

	var viaCEPResponse ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&viaCEPResponse); err != nil {
		return "", err
	}

	//Retorna 'erro: true' para CEPs não encontrados.
	if viaCEPResponse.Erro {
		return "", usecase.ErrZipCodeNotFound
	}

	return viaCEPResponse.Localidade, nil
}
