#ip
ipconfig.io clone written in golang using [gorilla/mux](https://github.com/gorilla/mux). Separate handlers for `"Accept: application/json"`, headless and browser requests. 

![master](https://github.com/iofq/ip/actions/workflows/ghcr.yml/badge.svg)
### usage

```
$ curl ip.iofq.net
127.0.0.1

$ wget -qO- ip.iofq.net
127.0.0.1

$ bat -print=b ip.iofq.net
127.0.0.1
```

Or, to run:
```
docker run -p 8080:8080 -d ghcr.io/iofq/ip
```

