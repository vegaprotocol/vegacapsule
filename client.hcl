plugin "docker" {
  config {
    volumes {
      enabled = true
    }
    auth {
      helper = "osxkeychain"
    }
  }
}

server {
  enable_event_broker = true
}