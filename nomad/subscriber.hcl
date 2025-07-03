job "aeron-poc-subscriber" {
  datacenters = ["aws-main"]

  constraint {
    attribute = "${attr.unique.hostname}"
    operator  = "="
    value     = "ryan-85"
  }

  group "aeron-poc-subscriber" {
    count = 1

    spread {
      attribute = "${unique.hostname}"
      weight    = 100
    }

    ephemeral_disk {
      migrate = false
      size    = 50
      sticky  = false
    }

    network {
      mode = "host"
      port "gpp" {
        static = 10010
        to     = 10010
      }
    }

    restart {
      attempts = 15
      interval = "5m"
      delay    = "15s"
      mode     = "delay"
    }
    shutdown_delay = "15s"

    task "aeron-poc-subscriber" {
      logs {
        max_files     = 4
        max_file_size = 10
      }

      driver = "docker"

      config {
        image    = "10.1.9.155:5000/aeron-go-poc:f8efbd5"
        ports = ["gpp"]
        shm_size = 536870912
      }

      env = {
        "INJ_COMMAND" = "basicSubscriber",
        "INJ_CHANNEL" = "aeron:udp?endpoint=0.0.0.0:10010"
        "INJ_IDLE"    = "yield"
      }
      resources {
        cpu        = 100
        memory     = 1000
        memory_max = 1000
      }
      #      kill_timeout = "120s"
    }

  }
}