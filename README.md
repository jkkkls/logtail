# logtail
日志收集工具

### 执行方式

``` shell
$logtail -d /data/logtail/

```

### 配置


```
name: game1
file: /data/games/slog/{date}/client.log
# 指定分隔符会按照fields拆分数据
separator: ;
fields:
    -
        name: time
        type: string
    -
        name: channel
        type: int
    -
        name: amount
        type: float

out:
    kafak:
        enabled: true
        hosts: ["kafka01:8101"]
        topic: game1:client
        username:
        password:
```