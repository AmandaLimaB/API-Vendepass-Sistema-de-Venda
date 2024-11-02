# ✈️ API-Vendepass - Sistema de Venda

---

## 📖 Índice
1. [Introdução](#introdução)
2. [Metodologia Utilizada](#metodologia-utilizada)
3. [Arquitetura da Solução](#arquitetura-da-solução)
4. [Protocolo de Comunicação](#protocolo-de-comunicação)
5. [Roteamento](#roteamento)
6. [Discussão e Resultados](#discussão-e-resultados)
   - [Concorrência Distribuída](#concorrência-distribuída)
   - [Confiabilidade da Solução](#confiabilidade-da-solução)
   - [Emprego do Docker](#emprego-do-docker)
7. [Conclusão](#conclusão)
8. [Autores](#autores)

---

## 🌟 Introdução
Este projeto apresenta o desenvolvimento de um **sistema distribuído de reservas de passagens aéreas**, onde clientes podem adquirir passagens de diferentes companhias aéreas através de uma interface de comunicação entre servidores. O sistema permite que um cliente conecte-se ao servidor de uma companhia aérea (por exemplo, **Companhia A**) e, a partir desse servidor, possa comprar passagens de outras companhias (por exemplo, **Companhia C**) de forma distribuída.

**Objetivo:** Demonstrar a viabilidade de um ambiente distribuído com controle de concorrência, garantindo que múltiplas transações de reserva possam ocorrer simultaneamente sem conflitos.

---

## 🛠️ Metodologia Utilizada
A implementação foi realizada em **Golang**, com o uso da biblioteca **Gin** para o gerenciamento de rotas HTTP, que permite o roteamento e tratamento das requisições de maneira eficiente e organizada. A arquitetura segue uma abordagem modular, dividida em pastas específicas para dados, manipuladores de requisições (handlers), modelos, repositórios e roteadores, de modo a facilitar a manutenção e extensibilidade do código.

---

## 🏗️ Arquitetura da Solução
A arquitetura do sistema é composta por vários componentes interligados, cada um desempenhando um papel essencial:
- **Servidor Principal (Companhia A)**: Recebe as requisições dos clientes e intermedia a comunicação com outros servidores de companhias aéreas.
- **Servidores de Companhias Aéreas (Companhia C, etc.)**: Armazenam as informações de voos e executam operações de reserva em suas respectivas bases de dados.
- **Data Store (JSON)**: Utiliza arquivos JSON para armazenar informações sobre clientes e rotas de voos, que são acessados e manipulados pelos repositórios.

A arquitetura pode ser classificada como uma **arquitetura distribuída RESTful**, na qual cada servidor funciona como um serviço independente, mas que se comunica com os demais para concluir transações distribuídas.

#### Arquivos do repositório e suas funções:

```bash
server
│   go.mod                   // Gerencia as dependências do projeto em Go.
│   go.sum                   // Registro das versões das dependências do projeto.
│   main.go                  // Ponto de entrada da aplicação, configura o servidor e as rotas.

├───data
│       clients.json         // Armazena dados dos clientes (nome, identificação, contato).
│       routes.json          // Contém informações sobre rotas de voos (origem, destino, horários).

├───handler
│       cancelReservationHandler.go // Manipulador para cancelar reservas de passagens.
│       handler.go                  // Contém funções comuns para o gerenciamento das requisições.
│       listFlightsHandler.go       // Manipulador para listar voos disponíveis.
│       registerClientHandler.go    // Manipulador para cadastrar novos clientes.
│       request.go                  // Estruturas e funções para tratar as requisições HTTP.
│       reserveSeatHandler.go       // Manipulador para reservar assentos em voos.

├───models
│       model.go                   // Define as estruturas de dados e modelos utilizados no sistema.

├───repository
│       fileRepository.go          // Implementação de repositório para leitura e escrita em arquivos JSON.

└───router
        router.go                  // Configuração das rotas e ligação entre endpoints e manipuladores.
        routes.go                  // Define as rotas disponíveis na aplicação, ligando endpoints aos manipuladores.
```

---

## 📡 Protocolo de Comunicação
A comunicação entre os servidores é implementada por meio de **APIs RESTful**, que permitem a realização de operações distribuídas para efetuar reservas. Os métodos incluem:

- **Reserva de Assento (POST /reserve)**  
  - **Descrição**: Reserva um assento em um voo específico.  
  - **Parâmetros**: Identificador do voo, informações do cliente (nome, identificação).  
  - **Retorno**: Status da reserva (confirmada ou não).

- **Cancelamento de Reserva (POST /cancelReservation)**  
  - **Descrição**: Cancela uma reserva existente.  
  - **Parâmetros**: Identificador da reserva, informações do cliente.  
  - **Retorno**: Status do cancelamento (sucesso ou erro).

- **Cadastro de Cliente (POST /registerClient)**  
  - **Descrição**: Cadastra um novo cliente no sistema.  
  - **Parâmetros**: Dados do cliente (nome, identificação, contato).  
  - **Retorno**: Confirmação do cadastro (sucesso ou erro).

- **Listagem de Voos (GET /listFlights)**  
  - **Descrição**: Retorna a lista de voos disponíveis entre uma origem e um destino.  
  - **Parâmetros**: Origem, destino.  
  - **Retorno**: Lista de voos correspondentes, com detalhes como horário e disponibilidade de assentos.

Essas APIs garantem a comunicação eficiente entre os servidores e possibilitam que clientes do servidor da **Companhia A** possam reservar passagens em outros servidores.

---

## 🗺️ Roteamento
A implementação do cálculo das rotas entre origens e destinos será realizada utilizando uma abordagem baseada em **grafos**. Cada aeroporto será representado como um nó, e as rotas de voo como arestas, permitindo a utilização de algoritmos de busca, como o de **Dijkstra**, para identificar as melhores opções de voos. Essa funcionalidade será integrada à API por meio de uma nova rota, facilitando a consulta de voos disponíveis.

---

## 💬 Discussão e Resultados

### 🔄 Concorrência Distribuída
Para resolver o problema de concorrência em sistemas distribuídos e evitar conflitos nas reservas de passagens, será implementado um mecanismo de controle com **relógio vetorial**. Esse mecanismo permite que cada transação seja marcada com um timestamp, gerando uma ordem de eventos distribuída que evita conflitos de reserva. 

### 🔒 Confiabilidade da Solução
O sistema é projetado para continuar a funcionar mesmo com a desconexão de alguns servidores de companhias aéreas. Quando um servidor é desconectado, o sistema interrompe temporariamente as reservas na companhia aérea afetada, mas as demais operações continuam disponíveis. Isso aumenta a **resiliência** e **disponibilidade** do sistema.

### 🐳 Emprego do Docker
Foi planejado o uso de **contêineres Docker** para simplificar o deploy e o teste do sistema. A utilização do Docker permite que cada servidor de companhia aérea seja executado em contêineres separados, simulando um ambiente distribuído real. Com essa abordagem, é possível configurar múltiplas instâncias dos servidores, viabilizando testes de concorrência e robustez do sistema em um ambiente controlado.

---

## ✅ Conclusão
O desenvolvimento deste sistema distribuído de reservas de passagens aéreas abordou diversos aspectos da concorrência e comunicação distribuída, com o objetivo de garantir uma experiência de reserva confiável e consistente. Embora algumas funcionalidades estejam em fase de planejamento ou implementação, como o uso do relógio vetorial para controle de concorrência, o sistema já apresenta uma estrutura modular e funcional, com potencial para ser ampliado.

---

## ✍️ Autores
- [Gabriel Baptista](https://github.com/BaptistaGabriel)
- [Amanda Lima](https://github.com/AmandaLimaB)

---
