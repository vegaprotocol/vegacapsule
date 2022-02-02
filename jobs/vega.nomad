job "vega" {
  datacenters = ["dc1"]

  group "vega" {
    count = 1
    task "node" {
      driver = "exec"

      config {
        command = "/Users/karelmoravec/go/bin/vega"
        args = [
          "node",
          "--home", "testnet/vega/node0",
          "--nodewallet-passphrase-file", "testnet/vega/node0/node-vega-wallet-pass.txt",
        ]
      }
      resources {
        cpu    = 500 # 500 MHz
        memory = 1024 # 256MB
      }
    }
  }
}
