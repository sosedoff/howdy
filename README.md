# howdy

Experimental YAML-based service monitoring thingy.

## Config

All stuff goes into config:

```yml
name: myapp
enabled: true

checks:
  # check live urls
  web:
    # just simple ur
    - url: http://myapp.com/

    # specify expected http status code
    - url: http://myapp.com/login
      status: 403

    # specify code and format
    - url: http://api.myapp.com/v1/foobar
      status: 200
      format: json

  # check if host is reachable
  ping:
    - host: myapp.com

  # check if domains are resolving correctly
  dns:
    - domain: myapp.com
      ip: 1.2.3.4
```