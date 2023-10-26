# aws-cost-notify-go

## 開発手順

- lambda で実行
  - 月初から月末までの金額の合計を取得
  - サービスごとに集計
  - or タグごとに集計
- slack に投稿
- webhookurl を引き渡す

## デプロイ

```
cd aws-cost-notify-go

# 1 build binary
make build

# 2 deploy by pulumi
cd iac
pulumi up
```
