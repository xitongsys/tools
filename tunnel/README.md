# tunnel-proxy

This is a tool for tcp port forward and reverse proxy tunnel, which supports encryption.

## run

```go
Usage of ./tunnel_proxy:
  -D string
        run as daemon with address
  -L string
        log file, default is stdout/stderr
  -P string
        password (default "12345")
```

## command

```
ls
	show info

ls tunname
	show tun info

open tun tun_name remote_addr remote_password
	open a tun connection named tun_name to remote_addr with remote_password

open listen tun_name direction listen_addr forward_addr
	direction: l(listen on local) or r(listen on remote)
	listen on listen_addr and forward the connection package to forward_addr by the tun
	
close tun tun_name
	close tun

close listen tun_name direction listen_id
	close listen

exit
	exit
```

