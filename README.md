# Slime Proxy
![](https://github.com/hoveychen/slime/raw/main/docs/mascot.png =300x)

Slime Proxy is a robust HTTP reverse proxy based on a unique Hub-Agent architecture. It is designed to streamline the process of managing multiple isomorphic service providers, who provide the same API, by grouping them into a single gateway. This reverse proxy setup simplifies the process of managing multiple service providers and ensures a seamless user experience.
## Architecture
![](https://github.com/hoveychen/slime/raw/main/docs/architecture.png =600x)
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
## Getting Started
// Here you can add instructions on how to install and use your project.
## Contributing
// Here you can add instructions on how to contribute to your project.
## License
// Here you can add information about the license of your project.