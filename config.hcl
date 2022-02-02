vega_path = "/vega"
wallet_path = "./wallet"

network "testnet" {

  service "ganache" {
    port = "sd"
    config = "{{}}{{}}"
  }

  service "vega" "validator" {
    binary = "./path/to/bin"

    count = 3

    tendermint {
      initCmdArgs = [
        "--"
      ]
    
    }

    vega {
      configPlugins = [
        marketUserWillsPlugin,
        marketUserKarelsPlugin
        marketUserKarelsPlugin
      ]
    }

    port = "sd"
    config = "{{}}{{}}"
  }

   service "vega" "full" {
    count = 2

    tendermint {

    }

    vega {
      
    }

    port = "sd"
    config = "{{}}{{}}"
  }

}
wallet_path = "./wallet"

network "testnet" {

  service "eeth" {
    port = "sd"
    config = "{{}}{{}}"
  }

  service "vega" "validator" {
    binary = "./path/to/bin"

    count = 3

    tendermint {
      initCmdArgs = [
        "--"
      ]
    
    }

    vega {
      
      configPlugins = [
        marketUserWillsPlugin,
        marketUserKarelsPlugin
        marketUserKarelsPlugin
      ]

    }

    port = "sd"
    config = "{{}}{{}}"
  }

   service "vega" "full" {
    count = 2

    tendermint {

    }

    vega {
      
    }

    port = "sd"
    config = "{{}}{{}}"
  }

}


service "http" "web_proxy" {
  listen_addr = "127.0.0.1:8080"
  
  process "version" {
    command = ["echo", "version jede"]
  }

  process "web_proxy" {
    command = ["echo", "starting web proxy ${pid}"]
  }
}