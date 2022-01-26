job "ganache" {
  datacenters = ["dc1"]
  group "ganache" {

network {
  port "ganache_port" {
    static = 8545
    to     = 8545
  }
}
    count = 1
    task "mydvbits-ganache" {
      driver = "docker"
      config {
        hostname = "mydvbits-ganache"
        image = "trufflesuite/ganache-cli:v6.12.2"
        args  = [
          "--blockTime", "1",
          "--chainId", "1440",
          "--networkId", "1441",
          "-h", "0.0.0.0",
          "-p", "8545",
          "-m", "cherry manage trip absorb logic half number test shed logic purpose rifle"
        ]
        ports = ["ganache_port"]
      }
      resources {
        cpu    = 500 # 500 MHz
        memory = 1024 # 256MB
      }
    }
  }
}
