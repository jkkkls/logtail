# logtail
日志收集工具

### 执行方式

``` shell
$logtail -d /data/logtail/

```

### 配置


```
file: /data/games/slog/{date}/client.log
# json, txt, txt需要指定分隔符和fields
format: txt
separator: ;
fields:
    -
        name: time
        type: string
        default: xx
    -
        name: channel
        type: int
        default: 2
    -
        name: amount
        type: float
        default: 1.1

out:
    console:
        enabled: true
    kafak:
        enabled: true
        hosts: ["kafka01:8101"]
        topic: game1:client
        username:
        password:
```