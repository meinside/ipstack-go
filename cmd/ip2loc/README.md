# ip2loc

Prints simple location information of given ip address.

## install

```bash
$ go get -u github.com/meinside/ipstack-go/cmd/ip2loc
```

## setup

Get your access key from [here](https://ipstack.com/),

then create a file named `ip2loc.json` in your `$HOME/.config/` directory:

```json
{
	"access_key": "PUT_YOUR_ACCESS_KEY_HERE",
	"is_premium": false
}
```

## run

### run without any param:

```bash
$ $GOPATH/bin/ip2loc
```

then it will print the location info of the external ip address of your machine.

### run with param(s):

```bash
# one ip address
$ $GOPATH/bin/ip2loc 1.1.1.1

# multiple ip addresses (not supported on free plan)
$ $GOPATH/bin/ip2loc 8.8.8.8 8.8.4.4
```

