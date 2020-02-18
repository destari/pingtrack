# pingtrack
Track latency and packet loss with a minimal web UI.

Built with:
* Golang
* React
* Gatsbyjs
* Bulma (rbx)


## Obligatory screenshot:
![Example screenshot of graphs](https://github.com/destari/pingtrack/blob/master/Screen%20Shot%202020-02-17%20at%208.26.48%20PM.png)

## Execution
```./pingtrack -H localhost -p 8080 -i 5 hosts google.com 192.168.2.1 192.168.2.10```

You can skip all the parameters except the hosts (for now). So you can also do this:

```./pingtrack hosts google.com```

In the web app, you can add/remove hosts dynamically.

## Building
In order to build, you will need Gatsby installed, then you can run `make` in the `cmd/` directory to build the full static binary.
