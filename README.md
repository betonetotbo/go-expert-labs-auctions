# Abertura e fechamento do Leilão - Go Routines

Pré-requisitos:
* Ter docker e docker compose na máquina

Para executar siga os passos:
* Na raiz deste repositório execute: `docker compose up -d`
  * Este comando irá fazer o *build* da aplicação e executá-la localmente
  * Serão provisionados:
    * A aplicação de auctions na porta http `8080`
    * O MongoDB na porta `27017`
    * O [mongo-express](https://github.com/mongo-express/mongo-express) na porta http `8081` (interface gráfica para explorar o mongodb)
  * As configurações da aplication de auctions estão em `cmd/auction/.env` 
* Para executar os testes existem [scripts HTTP (Intellij)](https://www.jetbrains.com/help/idea/http-client-in-product-code-editor.html) na pasta `scripts`
  * O script `test.http` tem testes para:
    * Criação de um *auction*
    * Criação de um *bid* para o *auction* criado anteriormente
    * Consulta do *bid* vencedor do *auction* criado anteriormente
    * Consulta do *auction* criado anteriormente
    * Listagem de *auctions*

### Executando os testes com `curl`

#### Criação de um auction

```shell
AUCTION=$(curl -v http://localhost:8080/auction -d '{"product_name":"produto 4","category":"categoria","description":"desicricao descricao","condition":1}' | jq -rM .id)
```

#### Criação de bid para o auction criado com o comando anterior

```shell
curl -v http://localhost:8080/bid -d "{\"user_id\":\"00000000-0000-0000-0000-000000000000\",\"auction_id\":\"$AUCTION\",\"amount\":1.53}"
```

> Execute mais de uma vez esse comando para criar vários BIDs

#### Consulta do bid vencedor do auction criado com o comando anterior

```shell
curl "http://localhost:8080/auction/winner/$AUCTION"
```

#### Consulta do ultimo auction criado

```shell
curl "http://localhost:8080/auction/$AUCTION"
```

#### Listagem dos auctions

```shell
curl "http://localhost:8080/auction?status=0"
```

### O que foi feito?

* Criação do método `autocloseAuction` em `internal/infra/database/auction/create_auction.go`
* Alguns refactors para possibilitar reuso de código
* Criação do teste `internal/infra/database/auction/create_auction_test.go` para validar o **auto fechamento** de um leilão
* Incluído no `docker-compose.yml` o mongo-express (interface para explorar o mongoDB)
* Scripts HTTP em `scripts/` e esta documentação