# Asset and Debt Management API

Este projeto é uma API para gerenciar ativos e dívidas de usuários, oferecendo operações de criação, atualização, exclusão e consulta. A API inclui autenticação, controle de permissões (administrador e usuário), limites de requisição e uma funcionalidade de cálculo de pontuação do usuário com cache para otimização.

## Estrutura do Projeto

- **`/handlers`**: Contém os controladores das requisições HTTP e as regras de middleware para autenticação e permissões.
- **`/services`**: Contém a lógica de negócio para as operações de ativos e dívidas, além do cálculo de pontuação do usuário.
- **`/models`**: Define os modelos de dados, como `Asset` e `Debt`.
- **`/auth`**: Funções de autenticação e validação de tokens JWT.
- **`main.go`**: Ponto de entrada para inicializar e configurar o servidor da API.
- **`/cache`**: Gerenciamento de cache com Redis para otimização das operações de pontuação.

## Funcionalidades

- **Autenticação (`/login`)**: Autenticação de usuário usando JWT. Retorna um token que deve ser usado para acessar as rotas autenticadas.
- **Gerenciamento de Ativos (`/assets`)**: 
  - Criar, atualizar, deletar e listar ativos.
- **Gerenciamento de Dívidas (`/debts`)**: 
  - Criar, atualizar, deletar e listar dívidas.
- **Cálculo de Pontuação (`/users/{user_id}/score`)**:
  - Calcula e retorna a pontuação de um usuário com base nos ativos e dívidas.
  - Utiliza cache Redis para melhorar o desempenho e reduzir consultas repetidas.
- **Autenticação e Controle de Acesso**: 
  - Autenticação JWT para usuários.
  - Middleware para limitar o acesso de rotas específicas com base no nível de permissão (usuário ou administrador).
- **Limite de Requisições**: 
  - Rate Limiter para evitar sobrecarga de requisições.

## Técnologias

- [Golang](https://golang.org/) >= 1.15
- [MongoDB](https://www.mongodb.com/try/download/community)
- [Redis](https://redis.io/download) para cache de pontuação
- [Gorilla Mux](https://github.com/gorilla/mux) para roteamento
- [X/time/rate](https://pkg.go.dev/golang.org/x/time/rate) para controle de taxa de requisições

## Instalação

1. Para configurar e iniciar a aplicação com Docker, execute os comandos abaixo na raiz do projeto:

   ```bash
   docker compose build
   docker compose up
   ```
2. Já existe um arquivo `.env` configurado com as informações abaixo, mas você pode alterá-las conforme necessário:

    ```
    DATABASE_URL=mongodb://root:root@localhost:27017/desafio?authSource=admin
    PORT=8080
    JWT_SECRET=i0p64HiyQgW4XMqdb281fC5AGtocooV3viA/ET/76OA
    REDIS_PASSWORD=123
    REDIS_ADDR=redis:6379
    ```

3. O banco de dados MongoDB deve ser configurado manualmente, pois houve uma dificuldade em criar automaticamente o banco `desafio` via Docker. O script de inicialização se encontra em `docker/mongodb/init/init-users.js`. Certifique-se de criar o banco `desafio` manualmente no MongoDB.

4. Além disso, para utilizar o sistema é necessário criar dois usuários, um administrador e um usuário comum, com as credenciais abaixo:

   ### Registro dos Usuários

   Use a rota `POST localhost:8080/users/register` para registrar os usuários.

   - **Usuário Administrador**:

     ```json
     {
        "username": "admin",
        "password": "123",
        "is_admin": true
     }
     ```

   - **Usuário Comum**:

     ```json
     {
        "username": "user",
        "password": "123",
        "is_admin": false
     }
     ```

> **Nota:** Esta rota foi criada como auxílio para facilitar os testes do sistema. 

5. Para facilitar os testes, foi incluído um arquivo do Postman chamado `Desafio.postman_collection.json`. Esse arquivo pode ser importado diretamente no Postman, contendo as requisições configuradas para testar as funcionalidades da aplicação.

## Rotas da API

### Autenticação

- **POST `/login`** - Autentica o usuário e retorna um token JWT.
  - **Corpo da requisição**:
    ```json
    {
      "username": "string",
      "password": "string"
    }
    ```
  - **Resposta de sucesso**:
    ```json
    {
      "token": "string"
    }
    ```
### Gerenciamento de Ativos

- **POST `/assets`** - Cria um novo ativo. (Apenas Usuários Autenticados)
  - **Corpo da requisição**:
    ```json
    {
      "Valor": 100,
      "Tipo": "Automóvel",
      "user_id": "123"
    }
    ```

- **PUT `/assets/{id}`** - Atualiza um ativo existente. (Apenas Usuários Autenticados)
  - **Corpo da requisição**:
    ```json
    {
      "Valor": 100,
      "Tipo": "Imóvel"
    }
    ```

- **DELETE `/assets/{id}`** - Deleta um ativo. (Apenas Usuários Autenticados)

- **GET `/assets`** - Lista todos os ativos do usuário autenticado.


### Gerenciamento de Dívidas

- **POST `/debts`** - Cria uma nova dívida para um usuário específico. (Apenas Administradores)
  - **Corpo da requisição**:
    ```json
    {
      "user_id": "123",
      "debt": {
        "Valor": 100,
        "Tipo": "Automóvel"
      }
    }
    ```

- **PUT `/debts/{id}`** - Atualiza uma dívida específica. (Apenas Administradores)
  - **Corpo da requisição**:
    ```json
    {
      "Valor": 100,
      "Tipo": "Imóvel"
    }
    ```

- **DELETE `/debts/{id}`** - Deleta uma dívida específica. (Apenas Administradores)

- **GET `/debts/{user_id}`** - Lista todas as dívidas de um usuário específico. (Apenas Administradores)


### Cálculo de Pontuação

- **GET `/users/{user_id}/score`** - Retorna a pontuação de um usuário com base em seus ativos e dívidas. (Apenas Administradores)
  - A pontuação é calculada com base no valor total de ativos e dívidas do usuário.
  - O resultado é armazenado em cache (Redis) por 15 minutos para otimizar o desempenho.

## Middlewares

- **AuthMiddleware**: Middleware de autenticação JWT para verificar o token do usuário.
- **AdminMiddleware**: Middleware de permissão para rotas que requerem nível de administrador.
- **UserMiddleware**: Middleware de permissão para rotas restritas a usuários não-administradores.
- **RateLimitMiddleware**: Limita o número de requisições a 1000 por segundo para evitar sobrecarga.

## Estrutura de Dados

### Asset

```json
{
    "id": "string",
    "userID": "string",
    "valor": "float",
    "tipo": "string"
}
```

### Debt

```json
{
    "id": "string",
    "userID": "string",
    "valor": "float",
    "tipo": "string",
    "descricao": "string"
}
```

## Score Calculation

O sistema calcula o `score` de um usuário com base nos ativos (`assets`) e dívidas (`debts`). Abaixo está uma explicação detalhada de como o cálculo é realizado.

### Fórmula

A função `CalculateScore` calcula o `score` com os seguintes passos:

1. **Condição Inicial**:  
   - Se o usuário não tiver ativos nem dívidas, o `score` padrão será `500`.

2. **Somatório de Ativos e Dívidas**:
   - O valor total dos ativos (`totalAssetValue`) é a soma dos valores de todos os ativos do usuário.
   - O valor total das dívidas (`totalDebtValue`) é a soma dos valores de todas as dívidas do usuário.

3. **Cálculo do `assetScore`**:
   ```go
   assetScore := math.Log(totalAssetValue + 1) * math.Log(float64(assetCount) + 1)
   ```
   - Calcula o logaritmo natural do valor total dos ativos mais 1 (`totalAssetValue + 1`), multiplicado pelo logaritmo natural da quantidade de ativos (`assetCount + 1`).
   - Essa fórmula pondera positivamente o `score` conforme o valor e a quantidade de ativos aumentam.

4. **Cálculo do `debtScore`**:
   ```go
   debtScore := math.Log(totalDebtValue + 1) * math.Log(float64(debtCount) + 1) * 1.5
   ```
   - Calcula o logaritmo natural do valor total das dívidas mais 1 (`totalDebtValue + 1`), multiplicado pelo logaritmo natural da quantidade de dívidas (`debtCount + 1`).
   - Multiplicado por `1.5` para que as dívidas impactem negativamente o `score` com um peso maior.

5. **Cálculo do `rawScore`**:
   ```go
   rawScore := 500 + (assetScore - debtScore) * 50
   ```
   - Começa com uma pontuação base de `500`.
   - O `assetScore` aumenta o `score`, enquanto o `debtScore` reduz a pontuação.
   - Multiplica a diferença por `50` para definir a escala.

6. **Limite de Pontuação**:
   - O `score` final é limitado entre `0` e `1000`.
   - Se o cálculo resultar em um valor menor que `0`, o `score` final será `0`.
   - Se for maior que `1000`, será ajustado para `1000`.

### Exemplo de Cálculo

```go
assets := []models.Asset{{Valor: 1000.0}, {Valor: 500.0}}
debts := []models.Debt{{Valor: 300.0}}
score := scoreService.CalculateScore(assets, debts)
fmt.Println(score) // Exemplo de saída do cálculo
```

Essas etapas e fórmulas ajudam a fornecer um `score` balanceado, recompensando os ativos e penalizando as dívidas de maneira justa.


### Exemplo de Cálculo com cURL

```bash
# Obter pontuação do usuário
curl -X GET http://localhost:8080/users/{user_id}/score \
     -H "Authorization: Bearer <token>" \
     -H "Content-Type: application/json"
```

## Tratamento de Erros

O projeto inclui tratamento de erros para:

- **Erros de autenticação**: Exemplo, token inválido ou não fornecido.
- **Erros de permissão**: Usuário sem permissão para acessar a rota.
- **Erros de validação**: Dados inválidos ou ausentes.
- **Erros internos do servidor**: Falhas no processamento das operações no banco de dados.

## Contribuições

Contribuições são bem-vindas! Sinta-se à vontade para abrir uma _issue_ ou enviar um _pull request_.

## Licença

Distribuído sob a licença MIT. Consulte o arquivo `LICENSE` para mais informações.
