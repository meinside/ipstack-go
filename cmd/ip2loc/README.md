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

#### one ip address

```shell
$ $GOPATH/bin/ip2loc 1.1.1.1
```

```bash
1.1.1.1 / one.one.one.one (Australia)
```

#### multiple ip addresses (not supported on free plan)

```shell
$ $GOPATH/bin/ip2loc 8.8.8.8 8.8.4.4
```

```bash
8.8.8.8 / google-public-dns-a.google.com (United States)
8.8.4.4 / google-public-dns-a.google.com (United States)
```
