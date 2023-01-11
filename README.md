# locks-and-concurrency

Load testing RDBMS locks in high concurrency using different queries to check data integrity, latency and throughput. 


It simulates a transfer money feature between accounts

Scenarios:
- Inconsistent
- Pessimistic
- Optimistic
- Optimized

# Inconsistent
Queries that can't guarantee data integrity 

```sql
BEGIN

-- GET Accounts 1 and 2
SELECT * FROM account WHERE id=$1;
SELECT * FROM account WHERE id=$1;

-- Create Transfer
INSERT INTO transfer (amount, from_account_id, to_account_id) VALUES ($1, $2, $3) RETURNING *;

-- UPDATE accounts setting its final balance
UPDATE account SET balance=$1, version=version+1 WHERE id=$2 RETURNING *;
UPDATE account SET balance=$1, version=version+1 WHERE id=$2 RETURNING *;

COMMIT
```

✓ status was 200

     checks.........................: 100.00% ✓ 3865       ✗ 0
     data_received..................: 896 kB  43 kB/s
     data_sent......................: 703 kB  34 kB/s
     http_req_blocked...............: avg=34.83µs min=0s     med=3µs    max=1.16ms   p(90)=9µs     p(95)=330µs
     http_req_connecting............: avg=25.66µs min=0s     med=0s     max=939µs    p(90)=0s      p(95)=270.79µs
     http_req_duration..............: avg=9.22ms  min=3.6ms  med=5.41ms max=146.68ms p(90)=17.38ms p(95)=24.45ms
       { expected_response:true }...: avg=9.22ms  min=3.6ms  med=5.41ms max=146.68ms p(90)=17.38ms p(95)=24.45ms
     http_req_failed................: 0.00%   ✓ 0          ✗ 3865
     http_req_receiving.............: avg=32.21µs min=8µs    med=25µs   max=254µs    p(90)=60µs    p(95)=71µs
     http_req_sending...............: avg=19.88µs min=4µs    med=14µs   max=534µs    p(90)=37µs    p(95)=50µs
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=9.17ms  min=3.58ms med=5.36ms max=146.65ms p(90)=17.35ms p(95)=24.42ms
     http_reqs......................: 3865    187.151998/s
     iteration_duration.............: avg=1s      min=1s     med=1s     max=1.14s    p(90)=1.01s   p(95)=1.02s
     iterations.....................: 3865    187.151998/s
     vus............................: 33      min=33       max=299
     vus_max........................: 300     min=300      max=300

# Pessimistic
Queries using pessimistic lock to guarantee data integrity 

```sql
BEGIN

-- GET and LOCK Accounts 1 and 2
SELECT * FROM account WHERE id=$1 FOR UPDATE; 
SELECT * FROM account WHERE id=$1 FOR UPDATE;

-- Create Transfer
INSERT INTO transfer (amount, from_account_id, to_account_id) VALUES ($1, $2, $3) RETURNING *;

-- UPDATE accounts setting its final balance
UPDATE account SET balance=$1, version=version+1 WHERE id=$2 RETURNING *;
UPDATE account SET balance=$1, version=version+1 WHERE id=$2 RETURNING *;

COMMIT
```


✓ status was 200

     checks.........................: 100.00% ✓ 3802       ✗ 0
     data_received..................: 875 kB  42 kB/s
     data_sent......................: 688 kB  33 kB/s
     http_req_blocked...............: avg=42.95µs min=0s     med=4µs    max=8.16ms   p(90)=9µs     p(95)=373µs
     http_req_connecting............: avg=31.93µs min=0s     med=0s     max=8.05ms   p(90)=0s      p(95)=296.94µs
     http_req_duration..............: avg=27.23ms min=3.57ms med=6.13ms max=590.62ms p(90)=42.13ms p(95)=206.72ms
       { expected_response:true }...: avg=27.23ms min=3.57ms med=6.13ms max=590.62ms p(90)=42.13ms p(95)=206.72ms
     http_req_failed................: 0.00%   ✓ 0          ✗ 3802
     http_req_receiving.............: avg=41.43µs min=9µs    med=37µs   max=401µs    p(90)=68.9µs  p(95)=80µs
     http_req_sending...............: avg=24.68µs min=4µs    med=21µs   max=437µs    p(90)=41µs    p(95)=59µs
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=27.16ms min=3.52ms med=6.06ms max=590.56ms p(90)=42.1ms  p(95)=206.67ms
     http_reqs......................: 3802    182.987002/s
     iteration_duration.............: avg=1.02s   min=1s     med=1s     max=1.59s    p(90)=1.04s   p(95)=1.2s
     iterations.....................: 3802    182.987002/s
     vus............................: 27      min=27       max=299
     vus_max........................: 300     min=300      max=300

# Optimistic
Queries using optimistic lock to guarantee data integrity

```sql
BEGIN

-- GET Accounts 1 and 2
SELECT * FROM account WHERE id=$1;
SELECT * FROM account WHERE id=$1;

-- Create Transfer
INSERT INTO transfer (amount, from_account_id, to_account_id) VALUES ($1, $2, $3) RETURNING *;

-- UPDATE accounts using version column
UPDATE account SET balance=$1, version=version+1 WHERE id=$2 AND version=$3 RETURNING *;
UPDATE account SET balance=$1, version=version+1 WHERE id=$2 AND version=$3 RETURNING *;

COMMIT
```

✗ status was 200
↳  45% — ✓ 1761 / ✗ 2091

     checks.........................: 45.71% ✓ 1761       ✗ 2091
     data_received..................: 750 kB 36 kB/s
     data_sent......................: 693 kB 34 kB/s
     http_req_blocked...............: avg=34.03µs min=0s     med=2µs    max=2.06ms   p(90)=9µs     p(95)=312.44µs
     http_req_connecting............: avg=25.4µs  min=0s     med=0s     max=1.14ms   p(90)=0s      p(95)=262.44µs
     http_req_duration..............: avg=11.18ms min=3.28ms med=6.3ms  max=114.52ms p(90)=22.91ms p(95)=30.69ms
       { expected_response:true }...: avg=8.47ms  min=3.91ms med=5.67ms max=99.91ms  p(90)=14.85ms p(95)=19.78ms
     http_req_failed................: 54.28% ✓ 2091       ✗ 1761
     http_req_receiving.............: avg=31.76µs min=6µs    med=25µs   max=2.55ms   p(90)=59µs    p(95)=73µs
     http_req_sending...............: avg=17.01µs min=3µs    med=12µs   max=246µs    p(90)=33µs    p(95)=47.44µs
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=11.13ms min=3.27ms med=6.26ms max=114.43ms p(90)=22.79ms p(95)=30.65ms
     http_reqs......................: 3852   186.558662/s
     iteration_duration.............: avg=1.01s   min=1s     med=1s     max=1.11s    p(90)=1.02s   p(95)=1.03s
     iterations.....................: 3852   186.558662/s
     vus............................: 31     min=31       max=299
     vus_max........................: 300    min=300      max=300

# Optimized-transfer
Optimized queries to guarantee data integrity and to avoid deadlocks with high throughput and low latency

```sql
BEGIN

-- Create Transfer
INSERT INTO transfer (amount, from_account_id, to_account_id) VALUES ($1, $2, $3) RETURNING *;

-- DEBIT from account
UPDATE account SET balance=balance - $1, version=version+1 WHERE id=$2 AND balance >= $1 RETURNING *;

-- CREDIT into account
UPDATE account SET balance=balance + $1, version=version+1 WHERE id=$2 AND RETURNING *;

COMMIT
```

✓ status was 200

     checks.........................: 100.00% ✓ 3856       ✗ 0
     data_received..................: 885 kB  43 kB/s
     data_sent......................: 652 kB  32 kB/s
     http_req_blocked...............: avg=36.92µs min=0s     med=3µs    max=1.35ms   p(90)=9µs     p(95)=304.74µs
     http_req_connecting............: avg=26.96µs min=0s     med=0s     max=838µs    p(90)=0s      p(95)=250µs
     http_req_duration..............: avg=9.64ms  min=2.58ms med=3.91ms max=339.51ms p(90)=13.85ms p(95)=22.61ms
       { expected_response:true }...: avg=9.64ms  min=2.58ms med=3.91ms max=339.51ms p(90)=13.85ms p(95)=22.61ms
     http_req_failed................: 0.00%   ✓ 0          ✗ 3856
     http_req_receiving.............: avg=36.17µs min=7µs    med=31µs   max=358µs    p(90)=62µs    p(95)=76µs
     http_req_sending...............: avg=23.13µs min=3µs    med=18µs   max=617µs    p(90)=40µs    p(95)=51µs
     http_req_tls_handshaking.......: avg=0s      min=0s     med=0s     max=0s       p(90)=0s      p(95)=0s
     http_req_waiting...............: avg=9.58ms  min=2.56ms med=3.84ms max=339.49ms p(90)=13.76ms p(95)=22.56ms
     http_reqs......................: 3856    186.358986/s
     iteration_duration.............: avg=1.01s   min=1s     med=1s     max=1.33s    p(90)=1.01s   p(95)=1.02s
     iterations.....................: 3856    186.358986/s
     vus............................: 29      min=29       max=299
     vus_max........................: 300     min=300      max=300


