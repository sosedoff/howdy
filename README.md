# howdy

Experimental YAML-based service monitoring thingy.

## Install

Via go:

```
go get github.com/sosedoff/howdy
```

Or download a [binary release](releases):

```
wget -O /usr/local/bin/howdy RELEASE_URL
chmod +x /usr/local/bin/howdy
```

## Config

All stuff goes into config:

```yml
name: myapp
enabled: true

notify:
  slack:
    webhook: "https://hooks.slack.com/services/..."
    # channel is optional
    channel: "#myapp"

checks:
  # check live urls
  web:
    # simple url
    - url: http://myapp.com/

    # specify expected http status code
    - url: http://myapp.com/login
      status: 403

    # specify code and format
    - url: http://api.myapp.com/v1/foobar
      status: 200 # default
      format: json

    # formats are: html, json, xml
    # or it could be anything that's included into Content-Type response header

  # check if host is reachable
  ping:
    - host: myapp.com

  # check if domains are resolving correctly
  dns:
    - domain: myapp.com
      ip: 1.2.3.4

  # check if port is reachable. default timeout is 5 seconds
  port:
    - host: myapp.com
      port: 22

    - host: 127.0.0.1
      port: 5000
      net: udp # default is tcp
```

## Usage

See all available options:

```
$ howdy -v

Usage of howdy:
  -c="": Path to all configs
  -n=true: Send notifications
  -t=false: Test mode
  -v=false: Show version
```

Run checks on all configs:

```
$ howdy -c /path/to/configs
```

Or run checks on a single config:

```
$ howdy /path/to/config.yml
```

## TODO

- Add support for notifications via email
- Add ability to detect service recovery. That will require having a DB.

## License

MIT