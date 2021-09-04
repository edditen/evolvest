# evolvest

kv storage

## quick start

1. start server:

```shell script
make build
bin/evolvestd -c conf/config.yaml
```

2. start client(support redis-cli):

```shell script
redis-cli -p 8762
```

3. checking data via client:

```shell script
bin/evolvestcli -a 127.0.0.1:8763
```
