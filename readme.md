## 💨 Quick test

**Test connectivity and exit**
```shell
kubectl testnet [endpoint:port]
```

## 🥡 Test connectivity on-demand

**See testnet logs:**
```bash
kubectl testnet client -l app=toto --log
```

**Launch a container and test connectivity on-demand:**
```shell
# launch on-demand testnet pod 
kubectl testnet server -l app=toto
# trigger test to google.com:80
kubectl testnet client google.fr:80 -l app=toto
```

## Installation

Put `kubectl-testnet` in your `$PATH`
