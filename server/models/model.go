package models

type Seat struct {
	IsReserved bool   `json:"is_reserved"`
	CustomerID string `json:"customer_id"`
}

type Flight struct {
	FlightId    int    `json:"flightId"`
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Seats       []Seat `json:"seats"`
}

type Client struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}

type Config struct {
    CompanyName       string   `json:"companyName"`
    Port              int      `json:"port"`
    ExternalCompanies []string `json:"externalCompanies"`
}

type ReserveSeatRequest struct {
	FlightId int    `json:"flightId"`
	SeatID   int    `json:"seatID"`
	CPF      string `json:"cpf"`
}

type CancelReservationRequest struct {
	FlightId int `json:"flightId"` // ID do voo
	SeatID   int `json:"seatID"`   // ID do assento
	CPF      string `json:"cpf"`    // CPF do cliente
}
