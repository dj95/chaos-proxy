<h1 align="center">chaos-proxy</h1>

<p align="center">
  A proxy for testing out the behaviour of network applications with bandwidth limitations, latency and packet loss.
  <br><br>
  <a href="https://cloud.drone.io/dj95/chaos-proxy">
    <img alt="BuildStatus" src="https://cloud.drone.io/api/badges/dj95/chaos-proxy/status.svg" />
  </a>
  <a href="https://godoc.org/github.com/dj95/chaos-proxy">
    <img alt="GoDoc" src="https://godoc.org/github.com/dj95/chaos-proxy?status.svg" />
  </a>
  <a href="https://goreportcard.com/report/github.com/dj95/chaos-proxy">
    <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/dj95/chaos-proxy" />
  </a>
  <a href="https://github.com/dj95/chaos-proxy/releases">
    <img alt="latest version" src="https://img.shields.io/github/tag/dj95/chaos-proxy.svg" />
  </a>
  <a href="https://codecov.io/gh/dj95/chaos-proxy">
    <img src="https://codecov.io/gh/dj95/chaos-proxy/branch/master/graph/badge.svg" />
  </a>
</p>



## 📦 Requirements

- 🐳 Docker
- 🐙 docker-compose
- Golang(>=1.11)
- Make


## 🐳 Docker

Running the proxy server in docker requires to build the docker image.
Just navigate to [./deployments/docker](./deployments/docker) and run `docker-compose build`.

In order to configure the container you can use a config file, which will look like:

- Customize the [./configs/config.yml](./configs/config.yml)
- Run `docker run -p "8080:8080" -v "${PWD}/config.yml:/config.yml" chaos-proxy` in the same directory with the *config.yml*


## 🔧 Usage

- Customize the [./configs/config.yml](./configs/config.yml)
- Run `make run` in order to start the server
- When running the server the running proxy connections can be configured on the fly via an api. The api runs on a configured port under `/api`. For the documentation take a look at [./api/swagger.yml](./api/swagger.yml) or open the `/doc` route when running in development mode(log_level set to debug).


## 🤝 Contributing

If you are missing features or find some annoying bugs please feel free to submit the as issue and a bugfix with in a pull request :)


## 🚧 TODO

- implement udp handler


## 📝 License

(c) 2019 Daniel Jankowski


This project is licensed under the MIT license.
