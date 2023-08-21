# 史莱姆代理 Slime Proxy

![go build](https://github.com/hoveychen/slime/actions/workflows/go.yml/badge.svg)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hoveychen/slime)
[![GoReportCard example](https://goreportcard.com/badge/github.com/hoveychen/slime)](https://goreportcard.com/report/github.com/hoveychen/slime)
[![Coverage Status](https://coveralls.io/repos/github/hoveychen/slime/badge.svg?branch=main)](https://coveralls.io/github/hoveychen/slime?branch=main)

<img src="https://github.com/hoveychen/slime/raw/main/docs/mascot.png" width="300px">

史莱姆代理是一个基于独特的Hub-Agent架构的强大的HTTP反向代理。它旨在通过将多个同构服务提供商组合到一个单一的网关中，简化管理多个服务提供商的过程。这种反向代理设置简化了管理多个服务提供商的过程，并确保了无缝的用户体验。

## 架构

<img src="https://github.com/hoveychen/slime/raw/main/docs/architecture.png" width="600x">

该架构由三个主要组件组成：

- **Hub**：Hub服务器是架构的核心。它接受HTTP请求并智能地将它们转发给适当的代理。这样可以实现高效的负载分配，并确保最佳性能。

- **Agent**：代理服务器充当中间人，将流量代理到上游服务器。它们与Hub保持持久连接，准备处理传入的请求。

- **Upstream**：这些是为应用程序请求生成响应的实际服务提供商。它们可以是任何提供所需API的服务器。

## 主要特点

史莱姆代理提供了几个独特的功能，使其与传统的负载均衡器有所区别：

- **动态代理**：Hub组织多个动态代理。这些代理使用拉取机制与服务器连接，消除了从Hub进行公共访问的需要。这样可以实现更安全和灵活的架构。

- **连接池**：Hub维护一个可用的连接池，限制活动并发应用程序请求。这种设计允许代理在单个线程上安全运行，即使应用程序进行并行调用也能减少资源开销。

## 使用方法

史莱姆代理非常适合在您有多个提供相同API的服务提供商，并且希望通过一个单一的网关来管理它们的场景。它非常适合微服务架构、混合云环境，或者仅仅是为了确保服务的高可用性和性能。

### 安装

要安装该项目，请确保您的Golang版本为1.20或更高，并启用了模块。然后可以使用以下命令安装该项目：
```bash
go install -u github.com/hoveychen/slime@latest
```
或者，您可以从提供的链接下载预编译的二进制文件。

[https://github.com/hoveychen/slime/releases](https://github.com/hoveychen/slime/releases)

## 入门指南

要初始化代理，至少需要一个Hub和一个Agent。

### Hub配置

首先，为Hub生成一个`<secret>`。这个`<secret>`可以是任何字符串，最好是从随机密码生成器生成的。它应该被安全地和私下地存储。如果泄漏，代理将容易受到来自伪造代理的攻击。
接下来，使用以下命令执行Hub服务器：

```bash
slime hub run --secret <secret>
```

> [!NOTE]
> 1. 除了调整`host`和`port`参数，建议在生产环境中将`concurrent`参数设置为合理的值（例如`1024`），这有助于减轻潜在的分布式拒绝服务（DDoS）攻击。
> 2. 如果Hub托管在互联网上，请确保Hub、应用程序和代理之间的网络是绝对安全的。以下是一些常见的做法：
>    * 将Hub放在一个*HTTPS*代理后面，如Nginx、HAProxy。
>    * 设置（Web应用程序防火墙）WAF以保护Hub的安全。
>    * 设置`appPassword`标志以要求应用程序进行身份验证。

### Agent配置

首先，为代理生成一个*Agent Token*，以便访问Hub。可以使用以下命令完成：

```bash
slime hub register --secret <secret> --name <my agent name>
```

该命令将输出一个加密的代理令牌。虽然可以在多个代理之间重用代理令牌，但建议为每个代理分配一个唯一的代理令牌，以进行审计目的和令牌重新生成。
接下来，使用以下命令执行代理服务器：

```bash
slime agent run --token <agent token> --hub <hub address> --upstream <upstream address> 
```

> 通常，一个代理负责一个上游服务。要在单个命令中为多个上游服务配置多个代理，请指定以逗号分隔的多个上游地址，如下所示：
> ```bash
> slime agent run --token <agent token> --hub <hub address> --upstream <upstream1>,<upstream2>,<upstream3>
> ```
> 在这种情况下，为上游提供商设置相同数量的代理。

> [!NOTE]
> 默认配置假设服务提供商在单线程模式下运行（例如，使用GPU进行重负载生成性AI任务）。如果不是这种情况，可以通过指定`numWorker`标志来增加并行度。

### 应用程序请求
下游应用程序可以使用任何HTTP请求调用Hub。
* 如果Hub已经设置为需要`appPassword`，应用程序的HTTP请求应包含一个`Slime-App-Password`头。
* 然后将请求转发给远程服务提供商（如果有）。如果没有服务提供商，则返回状态码`503 Service Unavailable`。包含一个HTTP头`Slime-Block: 1`将阻塞请求，直到服务提供商可用。

## 贡献

欢迎贡献。请随时提出问题并提交合并请求。

## 许可证

Apache 2.0