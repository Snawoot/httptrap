httptrap
========

[![httptrap](https://snapcraft.io//httptrap/badge.svg)](https://snapcraft.io/httptrap)

Web-server which produces infinite chunked-encoded responses. It's intended to be a part of active defense against malicious HTTP clients to cancel their load impact and/or bruteforce efforts. Depending on settings, httptrap may keep busy most of attacker's resources, diverting them to wait slowly fed infinite response, or cause often OOM and botnet payload crash if httptrap sends response as fast as possible and client uses RAM to buffer entire response.

## Why?

Why just not collect attackers IP addresses and ban them? Or maybe even automate this process and establish it on ongoing basis?

Filtering approach has some drawbacks:

* It prevents access of legitimate users who share IP address with attacker (e.g. common ISP NAT gateway, compromised devices in network, common proxy or TOR exit node and so on). Consequently, it gives attacker an additional power to cut someone's access from service if one will manage originate flood traffic from same IP. Slowdown approach is more gentle and capable to not impact correct requests from legitimate users.
* Filtering may have lower stopping power if attacker posesses a huge pool of source addresses. In fact, if pool is sufficiently large and it's full rotation period is longer than IP cooldown time, attackers may conduct their activities continously, without downtimes.

## Installation

#### Pre-built binaries

Pre-built binaries available on [releases](https://github.com/Snawoot/httptrap/releases/latest) page.

#### From source

Alternatively, you may install httptrap from source:

```
go get github.com/Snawoot/httptrap
```

#### From Snap Store

[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/httptrap)

```sh
sudo snap install httptrap
```

#### Docker

```sh
docker run -it --rm -p 8008:8008 yarmak/httptrap
```

## Use Case

Consider following example. We have some web application which suffers from HTTP request flood or authorization bruteforce attempts. In this example such application represented by Python script [demo/webapp.py](demo/webapp.py). It serves HTTP requests on port 8080 and it is exposed to the outer world via nginx reverse proxy with simple server configuration section like this one:

```nginx
    server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  _;

        location / {
            proxy_pass http://127.0.0.1:8080;
        }

        error_page 404 /404.html;
            location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
            location = /50x.html {
        }
    }

```

Web application greets users which passed authorization:

```
$ curl -d 'login=admin&password=12345678' -D- http://localhost/
HTTP/1.1 200 OK
Server: nginx/1.16.1
Date: Sat, 11 Apr 2020 21:08:31 GMT
Content-Type: text/plain
Transfer-Encoding: chunked
Connection: keep-alive

You are in!
```

And it rejects unauthorized requests:

```
$ curl -d 'login=admin&password=badpass' -D- http://localhost/
HTTP/1.1 403 Forbidden
Server: nginx/1.16.1
Date: Sat, 11 Apr 2020 21:10:17 GMT
Transfer-Encoding: chunked
Connection: keep-alive

```

Let's model bruteforcer or flooder with Apache Benchmark `ab` utility:

```
$ time ab -c 50 -t 30 -s 5 -r -p post_data_bad.txt http://localhost/
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 5000 requests
Completed 10000 requests
Completed 15000 requests
Completed 20000 requests
Completed 25000 requests
Completed 30000 requests
Completed 35000 requests
Finished 39821 requests


Server Software:        nginx/1.16.1
Server Hostname:        localhost
Server Port:            80

Document Path:          /
Document Length:        0 bytes

Concurrency Level:      50
Time taken for tests:   30.000 seconds
Complete requests:      39821
Failed requests:        0
Non-2xx responses:      39821
Total transferred:      4141384 bytes
Total body sent:        6179850
HTML transferred:       0 bytes
Requests per second:    1327.36 [#/sec] (mean)
Time per request:       37.669 [ms] (mean)
Time per request:       0.753 [ms] (mean, across all concurrent requests)
Transfer rate:          134.81 [Kbytes/sec] received
                        201.17 kb/s sent
                        335.98 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       4
Processing:     1   19 209.2      4   27422
Waiting:        0   19 209.2      4   27422
Total:          1   19 209.2      4   27422

Percentage of the requests served within a certain time (ms)
  50%      4
  66%      5
  75%      5
  80%      6
  90%      7
  95%      9
  98%     13
  99%   1014
 100%  27422 (longest request)

real	0m30,046s
user	0m0,521s
sys	0m2,058s
```

Test lasted 30 seconds and allowed attacker to probe almost 40000 accounts (for purposes of discussion we omit that fact `ab` probed same single login and password pair). Obviously, Python server process hit 100% CPU limit during test, illustrating load impact caused by request flood.

And here is how we can derail such attacks without restrictions to normal users.

Let's add few lines to our nginx config so it'll be looking like this:

```nginx
    server {
        listen       80 default_server;
        listen       [::]:80 default_server;
        server_name  _;

        location / {
            proxy_intercept_errors on;
            proxy_pass http://127.0.0.1:8080;
        }

        error_page 404 /404.html;
        location = /40x.html {
        }

        error_page 500 502 503 504 /50x.html;
        location = /50x.html {
        }

        error_page 454 = @trap;
        location @trap {
            internal;
            proxy_pass http://127.0.0.1:8008;
            proxy_buffering off;
        }
    }
```

Namely, we've enabled error interception from backend responses and added named internal location which redispatches request to httptrap server. `error_page 454 = @trap` directive binds this internal location to HTTP response code 454. Now it is sufficient for backend to respond with HTTP code 454 and client will be dumped to endless response backend.

It's up to backend server how to identify malicious clients which should be locked out. It's possible to use complex behavior analysis or per-IP statistics - everything what fits your case. In our example web application uses simple stateless logic: 1% of requests which failed authorization responded with HTTP 454 error code and consequently redispatched to slow stream server. This way bruteforcers eventually will hit httptrap server. Let's see what impact such policy has on attacker:

```
$ time ab -c 50 -t 30 -s 5 -r -p post_data_bad.txt http://localhost/
This is ApacheBench, Version 2.3 <$Revision: 1843412 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
^C

Server Software:        nginx/1.16.1
Server Hostname:        localhost
Server Port:            80

Document Path:          /
Document Length:        0 bytes

Concurrency Level:      50
Time taken for tests:   368.656 seconds
Complete requests:      4565
Failed requests:        0
Non-2xx responses:      4565
Total transferred:      499165 bytes
Total body sent:        715325
HTML transferred:       18305 bytes
Requests per second:    12.38 [#/sec] (mean)
Time per request:       4037.849 [ms] (mean)
Time per request:       80.757 [ms] (mean, across all concurrent requests)
Transfer rate:          1.32 [Kbytes/sec] received
                        1.89 kb/s sent
                        3.22 kb/s total

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.3      0       4
Processing:     1   32 352.2      3    7657
Waiting:        0   32 352.2      3    7657
Total:          1   33 352.3      3    7659

Percentage of the requests served within a certain time (ms)
  50%      3
  66%      4
  75%      5
  80%      5
  90%      7
  95%     10
  98%     15
  99%   1023
 100%   7659 (longest request)

real	6m8,670s
user	0m0,140s
sys	0m0,438s
```

`ab` was started with 50 threads, 30 seconds total time limit and 5 second timeout per request. However, it didn't even completed at all: It made less than 5000 requests until all workers were locked out. It took 11 seconds till `ab` hang completely and after 6 minutes it was killed manually. Note that five second network operation timeout doesn't actually helps here because httptrap server feeds data continously, avoiding client timing out and leaving. All the time during test application was available to legitimate users, no per-address restrictions were imposed.

## Synopsis

```
$ ./httptrap -h
Usage of ./httptrap:
  -bind string
    	listen address (default ":8008")
  -cert string
    	enable HTTPS and use certificate
  -ct string
    	Content-Type value for responses (default "text/html")
  -interval duration
    	interval between chunks (default 1s)
  -key string
    	key for TLS certificate
  -string string
    	hex-encoded representation of byte string repeated in responses (default "0a")
  -verbosity int
    	logging verbosity (10 - debug, 20 - info, 30 - warning, 40 - error, 50 - critical) (default 20)
```
