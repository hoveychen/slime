# Slime Proxy

![go build](https://github.com/hoveychen/slime/actions/workflows/go.yml/badge.svg)
[![GoDoc reference example](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/hoveychen/slime)
[![GoReportCard example](https://goreportcard.com/badge/github.com/hoveychen/slime)](https://goreportcard.com/report/github.com/hoveychen/slime)
[![Coverage Status](https://coveralls.io/repos/github/hoveychen/slime/badge.svg?branch=main)](https://coveralls.io/github/hoveychen/slime?branch=main)

<img src="https://github.com/hoveychen/slime/raw/main/docs/mascot.png" width="300px">

Slime Proxyは、ユニークなHub-Agentアーキテクチャに基づく堅牢なHTTPリバースプロキシです。同じAPIを提供する複数の同型サービスプロバイダを一つのゲートウェイにまとめることで、サービスプロバイダの管理プロセスを効率化することを目指して設計されています。このリバースプロキシの設定は、複数のサービスプロバイダの管理プロセスを簡素化し、シームレスなユーザーエクスペリエンスを保証します。

## アーキテクチャ

<img src="https://github.com/hoveychen/slime/raw/main/docs/architecture.png" width="600x">

アーキテクチャは主に三つのコンポーネントで構成されています：

- **Hub**: Hubサーバーはアーキテクチャの中心です。HTTPリクエストを受け付け、適切なエージェントにインテリジェントに転送します。これにより、効率的な負荷分散が可能となり、最適なパフォーマンスを保証します。

- **Agent**: エージェントサーバーは仲介者として機能し、トラフィックを上流サーバーにプロキシします。彼らはハブとの持続的な接続を維持し、着信リクエストを処理する準備をしています。

- **Upstream**: これらはアプリケーションリクエストのレスポンスを生成する実際のサービスプロバイダです。必要なAPIを提供する任意のサーバーになることができます。

## 主な特徴

Slime Proxyは、伝統的なロードバランサーとは異なるいくつかのユニークな特徴を提供します：

- **Dynamic Agents**: ハブは複数のダイナミックエージェントを組織します。これらのエージェントは、サーバーに接続するためのプルメカニズムを使用し、ハブからの公開アクセスの必要性を排除します。これにより、より安全で柔軟なアーキテクチャが可能となります。

- **Connection Pool**: ハブは利用可能な接続プールを維持し、アクティブな同時アプリケーションリクエストを制限します。この設計により、エージェントはアプリケーションが並列呼び出しを行っても、単一のスレッドで安全に実行することができ、リソースのオーバーヘッドを削減します。

## 使用法

Slime Proxyは、同じAPIを提供する複数のサービスプロバイダを持ち、それらを一つのゲートウェイを通じて管理したいシナリオに最適です。マイクロサービスアーキテクチャ、マルチクラウド環境、あるいは単純にサービスの高可用性とパフォーマンスを確保したい場合に最適です。

### インストール

プロジェクトをインストールするには、モジュールが有効化されたGolangバージョン1.18以上が必要です。次のコマンドを使用してプロジェクトをインストールできます：
```bash
go install -u github.com/hoveychen/slime@latest
```
または、提供されたリンクから事前にコンパイルされたバイナリをダウンロードすることもできます。

[https://github.com/hoveychen/slime/releases](https://github.com/hoveychen/slime/releases)

## はじめに

プロキシを初期化するには、最低でも一つのハブと一つのエージェントが必要です。

### ハブの設定

まず、ハブのための`<secret>`を生成します。この`<secret>`は任意の文字列で、ランダムなパスワードジェネレータから生成されることが望ましいです。それは安全に、そしてプライベートに保管されるべきです。もし漏洩した場合、プロキシは偽造されたエージェントからの攻撃に対して脆弱になります。
次に、次のコマンドを使用してハブサーバーを実行します：

```bash
slime hub run --secret <secret>
```
> 本番環境では、`host`と`port`の設定を明示的にバインドするフラグに加えて、`concurrent`フラグを適切な値（例えば、`1024`）に設定することを推奨します。これにより、分散型サービス拒否（DDoS）攻撃の可能性を軽減するのに役立ちます。

### エージェントの設定
まず、エージェントがハブにアクセスするための*Agent Token*を生成します。これは次のコマンドを使用して行うことができます：
```bash
slime hub register --secret <secret> --name <my agent name>
```
このコマンドは暗号化されたエージェントトークンを出力します。複数のエージェント間でエージェントトークンを再利用することは可能ですが、監査目的とトークンの再ロールのために、各エージェントにユニークなエージェントトークンを割り当てることをお勧めします。
次に、次のコマンドを使用してエージェントサーバーを実行します：
```bash
slime agent run --token <agent token> --hub <hub address> --upstream <upstream address> 
```
> 通常、一つのエージェントは一つの上流サービスを担当します。一つのコマンドで複数の上流サービスに対して複数のエージェントを設定するには、以下に示すように、カンマで区切られた複数の上流アドレスを指定します：
> ```bash
> slime agent run --token <agent token> --hub <hub address> --upstream <upstream1>,<upstream2>,<upstream3>
> ```
> このシナリオでは、上流プロバイダーに対して同数のエージェントが設定されます。

> [!NOTE]
> デフォルトの設定では、サービスプロバイダーが単一スレッドモード（例えば、GPUを使用したヘビーロードの生成的AIタスク）で動作することを想定しています。これが該当しない場合は、`numWorker`フラグを指定して並列度を増やすことができます。

## 貢献
貢献は歓迎されます。気軽に問題を開き、マージリクエストを提出してください。

## ライセンス
Apache 2.0