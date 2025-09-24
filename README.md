# Teste de Stress

## Como usar

Buildar e subir os containers:
```cmd
docker compose up --build -d

# logs
docker compose logs -f app
docker compose logs -f server
```

Executar o Stress Test:
```cmd
docker run --rm stress-test run -u http://server:8080 -r 1000 -c 10
```

