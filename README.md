HTTP-Verify
===========

A custom URI scheme that adds ECDSA based signatures to vanilla HTTP 1.1 requests.
The scheme can run over tcp or tls and is transparent to the underlying HTTP server.
Within this repository you will find simple examples, a proxy server and a library with useful constructs to build your very own httpv implementation!

Checkout the [wiki](https://github.com/NSkelsey/httpv/wiki) for a formal explanation!

##New Schemes

###httpv
Runs over http. Example: 

`httpv://localhost/api/status`

###httpsv
Runs over https. Example: 

`httpsv://ahimsa.io:1060/api/status`
