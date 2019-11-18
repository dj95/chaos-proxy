<h1 align="center">chaos-proxy</h1>

<p align="center">
  A proxy for testing out the behaviour of network applications with bandwidth limitations, latency and packet loss.
  <br><br>
  <a href="https://cloud.drone.io/dj95/chaos-proxy">
    <img alt="BuildStatus" src="https://cloud.drone.io/api/badges/dj95/chaos-proxy/status.svg" />
  </a>
  <a href="https://github.com/dj95/chaos-proxy/actions?query=workflow%3AGo">
    <img alt="GoActions" src="https://github.com/dj95/chaos-proxy/workflows/Go/badge.svg" />
  </a>
  <a href="https://godoc.org/github.com/dj95/chaos-proxy/pkg/proxy">
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
  <a href="https://sonarcloud.io/dashboard?id=dj95:chaos-proxy">
    <img src="https://sonarcloud.io/api/project_badges/measure?project=dj95:chaos-proxy&metric=sqale_rating" />
  </a>
  <a href="https://sonarcloud.io/dashboard?id=dj95:chaos-proxy">
    <img src="https://sonarcloud.io/api/project_badges/measure?project=dj95:chaos-proxy&metric=alert_status" />
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

If you are missing features or find some annoying bugs please feel free to submit an issue or a bugfix within a pull request :)


## 🚧 TODO

- implement udp handler


## 📝 License

© 2019 Daniel Jankowski


This project is licensed under the MIT license.


Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:


The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.


THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

