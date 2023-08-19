# Slime Proxy

![go build](https://github.com/hoveychen/slime/actions/workflows/go.yml/badge.svg)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hoveychen/slime)
[![GoReportCard example](https://goreportcard.com/badge/github.com/hoveychen/slime)](https://goreportcard.com/report/github.com/hoveychen/slime)

<img src="https://github.com/hoveychen/slime/raw/main/docs/mascot.png" width="300px">

Slime Proxy is a robust HTTP reverse proxy based on a unique Hub-Agent architecture. It is designed to streamline the process of managing multiple isomorphic service providers, who provide the same API, by grouping them into a single gateway. This reverse proxy setup simplifies the process of managing multiple service providers and ensures a seamless user experience.

## Architecture
<img src="https://github.com/hoveychen/slime/raw/main/docs/architecture.png" width="600x">
The architecture consists of three main components:

- **Hub**: The hub server is the heart of the architecture. It accepts HTTP requests and intelligently forwards them to the appropriate agents. This allows for efficient load distribution and ensures optimal performance.

- **Agent**: The agent servers act as intermediaries, proxying the traffic to the upstream server. They maintain a persistent connection with the hub, ready to handle incoming requests.

- **Upstream**: These are the actual service providers that generate responses for the application requests. They can be any servers that provide the necessary API.


## Key Features
Slime Proxy offers several unique features that set it apart from traditional load balancers:

- **Dynamic Agents**: The hub organizes multiple dynamic agents. These agents use a pull mechanism to connect to the server, eliminating the need for public access from the hub. This allows for a more secure and flexible architecture.

- **Connection Pool**: The hub maintains an available connection pool, which restricts the active concurrent application requests. This design allows the agent to safely run on a single thread, reducing resource overhead, even when the application makes parallel calls.

## Usage
Slime Proxy is ideal for scenarios where you have multiple service providers offering the same API and you want to manage them through a single gateway. It's perfect for microservice architectures, multi-cloud environments, or simply when you want to ensure high availability and performance for your services.

### Installation
To install the project, ensure that you have Golang version 1.18 or higher with module enabled. You can then install the project using the following command:
```bash
go install -u github.com/hoveychen/slime@latest
```
Alternatively, you can download the pre-compiled binary from the provided link.

[https://github.com/hoveychen/slime/releases]()


## Getting Started
To initialize the proxy, a minimum of one hub and one agent is required.

### Hub Configuration
Firstly, generate a `<secret>` for the hub. This `<secret>` can be any string, preferably generated from a random password generator. It should be stored securely and privately. If leaked, the proxy becomes susceptible to attacks from forged agents.
Next, execute the hub server using the following command:
```bash
slime hub run --secret <secret>
```
> It is recommended to set the `concurrent` flag to a reasonable value (e.g., `1024`) in a production environment, in addition to the explicit flags binding the `host` and `port` configurations. This helps to mitigate potential Distributed Denial of Service (DDoS) attacks.

### Agent Configuration
Firstly, generate an *Agent Token* for the agent to access the hub. This can be done using the following command:
```bash
slime hub register --secret <secret> --name <my agent name>
```
This command will output an encrypted agent token. While it is possible to reuse the agent token across multiple agents, it is advisable to assign a unique agent token to each agent for auditing purposes and token reroll.
Next, execute the agent server using the following command:
```bash
slime agent run --token <agent token> --hub <hub address> --upstream <upstream address> 
```
> Typically, one agent is responsible for one upstream service. To configure multiple agents for multiple upstream services in a single command, specify multiple upstream addresses separated by commas as shown below:
> ```bash
> slime agent run --token <agent token> --hub <hub address> --upstream <upstream1>,<upstream2>,<upstream3>
> ```
> In this scenario, an equal number of agents are set up for the upstream providers.

> [!NOTE]
> The default configuration assumes that the service provider operates in a single-threaded mode (e.g., heavy-load generative AI tasks using GPU). If this is not the case, you can increase the degree of parallelism by specifying the `numWorker` flag.

## Contributing
Contributions are welcome. Feel free to open issues and submit merge requests.

## License
Apache 2.0
