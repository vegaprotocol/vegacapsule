job "{{ .RemoteCommandRunner.Name }}" {
  // Currently impossible to wildcard datacenters so we have to list all our DCs
  // Ref: https://github.com/hashicorp/nomad/issues/9024
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