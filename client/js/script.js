document.addEventListener("DOMContentLoaded", () => {
    const registerForm = document.getElementById("registerForm");
    const flightsList = document.getElementById("flightsList");
    const registerMessage = document.getElementById("registerMessage");

    // Função para registrar um cliente
    registerForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        const cpf = document.getElementById("cpf").value;
        const password = document.getElementById("password").value;

        try {
            const response = await fetch("http://localhost:8080/register-client", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ cpf, password }),
            });
            if (response.ok) {
                registerMessage.innerText = "Cadastro realizado com sucesso!";
                registerForm.reset();
            } else {
                registerMessage.innerText = "Falha no cadastro.";
            }
        } catch (error) {
            console.error("Erro no cadastro:", error);
            registerMessage.innerText = "Erro no cadastro. Tente novamente.";
        }
    });

    // Função para listar voos
    async function fetchFlights() {
        try {
            const response = await fetch("http://localhost:8080/list-flights");
            const flights = await response.json();

            flightsList.innerHTML = "";
            flights.forEach((flight) => {
                const flightDiv = document.createElement("div");
                flightDiv.classList.add("flight");

                let seatsStatus = "";
                flight.seats.forEach((seat, index) => {
                    seatsStatus += `Assento ${index + 1}: ${seat.is_reserved ? "Reservado" : "Disponível"}\n`;
                });

                flightDiv.innerHTML = `
                    <h3>Voo ${flight.flightId}</h3>
                    <p>Origem: ${flight.origin}</p>
                    <p>Destino: ${flight.destination}</p>
                    <pre>${seatsStatus}</pre>
                `;
                flightsList.appendChild(flightDiv);
            });
        } catch (error) {
            console.error("Erro ao buscar voos:", error);
            flightsList.innerHTML = "Não foi possível carregar os voos.";
        }
    }

    // Chamar a função de busca de voos ao carregar a página
    fetchFlights();
});
