# ✈️ API-Vendepass - Sistema de Venda

---

## 📖 Índice
1. [Introdução](#🌟-introdução)
2. [Metodologia Utilizada](#🛠️-metodologia-utilizada)
3. [Arquitetura da Solução](#🏗️-arquitetura-da-solução)
   - [Estrutura do Repositório](#estrutura-do-repositório)
4. [Protocolo de Comunicação](#📡-protocolo-de-comunicação)
   - [APIs Principais](#apis-principais)
5. [Roteamento e Token Ring](#🔄-roteamento-e-token-ring)
   - [Token Ring](#token-ring)
   - [Configuração do Sistema](#configuração-do-sistema)
6. [Discussão e Resultados](#💬-discussão-e-resultados)
   - [Concorrência Distribuída](#concorrência-distribuída)
   - [Possível Solução e Melhoria Futura](#possível-solução-e-melhoria-futura)
   - [Confiabilidade da Solução](#confiabilidade-da-solução)
   - [Emprego do Docker](#emprego-do-docker)
7. [Arquivo de Testes](#🔎-arquivo-de-testes)
8. [Como Rodar o Projeto](#🚀-como-rodar-o-projeto)
   - [Pré-requisitos](#pré-requisitos)
   - [Passo a Passo](#passo-a-passo)
9. [Conclusão](#🏆-conclusão)
10. [Autores](#✍️-autores)

---

## 🌟 Introdução
Este projeto apresenta o desenvolvimento de um **sistema distribuído de reservas de passagens aéreas**, no qual clientes podem adquirir passagens de diferentes companhias aéreas por meio de uma interface de comunicação entre servidores. O sistema permite que um cliente conecte-se ao servidor de uma companhia aérea (por exemplo, **Companhia A**) e, a partir deste servidor, possa comprar passagens de outras companhias (por exemplo, **Companhia C**) de forma distribuída.

**Objetivo:** Demonstrar a viabilidade de um ambiente distribuído com controle de concorrência, garantindo que múltiplas transações de reserva possam ocorrer simultaneamente sem conflitos.

---

## 🛠️ Metodologia Utilizada
A implementação foi realizada em **Golang**, com a biblioteca **Gin** para o gerenciamento de rotas HTTP, permitindo o roteamento e tratamento das requisições de forma eficiente e organizada. A arquitetura é modular, dividida em pastas específicas para dados, manipuladores de requisições (handlers), modelos, repositórios e roteadores, o que facilita a manutenção e extensibilidade do código.

---

## 🏗️ Arquitetura da Solução
A arquitetura do sistema é composta por vários componentes interligados, cada um desempenhando um papel essencial:

- **Servidor Principal (Companhia A)**: Recebe as requisições dos clientes e intermedia a comunicação com outros servidores de companhias aéreas.
- **Servidores de Companhias Aéreas (Companhia C, etc.)**: Armazenam as informações de voos e executam operações de reserva em suas respectivas bases de dados.
- **Armazenamento de Dados (JSON)**: Utiliza arquivos JSON para armazenar informações sobre clientes (`clients.json`), rotas de voo (`routes.json`), e configurações de rede e servidores (`config.json`).

---
#### Estrutura do Repositório:

```bash
server
│
├───data
│       clients.json         // Armazena dados dos clientes (nome, identificação, contato).
│       config.json          // Contém configurações de rede e das companhias aéreas.
│       routes.json          // Contém informações sobre rotas de voos (origem, destino, horários).
├───handler
│       cancelReservationHandler.go // Manipulador para cancelar reservas de passagens.
│       listFlightsHandler.go       // Manipulador para listar voos disponíveis.
│       registerClientHandler.go    // Manipulador para cadastrar novos clientes.
│       reserveSeatHandler.go       // Manipulador para reservar assentos em voos.
├───models
│       model.go                    // Define as estruturas de dados e modelos utilizados no sistema.
├───repository
│       fileRepository.go           // Implementação de repositório para leitura e escrita em arquivos JSON.
│       tokenRing.go                // Implementação do Token Ring para coordenação das requisições distribuídas.
├───router
│       router.go                   // Configuração das rotas e ligação entre endpoints e manipuladores.
│       routes.go                   // Define as rotas disponíveis na aplicação, ligando endpoints aos manipuladores.
├───test
│       main.go                     // Testes do sistema.
│
│   go.mod                   // Gerencia as dependências do projeto em Go.
│   go.sum                   // Registro das versões das dependências do projeto.
│   main.go                  // Ponto de entrada da aplicação, configura o servidor e as rotas.
│   docker-compose.yml              // Configuração do Docker Compose para o ambiente distribuído.
│   dockerfile                      // Dockerfile para criação dos contêineres.
│   .gitignore                      // Arquivos e pastas ignorados pelo Git.
│   README.md                       // Documentação principal do sistema.
```

---

## 📡 Protocolo de Comunicação
A comunicação entre os servidores utiliza APIs RESTful, possibilitando operações distribuídas para efetuar reservas. Foi introduzido um cabeçalho customizado nas requisições (`X-Request-Type`) para diferenciar requisições de máquinas (servidores) e de clientes. Requisições de máquinas usam o valor `machine`, enquanto as requisições de clientes usam `client`.

### APIs Principais
- **Reserva de Assento (POST /reserve)**  
  - **Descrição**: Reserva um assento em um voo específico.
  - **Parâmetros**: Identificador do voo, informações do cliente.
  - **Retorno**: Status da reserva (confirmada ou não).

- **Cancelamento de Reserva (POST /cancelReservation)**  
  - **Descrição**: Cancela uma reserva existente.
  - **Parâmetros**: Identificador da reserva, informações do cliente.
  - **Retorno**: Status do cancelamento.

- **Cadastro de Cliente (POST /registerClient)**  
  - **Descrição**: Cadastra um novo cliente.
  - **Parâmetros**: Dados do cliente (nome, identificação, contato).
  - **Retorno**: Confirmação do cadastro.

- **Listagem de Voos (GET /listFlights)**  
  - **Descrição**: Retorna a lista de voos disponíveis entre uma origem e um destino.
  - **Parâmetros**: Origem, destino.
  - **Retorno**: Lista de voos com detalhes como horário e disponibilidade de assentos.

Essas APIs garantem a comunicação eficiente entre os servidores e permitem que clientes conectados ao servidor de uma companhia aérea possam reservar passagens em outras companhias.

---

## 🔄 Roteamento e Token Ring

### Token Ring
Para controlar a comunicação e coordenar as transações entre diferentes servidores, foi implementado um mecanismo de **Token Ring**. Esse mecanismo funciona passando um "token" entre os servidores participantes, garantindo que apenas um servidor tenha permissão para realizar uma operação em determinado momento. Isso evita conflitos e melhora o controle de concorrência, especialmente em transações distribuídas, como a reserva de assentos.

### Configuração do Sistema
O arquivo `config.json` armazena as configurações de rede, incluindo os endereços dos servidores e informações de inicialização para cada companhia. Esse arquivo permite que o sistema ajuste automaticamente as configurações para incluir novos servidores sem precisar alterar o código diretamente.

---

## 💬 Discussão e Resultados

### Concorrência Distribuída

Com a implementação do **Token Ring**, o sistema distribui eficientemente as transações de reserva de forma controlada, minimizando conflitos e mantendo a consistência dos dados. Cada servidor, ao receber o token, tem permissão exclusiva para executar transações críticas, garantindo a integridade dos dados em um ambiente distribuído

No entanto, o **Token Ring** como implementado apresenta uma limitação: **se o servidor que possui o token cair**, o sistema perde temporariamente o controle sobre as operações críticas. Isso pode resultar em inconsistências e na paralisação do sistema até que o token seja recuperado ou regenerado.

### Possível Solução e Melhoria Futura

Para resolver essa limitação, deve-se implementar um mecanismo de **detecção de falhas** e **geração de token atualizado**. Esse mecanismo incluirá:

1. **Timeout do Token**:  
   Caso um servidor não receba o token dentro de um período definido, ele assume que o servidor anterior falhou.

2. **Reconfiguração Dinâmica do Anel**:  
   O sistema deve ser capaz de reorganizar o anel de servidores automaticamente, excluindo o servidor que falhou, para que o token continue circulando.

3. **Geração de Novo Token**:  
   Se o servidor que detinha o token cair, o sistema gera um **token com valor maior** para substituir o anterior. Isso garante que todos os servidores reconheçam o novo token como válido.

4. **Mecanismo de Recuperação**:  
   Quando o servidor que caiu volta a operar, ele verifica o valor do token atual. Se perceber que seu token está desatualizado, ele o descarta e se reintegra ao anel, recebendo o token atualizado na próxima rodada.

Essas melhorias são fundamentais para garantir a **resiliência e disponibilidade do sistema** em cenários distribuídos e serão consideradas como parte do desenvolvimento futuro.

### Confiabilidade da Solução
O sistema foi projetado para continuar funcional mesmo quando algum servidor de companhia aérea fica temporariamente indisponível. Porém, por conta de como o token ring foi implementado, a reserva de assentos fica indisponível em qualquer companhia, mas as operações de listagem de vôos e assentos, cancelamento de assentos (de empresas online) e criação de usuários ainda estão disponíveis.

### Emprego do Docker
O uso de contêineres **Docker** possibilitou o desenvolvimento de um ambiente distribuído isolado, no qual cada servidor de companhia aérea opera em seu próprio contêiner. Essa configuração facilita a simulação de um sistema distribuído real e permite a escalabilidade, já que novos servidores podem ser adicionados rapidamente. Além disso, o **Docker Compose** simplifica a orquestração dos contêineres e o gerenciamento do ambiente de desenvolvimento e testes.

---

## 🔎 Arquivo de Testes

O arquivo `main.go` dentro da pasta `test` contém scripts para **testar o sistema** de reserva de passagens aéreas. Ele realiza testes automatizados de ponta a ponta para garantir que as operações de:
- Cadastro de clientes
- Listagem de voos
- Reserva e cancelamento de passagens

Os testes ajudam a identificar possíveis problemas de concorrência, falhas de comunicação e verificam a resposta do sistema a cada operação, garantindo a robustez e a confiabilidade da aplicação.

---

## 🚀 Como Rodar o Projeto

### Pré-requisitos

- **Docker** e **Docker Compose** instalados na máquina.

### Passo a Passo

1. **Clone o Repositório**
   ```bash
   git clone https://github.com/seu-repositorio/api-vendepass.git
   cd api-vendepass
   ```

2. **Configurar o Docker Compose** O arquivo docker-compose.yml já está configurado para levantar três contêineres, cada um representando um servidor de companhia aérea. Certifique-se de que as portas e endereços no config.json estão corretos para evitar conflitos de rede.
<br/>
3. **Iniciar o Docker Compose** Execute o comando abaixo para levantar todos os servidores ao mesmo tempo.

    ```bash
      docker-compose up --build
    ```
4. **Verificar a Inicialização** Após o comando acima, o Docker Compose exibirá os logs de inicialização de cada contêiner. Verifique se todos os servidores foram inicializados corretamente e se as portas configuradas estão abertas.
---

## 🏆 Conclusão

Este projeto demonstra a viabilidade de um sistema de reservas de passagens aéreas distribuído utilizando Docker e um protocolo de comunicação distribuído. A implementação do Token Ring e a arquitetura modular em Golang asseguram um controle de concorrência eficaz.

---

## ✍️ Autores
- [Gabriel Baptista](https://github.com/BaptistaGabriel)
- [Amanda Lima](https://github.com/AmandaLimaB)

---