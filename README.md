# brasil-weather-report-postal-code-cep-go
# API de Clima por CEP

Esta é uma API RESTful em **Go**, que recebe um CEP brasileiro (8 dígitos), identifica a cidade correspondente usando a API ViaCEP, e retorna a temperatura atual em **Celsius**, **Fahrenheit** e **Kelvin**. Foi projetada para execução em contêiner (Docker/Docker‑Compose) e deploy via **Google Cloud Run**, seguindo os princípios de Clean Architecture.

---

## Funcionalidades

- Recebe um CEP brasileiro de 8 dígitos via caminho ou query string.
- Valida o formato do CEP (somente números, exatamente 8 dígitos).
- Consulta ViaCEP para obter cidade e estado.
- Consulta WeatherAPI para obter temperatura atual.
- Converte e retorna as temperaturas em °C, °F, °K.
- Retorna os seguintes status HTTP:
  - 200 OK – sucesso;
  - 422 Unprocessable Entity – formato de CEP inválido;
  - 404 Not Found – CEP não encontrado;
  - 500 Internal Server Error – erro inesperado ou falha em API externa.

---

## Arquitetura

A aplicação segue Clean Architecture, separando:

- Domínio: entidades puras (ex: `WeatherOutput`).
- Casos de uso: lógica principal e orquestração, definindo interfaces.
- Infraestrutura: 
  - Adaptadores para ViaCEP e WeatherAPI,
  - Servidor HTTP (Chi v5),
  - Configuração com Viper.

Também são usados Testify (suite, mock) para testes unitários.

---

## Tecnologias

- Linguagem: Go 1.25+
- Roteador: Chi v5
- Configuração: Viper
- Testes: Testify
- Containerização: Docker, Docker‑Compose
- Deploy: Google Cloud Run
- CI/CD: Google Cloud Build + Artifact Registry

---

## Pré‑requisitos

- Go >= 1.25
- Docker + Docker‑Compose
- (Opcional) Google Cloud SDK — apenas para deploy

---

## Configuração do ambiente

1. Clone o repositório:

   ```bash
   git clone https://github.com/gabrielfeb/brasil-weather-report-postal-code-cep-go.git
   cd brasil-weather-report-postal-code-cep-go
   git checkout feature/develop
   ```

2. Crie arquivo `.env` com sua chave:

   ```env
   WEATHER_API_KEY="SUA_CHAVE_DA_WEATHERAPI_AQUI"
   ```

---

## Executando localmente

Use Docker Compose para build e execução:

```bash
docker-compose up --build
```

A API ficará disponível em:

```
http://localhost:8080
```

---

## Executando testes

Rode:

```bash
go test -v -cover ./...
```

---

## Endpoints

### GET /weather/{cep}

- Parâmetro: `cep` (string, 8 dígitos)
- Resposta 200 (JSON):

  ```json
  {
    "temp_C": 25.0,
    "temp_F": 77.0,
    "temp_K": 298.15
  }
  ```

- 404 Not Found  
  Corpo: `can not find zipcode`

- 422 Unprocessable Entity  
  Corpo: `invalid zipcode`

- 500 Internal Server Error  
  Erro inesperado ocorreu ao consultar APIs externas.

---

## Deploy no Google Cloud Run

1. Execute:

   ```bash
   gcloud run deploy clima-cep-service      --source .      --region us-central1      --allow-unauthenticated      --set-env-vars="WEATHER_API_KEY=$WEATHER_API_KEY"
   ```

2. A URL pública será algo como:

```
https://clima-cep-service-<id>.run.app
```

### Exemplo de chamada

```bash
curl https://clima-cep-service-1075234091446.us-central1.run.app/weather/01001000
```

ou no navegador:

```
https://clima-cep-service-1075234091446.us-central1.run.app/weather/01001000
```

---

## Observações

- As conversões são:
  - °C → °F: °F = °C * 1.8 + 32
  - °C → K: K = °C + 273.15
- Erros de comunicação com ViaCEP ou WeatherAPI resultam em 500 Internal Server Error.
- Ideal para uso em microserviços, projetos escolares ou POCs em nuvem.

---

## Licença

MIT License – veja o arquivo `LICENSE`.

---

## Possível Makefile

```makefile
dev-start:
	docker-compose up --build -d

dev-stop:
	docker-compose stop

dev-down:
	docker-compose down

test:
	go test ./... -cover

prod-start:
	docker-compose -f docker-compose.prod.yml up --build -d

prod-down:
	docker-compose -f docker-compose.prod.yml down
```