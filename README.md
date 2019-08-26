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

- Customize the [./configs/config.yml](./configs/config.yml)
- Run `docker run -p "8080:8080" -v "${PWD}/config.yml:/config.yml" dj95/deception-proxy` in the same directory with the *config.yml*


## 🔧 Usage

- Customize the [./configs/config.yml](./configs/config.yml)
- Run `make run` in order to start the server


## 🤝 Contributing

If you are missing features or find some annoying bugs please feel free to submit the as issue and a bugfix with in a pull request :)


## 📝 License

(c) 2019 Daniel Jankowski


This project is licensed under the MIT license.
