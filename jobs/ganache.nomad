job "ganache-1" {
  datacenters = ["dc1"]
  group "ganache" {

    network {
      port "http" {
        static = 8545
      }
    }

    update {
      health_check = "task_states"
    }

    count = 1
    task "mydvbits-ganache" {

      driver = "docker"

      config {
        ports = ["http"]
        hostname = "mydvbits-ganache"
        image = "ghcr.io/vegaprotocol/devops-infra/ganache:latest"
        command = "ganache-cli"
        args  = [
          "--blockTime", "1",
          "--chainId", "1440",
          "--networkId", "1441",
          "-h", "0.0.0.0",
          "-p", "8545",
          "-m", "cherry manage trip absorb logic half number test shed logic purpose rifle",
          "--db", "/app/ganache-db",
        ]
      }

      resources {
        cpu    = 500 # 500 MHz
        memory = 1024 # 256MB
      }
    }
  }
}
