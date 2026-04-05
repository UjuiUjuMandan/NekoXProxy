# NekoXProxy

## 原理

本程序作为本地HTTP代理运行。当Telegram客户端请求某个数据中心IP时，程序根据`dcmap.json`将该IP映射到中继域名的子域名，并将请求转发至中继服务器。

```
Telegram 客户端 → NekoXProxy → http(s)://<子域名>.<中继域名>/api → Cloudflare → Telegram DC
```

## 使用方法

```
./NekoXProxy -p http://mtproto.example.com -c dcmap.json -l 127.0.0.1:26641
```

| 参数 | 默认值            | 说明                                              |
| ---- | ----------------- | ------------------------------------------------- |
| `-p` | 必填              | 中继服务器 base URL，使用 `http://` 或 `https://` |
| `-c` | `dcmap.json`      | IP → 子域名映射配置文件                           |
| `-l` | `127.0.0.1:26641` | 本地监听地址                                      |

## dcmap.json 格式

键为 Telegram 服务器 IP（或 IP 前缀），值为对应的子域名。

```json
{
  "149.154.175.5": "dc1",
  "95.161.76.100": "dc2",
  "91.108.56.": "dc5"
}
```

若 `-p http://mtproto.example.com`，则 `149.154.175.5` 的请求会被转发到 `http://dc1.mtproto.example.com/api`。

IP 前缀（如 `91.108.56.`）用于匹配同一段内的所有地址。

## 自行搭建中继

1. 准备一个域名，如 `mtproto.example.com`
2. 为每个 Telegram 数据中心设置子域名 DNS 记录，指向对应的 Telegram DC IP：
   ```
   1.mtproto.example.com → 2001:b28:f23d:f001::a
   2.mtproto.example.com → 2001:67c:4e8:f002::a
   3.mtproto.example.com → 2001:b28:f23d:f003::a
   4.mtproto.example.com → 2001:67c:4e8:f004::a
   5.mtproto.example.com → 2001:b28:f23d:f005::a
   ...
   ```
3. 编辑 `dcmap.json`，填入各 IP 到子域名的映射

## 评论

猫耳逆变器 @tehcneko 在2022年2月[发布](https://t.me/NekoUpdates/223)了名为tcp2ws的闭源jar文件、后来又开发了GUI程序[WSProxy](https://github.com/Nekogram/WSProxy)。

世界 @nekohasekai 和 @arm64-v8a 重新实现并开源了它。

然而，所有对WS的提及是错误和误导性的、因为没有使用WebSocket，Cloudflare仅代理了HTTP Post请求。

只有Telegram Web版本会使用WebSocket传输，要将其转换为一般MTProxy，请见：

- https://github.com/Flowseal/tg-ws-proxy
- https://github.com/valnesfjord/tg-ws-proxy-rs
