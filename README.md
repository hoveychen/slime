# Slime Proxy

![go build](https://github.com/hoveychen/slime/actions/workflows/go.yml/badge.svg)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hoveychen/slime)
[![GoReportCard example](https://goreportcard.com/badge/github.com/hoveychen/slime)](https://goreportcard.com/report/github.com/hoveychen/slime)
[![Coverage Status](https://coveralls.io/repos/github/hoveychen/slime/badge.svg?branch=main)](https://coveralls.io/github/hoveychen/slime?branch=main)

<img src="https://github.com/hoveychen/slime/raw/main/docs/mascot.png" width="300px">

[中文](https://github.com/hoveychen/slime/blob/main/docs/README_cn.md)
[日本語](https://github.com/hoveychen/slime/blob/main/docs/README_jp.md)

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

## Installation
### Docker

Docker hub repoistory: [https://hub.docker.com/r/hoveychen/slime](https://hub.docker.com/r/hoveychen/slime)

```bash
docker pull hoveychen/slime:latest
```


### Pre-built binaries

Download binaries: [https://github.com/hoveychen/slime/releases](https://github.com/hoveychen/slime/releases)


### Build from source

Alternatively, you can build the project. ensure that you have Golang version 1.20 or higher with module enabled. You can then install the project using the following command:

```bash
go install -u github.com/hoveychen/slime@latest
```

## Getting Started
To initialize the proxy, a minimum of one hub and one agent is required.

### Hub Configuration
Firstly, generate a `<secret>` for the hub. This `<secret>` can be any string, preferably generated from a random password generator. It should be stored securely and privately. If leaked, the proxy becomes susceptible to attacks from forged agents.

Optionally, assign an `<appPassword>` for the applications to invoke with. It's used to provide the most basic authentication to the applications. It's recommended when the hub is exposed to the unsafe environment like the Internet.

Next, execute the hub server using the following command:
```bash
slime hub run --secret <secret> --appPassword <appPassword> --port <port>
```

or with docker
```bash
docker run -d --restart always --name slime-hub -e SECRET=<secret> -e APP_PASSWORD=<appPassword> -p <port>:8080 hoveychen/slime:latest hub run
```
> [!NOTE]
> 1. It is recommended to set the `concurrent` flag to a reasonable value (e.g., `1024`) in a production environment, in addition to the explicit flags binding the `host` and `port` configurations. This helps to mitigate potential Distributed Denial of Service (DDoS) attacks.
> 2. If the hub is hosting on the Internet, make sure the network between the hub, applications and agents are in absolute safe. Here are some common practices:
>    * Host the hub behind a *HTTPS* proxy, like Nginx, HAProxy.
>    * Setup (Web Application Firewall) WAF to keep the hub safe.
>    * Set `appPassword` flag to require the application to authenticate.

### Agent Configuration
Firstly, generate an *Agent Token* for the agent to access the hub. This can be done using the following command:
```bash
slime hub register --secret <secret> --name <my agent name>
```
or with docker
```bash
docker run --rm -e SECRET=<secret> hoveychen/slime hub register --name <my agent name>
```
This command will output an encrypted agent token. While it is possible to reuse the agent token across multiple agents, it is advisable to assign a unique agent token to each agent for auditing purposes and token reroll.
Next, execute the agent server using the following command:
```bash
slime agent run --token <agent token> --hub <hub address> --upstream <upstream address> 
```
or with docker
```bash
# when the upstream is not in the same node with the agent
# or the upsteram service is in the same docker bridge network
docker run -d --restart always --name slime-agent \ 
           -e TOKEN=<agent token> -e HUB=<hub address> \
           -e UPSTREAM=<upstream address> agent run

# when the upstream provider and the agent are in the same node in Linux
docker run -d --restart always --name slime-agent \
           -e TOKEN=<agent token> -e HUB=<hub address> \
           -e UPSTREAM=127.0.0.1:<upstream port> \
           --network host agent run

# when the upstream provider and the agent are in the same node in Windows/Mac
docker run -d --restart always --name slime-agent \
           -e TOKEN=<agent token> -e HUB=<hub address> \
           -e UPSTREAM=host.docker.internal:<upstream port> agent run
```
> Typically, one agent is responsible for one upstream service. To configure multiple agents for multiple upstream services in a single command, specify multiple upstream addresses separated by commas as shown below:
> ```bash
> slime agent run --token <agent token> --hub <hub address> --upstream <upstream1>,<upstream2>,<upstream3>
> ```
> In this scenario, an equal number of agents are set up for the upstream providers.

> [!NOTE]
> The default configuration assumes that the service provider operates in a single-threaded mode (e.g., heavy-load generative AI tasks using GPU). If this is not the case, you can increase the degree of parallelism by specifying the `numWorker` flag.

### Application request
The downstream applications are free to invoke the hub with any HTTP request. 

* If the hub has been setup to require an `appPassword`, the application HTTP request should include a header `Slime-App-Password`.
* The requests are then forwarded to the remote service providers if any available. If there are no service providers, status `503 Service Unavailable` will be returned. Including a HTTP header `Slime-Block: 1` will block the request until service providers become available.

## Contributing
Contributions are welcome. Feel free to open issues and submit merge requests.

## License
Apache 2.0
