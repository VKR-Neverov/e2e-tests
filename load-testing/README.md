# Reports

## auth-service

### 08.03.2024 

Requests | Workers | Time | RPS

1000 | 10 | 5.63s | 178

1000 | 100 | 653ms | 1530

10000 | 100 | 21.19s | 472 - Server is down, huge degradation, website was down for ~6 minutes

shut down graylog and jaeger

10000 | 10 | 56.8s | 176 = ok no degradation

15000 | 15 | 2m6s | 119 = ok

20000 | 50 | 24s | 833 = ok

50000 | 100 | 46s | 1086 = ok
