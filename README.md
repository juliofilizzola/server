# Server

Este repositório contém um servidor desenvolvido em Go, utilizando sqlc, tern e docker-compose para fins educativos e de demonstração.

## Funcionalidades

- Configuração e inicialização de um servidor básico em Go.
- Suporte a múltiplas rotas e endpoints.
- Manipulação de requisições e respostas.
- Integração com banco de dados usando sqlc e tern.
- Configuração e orquestração de containers com docker-compose.

## Estrutura do Projeto

```
server/
├── cmd/
│   └── tools
|         └── terndotenv
|                  └── main.go         # Arquivo para executar o tern
├── internal/
│   ├── handlers/       # Manipuladores das rotas
│   ├── models/         # Modelos gerados pelo sqlc
│   ├── store/
|          └── pgstore
│              ├── migrations/ # Scripts de migração do banco de dados
│              └── queries/    # Consultas SQL
├── config/
│   └── config.go       # Configurações do servidor
├── docker-compose.yml  # Configuração do docker-compose
├── .gitignore
├── go.mod
├── go.sum
└── README.md
```

## Pré-requisitos

- [Go 1.23](https://golang.org/dl/) instalado na máquina.
- [Docker](https://www.docker.com/get-started) e [Docker Compose](https://docs.docker.com/compose/install/) instalados.

## Instalação

1. Clone o repositório:
    ```sh
    git clone https://github.com/juliofilizzola/server.git
    cd server
    ```

2. Instale as dependências do Go:
    ```sh
    go mod tidy
    ```

3. Gere os arquivos de consulta sqlc:
    ```sh
    sqlc generate
    ```

## Uso

Para iniciar o servidor usando Docker Compose, execute:

```sh
docker-compose up --build
```

O servidor estará rodando em `http://localhost:3000`.

## Migrações do Banco de Dados

Para aplicar migrações do banco de dados usando tern, execute:

```sh
tern migrate
```

## Contribuição

Se você quiser contribuir com o projeto:

1. Faça um fork do repositório.
2. Crie uma branch para a sua feature (`git checkout -b feature/nova-feature`).
3. Commit suas mudanças (`git commit -am 'Adiciona nova feature'`).
4. Faça o push para a branch (`git push origin feature/nova-feature`).
5. Crie um novo Pull Request.

## Licença

Este projeto está licenciado sob a Licença MIT - veja o arquivo [LICENSE](LICENSE) para mais detalhes.

---

Sinta-se à vontade para ajustar conforme necessário!