package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ReserveSeatRequest struct {
	FlightId int    `json:"flightId"`
	SeatID   int    `json:"seatID"`
	CPF      string `json:"cpf"`
}

type CancelReservationRequest struct {
	FlightId int    `json:"flightId"`
	SeatID   int    `json:"seatID"`
	CPF      string `json:"cpf"`
}

var (
	baseURL           = "http://localhost:8080"
	numberOfSequences = 3 // Número de sequências completas de operações
	flightId          = 1   // ID do voo de teste
	seatID            = 0  // ID do assento de teste
	testCPF           = "12345678900" // CPF fictício para testes
	wg                sync.WaitGroup
)

func main() {
	fmt.Println("Iniciando testes de concorrência em sequência...")

	// Executa a sequência: Listar -> Reservar -> Cancelar para cada ciclo
	for i := 0; i < numberOfSequences; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			listFlights()
			makeReservation()
			cancelReservation()
		}()
	}

	wg.Wait()
	fmt.Println("Testes de concorrência em sequência concluídos")
}

func listFlights() {
	resp, err := http.Get(fmt.Sprintf("%s/list-flights", baseURL))
	if err != nil {
		fmt.Println("Erro ao listar voos:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Listagem de Voos - Status: %d\n", resp.StatusCode)
}

func makeReservation() {
	reqBody := ReserveSeatRequest{
		FlightId: flightId,
		SeatID:   seatID,
		CPF:      testCPF,
	}

	requestData, _ := json.Marshal(reqBody)
	resp, err := http.Post(fmt.Sprintf("%s/reserve-seat", baseURL), "application/json", bytes.NewBuffer(requestData))
	if err != nil {
		fmt.Println("Erro ao reservar assento:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Reserva - Status: %d\n", resp.StatusCode)
}

func cancelReservation() {
	reqBody := CancelReservationRequest{
		FlightId: flightId,
		SeatID:   seatID,
		CPF:      testCPF,
	}

	requestData, _ := json.Marshal(reqBody)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/cancel-reservation", baseURL), bytes.NewBuffer(requestData))
	if err != nil {
		fmt.Println("Erro ao cancelar reserva:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Erro ao cancelar reserva:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Cancelamento - Status: %d\n", resp.StatusCode)
}
