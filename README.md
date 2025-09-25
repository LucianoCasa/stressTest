# Teste de Stress

## Teste com o docker

Sobe o server para simular um servidor web simples enviando status aleatório e o app fica aguardando o comando manual para iniciar o teste de stress.
```cmd
docker compose up --build -d
```

Executar o Stress Test:
```cmd
docker compose run --rm app run -u http://server:8080 -r 1000 -c 10
```

## Teste com site externo

Teste em cima do google.com
```cmd
task s
```
