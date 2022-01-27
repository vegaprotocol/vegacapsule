locals {
  path = abspath("testnet/smartcontracts")
}

job "smartcontracts" {
  datacenters = ["dc1"]
  type = "batch"
  group "smartcontracts" {
      network {
        port "smart_contracts_port" { static = 80 }
        }
    count = 1
    task "vegacapsule-smartcontracts" {
      driver = "docker"
      config {
        network_mode = "host"
        work_dir = "/app"
        image = "ghcr.io/vegaprotocol/devops-infra/smartcontracts:docker"
        entrypoint = ["/app/run"]
        volumes = ["${local.path}:/mnt"]
        ports = [
          "smart_contracts_port",
        ]
      }
      env {
        GANACHE_HOSTNAME = "127.0.0.1"
      }
      resources {
        cpu    = 500 # 500 MHz
        memory = 1024 # 256MB
      }
    }
  }
}
