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