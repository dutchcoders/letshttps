# letshttps
Automatically requests Lets Encrypt certificates and forwards to webserver.

# Usage

```
$ go get github.com/dutchcoders/letshttps

$ letshttps -backend 127.0.0.1:80 -https :443 > /var/log/letshttps.log 2>&1 &
```

## Contributions

Contributions are welcome.

## Creators

**Remco Verhoef**
- <https://twitter.com/remco_verhoef>
- <https://twitter.com/dutchcoders>

## Copyright and license

Code and documentation copyright 2016 Remco Verhoef.
Code released under [the MIT license](LICENSE).
