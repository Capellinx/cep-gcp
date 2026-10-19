# Clima por CEP

Serviço em Go que recebe um CEP, descobre a cidade (ViaCEP) e retorna a temperatura atual (WeatherAPI) em Celsius, Fahrenheit e Kelvin.

**URL em produção (Cloud Run):** https://cep-gcp-64770695016.southamerica-east1.run.app

## Uso

```bash
curl <URL>/weather/01001000
# {"temp_C":28.5,"temp_F":83.3,"temp_K":301.5}
```

| Cenário | Status | Corpo |
|---|---|---|
| Sucesso | 200 | `{"temp_C":..,"temp_F":..,"temp_K":..}` |
| CEP inválido (não são 8 dígitos) | 422 | `invalid zipcode` |
| CEP inexistente | 404 | `can not find zipcode` |
| Falha em serviço externo | 502 | `upstream error` |

Conversões: `F = C × 1.8 + 32`, `K = C + 273`.

## Rodando localmente

Requer uma chave da [WeatherAPI](https://www.weatherapi.com/).

```bash
echo 'WEATHER_API_KEY=sua-chave' > .env

# testes
docker build --target test .
# ou, com Go instalado: go test ./...

# aplicação
docker build -t cep-gcp .
docker run --rm -p 8080:8080 --env-file .env cep-gcp

curl localhost:8080/weather/01001000
```

## Variáveis de ambiente

| Nome | Obrigatória | Padrão |
|---|---|---|
| `WEATHER_API_KEY` | sim | — |
| `PORT` | não | `8080` |

## Deploy no Cloud Run

```bash
gcloud run deploy cep-gcp --source . --region southamerica-east1 \
  --allow-unauthenticated --set-env-vars WEATHER_API_KEY=sua-chave
```
