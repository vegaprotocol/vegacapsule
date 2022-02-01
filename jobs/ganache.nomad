job "ganache-2" {
  datacenters = ["dc1"]
  group "ganache" {

    count = 1
    task "mydvbits-ganache" {
      driver = "docker"

      config {
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
        network {
          mbits = 10
          port "http" {
            static = 8545
          }
        }
      }
    }
  }
}
