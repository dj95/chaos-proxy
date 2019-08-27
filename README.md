<h1 align="center">deception-proxy 💀</h1>

<p align="center">
  A proxy for testing out the behaviour of network applications with bandwidth limitations, latency and packet loss.
  <br><br>
  <a href="https://godoc.org/github.com/dj95/deception-proxy">
    <img alt="GoDoc" src="https://godoc.org/github.com/dj95/deception-proxy?status.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/dj95/deception-proxy">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/dj95/deception-proxy" />
  </a>
  <a href="https://github.com/dj95/deception-proxy/releases">
    <img alt="latest version" src="https://img.shields.io/github/tag/dj95/deception-proxy.svg" />
  </a>
  <a href="https://hub.docker.com/r/dj95/deception-proxy">
    <img alt="Pulls from DockerHub" src="https://img.shields.io/docker/pulls/dj95/deception-proxy.svg?style=flat-square" />
  </a>
</p>


## 📦 Requirements

- 🐳 Docker
- 🐙 docker-compose
- Golang(>=1.11)
- Make


## 🐳 Docker

Running the proxy server in docker can be achieved with some different approaches.
One is to use pre-built docker images provided on docker hub.
If you don't feel comfortable with those images, feel free to build the docker images by yourself with the provided Dockerfile.
Just navigate to [./deployments/docker](./deployments/docker) and run `docker-compose build`.

In order to configure the container you can use a config file, which will look like:

- Customize the [./configs/config.yml](./configs/config.yml)
- Run `docker run -p "8080:8080" -v "${PWD}/config.yml:/config.yml" dj95/deception-proxy` in the same directory with the *config.yml*

Alternatively the container can be configured with environment variables, that are similar to the ones in the config file.
The commands therefore will be:

- `docker run -p "8080:8080" -e "DP_CORE_LOG_LEVEL=info" dj95/deception-proxy`

Every config variable that can be found in the config file will be available as environment variable as well.
Please bear in mind that set environment variables will override the values set in the config file.
The following environment variables are available:

- **DP_CORE_LOG_LEVEL** log_level of the server
    This key can be set to:
      *"debug"*
      *"info"*
    If another value is set, deception-proxy will fall
    back to the info level.
- **DP_CORE_LOG_FORMAT** specifies the format to log
    This key can be set to:
      *"text"*
      *"json"*
    If another value is set, the server will fall back
    to the text format.
- **DP_CORE_ADDRESS** address for the proxies socket
- **DP_CORE_PORT** port for the proxies socket
- **DP_CONN_TARGET** target url that should be proxied
- **DP_CONN_BANDWIDTH** sets the bandwidth in bit/sec.
    The default value is set to 8 Mbit/sec. If set to 0 or
    a negative value, the bandwidth will be unlimited.
- **DP_CONN_OVERHEAD** specifies the header overhead for the packages.
- **DP_CONN_LATENCY_MIN** configures the minimal latency for the requests
- **DP_CONN_LATENCY_MAX** configures the maximal latency for the requests
- **DP_CONN_LOSS_RATE** specifies the loss rate


## 🔧 Usage

- Customize the [./configs/config.yml](./configs/config.yml)
- Run `make run` in order to start the server


## 🤝 Contributing

If you are missing features or find some annoying bugs please feel free to submit the as issue and a bugfix with in a pull request :)


## 📝 License

(c) 2019 Daniel Jankowski


This project is licensed under the MIT license.
