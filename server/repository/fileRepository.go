// fileRepository.go
package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/AmandaLimaB/API-Vendepass-Sistema-de-Venda/tree/main/server/models"
)

const dataPath = "data/routes.json"

// LoadFlights carrega os voos do arquivo JSON.
func LoadFlights() ([]models.Flight, error) {
	var flights []models.Flight
	file, err := os.Open(filepath.Clean(dataPath))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&flights)
	if err != nil {
		return nil, err
	}

	return flights, nil
}

// SaveFlights salva os voos no arquivo JSON.
func SaveFlights(flights []models.Flight) error {
	file, err := os.Create(filepath.Clean(dataPath))
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(flights)
	if err != nil {
		return err
	}

	fmt.Println("Dados salvos com sucesso!")
	return nil
}

const clientsFile = "data/clients.json"

func SaveClient(client models.Client) error {
	clients, err := LoadClients()
	if err != nil {
		return err
	}

	// Verificar se o CPF já existe
	for _, c := range clients {
		if c.CPF == client.CPF {
			return errors.New("cliente com este CPF já existe")
		}
	}

	clients = append(clients, client)
	return saveClientsToFile(clients)
}

func LoadClients() ([]models.Client, error) {
	file, err := os.Open(clientsFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.Client{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var clients []models.Client
	if err := json.NewDecoder(file).Decode(&clients); err != nil {
		return nil, err
	}

	return clients, nil
}

func saveClientsToFile(clients []models.Client) error {
	file, err := os.Create(clientsFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(clients)
}

func CancelReservation(flightId int, seatID int, cpf string, externalCompanies []string) error {
	// Carregar os voos locais
	flights, err := LoadFlights()
	if err != nil {
		return err
	}

	// Tentar cancelar a reserva localmente
	var flight *models.Flight
	for i := range flights {
		if flights[i].FlightId == flightId {
			flight = &flights[i]
			break
		}
	}

	// Se o voo foi encontrado localmente
	if flight != nil {
		// Verificar se o assento está reservado
		if seatID < 0 || seatID >= len(flight.Seats) || !flight.Seats[seatID].IsReserved || flight.Seats[seatID].CustomerID != cpf {
			return errors.New("cancelamento não realizado") // Mensagem genérica
		}
		// Cancelar a reserva
		flight.Seats[seatID].IsReserved = false
		flight.Seats[seatID].CustomerID = ""

		// Salvar a atualização local
		if err := SaveFlights(flights); err != nil {
			return err
		}
		return nil // Cancelamento local bem-sucedido
	}

	// Se o voo não foi encontrado localmente, tentar nas companhias externas
	reqBody, _ := json.Marshal(models.CancelReservationRequest{
		FlightId: flightId,
		SeatID:   seatID,
		CPF:      cpf,
	})

	// Configurar cliente HTTP com timeout para evitar espera longa
	client := http.Client{Timeout: 5 * time.Second}
	for _, companyURL := range externalCompanies {
		resp, err := client.Post(fmt.Sprintf("%s/cancel-reservation", companyURL), "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Printf("Erro ao tentar cancelamento na companhia %s: %v\n", companyURL, err)
			continue // Ignorar erro e tentar a próxima companhia
		}
		defer resp.Body.Close()

		// Verificar resposta da companhia externa
		if resp.StatusCode == http.StatusOK {
			return nil // Cancelamento realizado com sucesso em uma companhia externa
		} else if resp.StatusCode == http.StatusConflict {
			// Se o assento não estiver reservado, retornar a mensagem genérica
			continue
		}
	}

	// Se não foi possível cancelar em nenhuma companhia
	return errors.New("cancelamento não realizado") // Mensagem genérica
}

func LoadConfig() models.Config {
    file, err := os.Open("data/config.json")
    if err != nil {
        log.Fatalf("Erro ao carregar configuração: %v", err)
    }
    defer file.Close()

    var config models.Config
    if err := json.NewDecoder(file).Decode(&config); err != nil {
        log.Fatalf("Erro ao decodificar configuração: %v", err)
    }
    return config
}


func ReserveSeat(flightId int, seatID int, cpf string, externalCompanies []string) error {
	flights, err := LoadFlights()
	if err != nil {
		return err
	}

	// Determina a companhia pelo prefixo do FlightId
	companyIndex := flightId / 1000 // Supondo que cada prefixo tem 1000 IDs
	if companyIndex >= 0 && companyIndex < len(externalCompanies) {
		companyURL := externalCompanies[companyIndex]

		// Se a companhia é local (FlightId começa com 0), tenta reservar localmente
		if companyIndex == 0 {
			for i := range flights {
				if flights[i].FlightId == flightId {
					// Verifica se o assento já está reservado
					if seatID < 0 || seatID >= len(flights[i].Seats) || flights[i].Seats[seatID].IsReserved {
						return errors.New("reserva não realizada") // Assento inválido ou já reservado
					}

					// Realiza a reserva
					flights[i].Seats[seatID].IsReserved = true
					flights[i].Seats[seatID].CustomerID = cpf

					// Salva a atualização
					return SaveFlights(flights)
				}
			}
			return errors.New("voo não encontrado")
		}

		// Se a companhia é externa, faz a requisição de reserva
		reqBody, _ := json.Marshal(models.ReserveSeatRequest{
			FlightId: flightId,
			SeatID:   seatID,
			CPF:      cpf,
		})
		

		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Post(fmt.Sprintf("%s/reserve-seat", companyURL), "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Printf("Erro ao tentar reservar na companhia %s: %v\n", companyURL, err)
			return errors.New("reserva não realizada")
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil // Reserva bem-sucedida em companhia externa
		}
		return errors.New("reserva não realizada")
	}

	return errors.New("companhia não encontrada para o voo")
}