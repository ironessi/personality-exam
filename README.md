# GoFrame Template For SingleRepo

Quick Start: 
- https://goframe.org/quick

## Local configuration

Real configuration files are ignored by Git. Before running locally, copy the
example files and replace the placeholders:

```bash
cp config/config.example.yaml config/config.yaml
cp hack/config.example.yaml hack/config.yaml
cp manifest/config/config.example.yaml manifest/config/config.yaml
cp manifest/config/payment.example.yaml manifest/config/payment.yaml
```

Do not commit real database passwords, Redis passwords, SMS keys, payment
private keys, certificates, or log files.
