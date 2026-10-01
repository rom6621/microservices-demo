**Online Boutique(オンラインブティック)** はクラウドファーストなマイクロサービスアプリケーションです。
WebベースのECアプリで、ユーザーは商品を探し、カートに入れ、購入することができます。

Googleはこのアプリケーションを、開発者が以下のGoogle Cloudのサービスを使ってモダンなエンタープライズアプリケーション開発を行う際のデモとして使用しています。

- Google Kubernets Engine (GKE)
- Cloud Service Mesh (CSM)
- gRPC
- Cloud Operations
- Spanner
- Memorystore
- AlloyDB
- Gemini

もしこのデモを使用する場合は、**★Star** をお願いします!

## アーキテクチャ

**Oneline Boutique** はそれぞれ異なる言語で書かれた11個のマイクロサービスがgRPCを通して通信しています。
[![Architecture of microservices](/docs/img/architecture-diagram.png)](/docs/img/architecture-diagram.png)

**Protobuf定義** は `./protos` ディレクトリを参照してください。

| サービス                                              | 言語            | 概要                                                                                                      |
| ----------------------------------------------------- | --------------- | --------------------------------------------------------------------------------------------------------- |
| [frontend](./src/frontend/)                           | Go              | WebサイトのHTTPサーバー。アカウント作成・ログインは不要で、自動的にセッションが作成される。               |
| [cartservice](./src/cartservice/)                     | C#              | Redisを使用し、ユーザーのカートへの保存と参照を行う。                                                     |
| [productcatalogservice](./src/productcatalogservice/) | Go              | JSONファイルから商品リストを提供し、検索・参照を行う。                                                    |
| [currencyservice](./src/currencyservice/)             | Node.js         | 通貨換算を行う。ヨーロッパ中央銀行から実際の値をフェッチする。QPS(Queries Per Second)が最も高いサービス。 |
| [paymentservice](./src/paymentservice/)               | Node.js         | クレジットカード情報(モック)を使用し、支払いを行う。支払い後、トランザクションIDを返す。                  |
| [shippingservice](./src/shippingservice/)             | Go              | カート情報から配送料の計算と、指定された住所への配送を行う。(モック)                                      |
| [emailservice](./src/emailservice/)                   | Python          | 注文確定のメールを送信する。(モック)                                                                      |
| [checkoutservice](./src/checkoutservice/)             | Go              | カート情報から、注文・決済・配送・メール通知の準備を行う。                                                |
| [recommendationservice](./src/recommendationservice/) | Python          | カート情報から商品のリコメンドを行う。                                                                    |
| [adservice](./src/adservice/)                         | Java            | 提供されたコンテキストワードからテキスト広告を提供する。                                                  |
| [loadgenerator](./src/loadgenerator/)                 | Python / Locust | ユーザーの購入までの流れを擬似的かつ継続的にフロントエンドに送り続ける。                                  |

## スクリーンショット

| ホーム画面                                                                                                                    | 決済画面                                                                                                              |
| ----------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------- |
| [![ストアホーム画面のスクリーンショット](/docs/img/online-boutique-frontend-1.png)](/docs/img/online-boutique-frontend-1.png) | [![決済画面のスクリーンショット](/docs/img/online-boutique-frontend-2.png)](/docs/img/online-boutique-frontend-2.png) |

## クイックスタート

1. 以下を用意してください

- Google Cloud Project
- `gcloud`, `git`, `kubectl`コマンド

2. 最新のメジャーバージョンをclone

```
git clone --depth 1 --branch v0
https://github.com/GoogleCloudPlatform/microservices-demo.git
cd microservices-demo/
```
