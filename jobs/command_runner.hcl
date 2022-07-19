job "{{ .RemoteCommandRunner.Name }}" {
  datacenters = [
    "dc1"
  ]

  group "command-runner" {
    task "runner" {
      driver = "raw_exec"

      config {
        command = "bash"
        args = [
          "-c",
          "for (( ; ; )); do sleep 3600; done;" # just run forever
        ]
      }
      
      resources {
        cpu    = 500
        memory = 512
      }
    }
  }
}