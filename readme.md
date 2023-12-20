
*You can test connectivity to pod (using IP or name) and external domain*

## 💨 Quick test

**Test connectivity and exit**
```shell
kubectl testnet [endpoint:port]
```

## 🥡 Test connectivity on-demand

**Launch a container and test connectivity on-demand:**
```shell
# launch on-demand testnet pod 
kubectl testnet server -l app=toto
# trigger test to google.com:80
kubectl testnet client google.fr:80 -l app=toto
```

**See testnet logs:**
```bash
kubectl testnet client  --log -l app=toto
```

## Installation

Put `kubectl-testnet` in your `$PATH`
