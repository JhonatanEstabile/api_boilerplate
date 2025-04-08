# 🚀 API Boilerplate com Go + Gin + sqlx + Generics

Este é um projeto base para criar APIs RESTful em Go usando `gin` e `sqlx`, com suporte a CRUD genérico, testes automatizados e um **gerador automático de domains** (modelos + rotas).

---

## 📦 Tecnologias usadas

- [Go](https://golang.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [sqlx](https://github.com/jmoiron/sqlx)
- [sqlite3](https://github.com/mattn/go-sqlite3)
- [Testify](https://github.com/stretchr/testify)
- [sqlmock](https://github.com/DATA-DOG/go-sqlmock)

---

## ▶️ Como rodar o projeto

1. **Clone o repositório e instale as dependências:**

```bash
go mod tidy
```

2. **Rode o projeto:**

```bash
go run main.go
```

3. **API disponível em:**

```
http://localhost:8080
```

---

## 🧪 Rodar os testes com cobertura

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

---

## ⚙️ Gerador de Domains

Você pode criar models automaticamente com rota e registro no sistema usando o script:

### 🔧 Como usar

```bash
go run cmd/create_domain/main.go <nome> [Campo:Tipo Campo:Tipo ...]
```

### 🔁 Exemplo:

```bash
go run cmd/create_domain/main.go car Name:string Brand:string Year:int
```

### 🧬 Isso irá:

1. Criar o arquivo `model/car.go` com a struct `Car`:

```go
type Car struct {
    ID    string  `db:"id" json:"id"`
    Name  string `db:"name" json:"name"`
    Brand string `db:"brand" json:"brand"`
    Year  int    `db:"year" json:"year"`
    CreatedAt   string `json:"created_at" db:"created_at"`
	UpdatedAt   string `json:"updated_at" db:"updated_at"`
}
```

2. Adicionar automaticamente a linha no `util/registry.go`:

```go
RegisterGenericResource[model.Car](r, db, "car", model.CarFields)
```

## Estrutura Padrão dos Models

Os models no projeto seguem a seguinte estrutura padrão:

- `id`: Identificador único do tipo ULID.
- `created_at`: Timestamp indicando quando o registro foi criado.
- `updated_at`: Timestamp indicando quando o registro foi atualizado pela última vez.

## Filtragem de Dados

A API suporta filtragem de dados através de parâmetros de consulta (query parameters). Os filtros disponíveis são:

- `eql`: Filtra registros onde o campo é igual ao valor especificado.
- `lik`: Filtra registros onde o campo corresponde ao padrão especificado utilizando `LIKE`.

### Exemplo de Uso

Para filtrar produtos com o nome exatamente igual a "Produto1":

```
GET /products?name=eql,Produto1
```

Para filtrar produtos cujo nome contém "Produto":

```
GET /products?name=lik,Produto
```

---

## 📁 Estrutura

```
.
├── model/               # Models gerados
├── controller/          # Controller genérico
├── service/             # Service genérico
├── repository/          # Repository genérico
├── util/registry.go     # Registro central dos domains
├── db/                  # Conexão com banco de dados
├── main.go              # Entrada principal
├── generate_domain.go   # Gerador de domínio automático
```

---

## 🛠 Exemplo de endpoint gerado

Após rodar `go run generate_domain.go book Title:string Pages:int`, você poderá acessar:

- `GET /book`
- `POST /book`
- `GET /book/:id`
- `PUT /book/:id`
- `DELETE /book/:id`

Tudo pronto, sem escrever código manual.

---

## 🧠 Dúvidas ou sugestões?

Abra uma issue ou contribua com um PR. 😉

---

Desenvolvido com ❤️ em Go.
