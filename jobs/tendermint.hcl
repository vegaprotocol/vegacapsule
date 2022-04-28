job "tendermint" {
  datacenters = ["dc1"]

  group "tendermint" {
    count = 1
    task "tm" {
      driver = "exec"

      config {
        command = "/Users/karelmoravec/go/bin/vega"
        args = [
          "tm",
          "--home", "testnet/tendermint/node0",
        ]
      }
      resources {
        cpu    = 500 # 500 MHz
        memory = 1024 # 256MB
      }
    }
  }
}
