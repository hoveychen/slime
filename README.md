# Slime Proxy
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
When you have golang >= 1.18 with module enabled, simply install the project
```
go install -u github.com/hoveychen/slime@latest
```

Or you may download the pre-built binary here
// TBC

## Getting Started

To bring up the proxy, at least one hub and one agent are required.

### Setup hub

First, you need to generate a `<secret>` for the hub. This `<secret>` can be any text, best generated from random password generator, should be kept safe and privately. Once leaked, the proxy is vulnerable to the forged agents.

Second, run the hub server by command:
```
slime hub run --secret <secret>
```

> Except for the explict flags binding config `host` and `port`, it's suggested to set the `concurrent` flag to a reasonable value like `1024` in the production environment, to prevent DDOS attack whatsoever.

### Setup agent

First, you need to generate an *Agent Token* for the agent to access to the hub. Execute command:

```
slime hub register --secret <secret> --name <my agent name>
```

The command will print a encrypted agent token to the standard output. Although you can reuse the agent token in multiple agents, it's suggested to assign one agent token per agent, to help audit and token reroll.

Second, run the agent server by command:
```
slime agent run --token <agent token> --hub <hub address> --upstream <upstream address> 
```

> Typically, one agent is responsible for one upstream service. To setup multiple agents for multiple upstream services in one command, you may specify multiple upstream addresses separated by comma like
> ```
> slime agent run --token <agent token> --hub <hub address> --upstream <upstream1>,<upstream2>,<upstream3>
> ```
> In such case, same number of agents to the upstream providers are setup.

> [!NOTE]
> The default setup assumes that the service provider works in single thread (for example, heavy-load generative AI task using GPU). If this is not the case, you may specify flag `numWorker` to increase the degree of parallelism. 

## Contributing
Feel free to open issues and merge requests. @hoveychen is actively maintaining this project.

## License
// Here you can add information about the license of your project.
