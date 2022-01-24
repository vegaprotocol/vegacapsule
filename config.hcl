vega_path = "async"

service "http" "web_proxy" {
  listen_addr = "127.0.0.1:8080"
  
  process "version" {
    command = ["echo", "version jede"]
  }

  process "web_proxy" {
    command = ["echo", "starting web proxy ${pid}"]
  }
}