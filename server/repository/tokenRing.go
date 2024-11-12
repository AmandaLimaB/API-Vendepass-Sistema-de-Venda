package repository

import (
    "fmt"
    "net/http"
)

// TokenRing struct para manter o token e os servidores vizinhos
type TokenRing struct {
    HasToken     bool
    NextServer   string
    ServerID     string
}

// Inicializa o TokenRing com os dados de ID e vizinho
func NewTokenRing(serverID string, nextServer string, hasToken bool) *TokenRing {
    return &TokenRing{
        ServerID:   serverID,
        NextServer: nextServer,
        HasToken:   hasToken,
    }
}

// Função para enviar o token ao próximo servidor
func (tr *TokenRing) PassToken() {
	if tr.HasToken {
		// Construir a URL do próximo servidor
		url := "http://servidor_b:8081/receive-token"

		// Criar uma nova requisição POST
		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			fmt.Printf("Erro ao criar a requisição: %v\n", err)
			return
		}

		// Fazer a requisição HTTP
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Erro ao enviar o token: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Verificar se o status da resposta é 200 OK
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Token enviado com sucesso para o próximo servidor.")
			tr.HasToken = false
		} else {
			fmt.Printf("Falha ao enviar o token: status %d\n", resp.StatusCode)
		}
	}
}

// Função para receber o token
func (tr *TokenRing) ReceiveToken() {
    tr.HasToken = true
    fmt.Println("Token recebido com sucesso.")
}
