job "aeron-poc-publisher" {
  datacenters = ["aws-main"]

  constraint {
    attribute = "${attr.unique.hostname}"
    operator  = "="
    value     = "ryan-85"
  }

  group "aeron-poc-publisher" {
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
    }

    restart {
      attempts = 15
      interval = "5m"
      delay    = "15s"
      mode     = "delay"
    }
    shutdown_delay = "15s"

    task "aeron-poc-publisher" {
      logs {
        max_files     = 4
        max_file_size = 10
      }

      driver = "docker"

      config {
        image    = "10.1.9.155:5000/aeron-go-poc:f8efbd5"
        shm_size = 536870912
      }

      env = {
        "INJ_COMMAND" = "basicPublisher",
        "INJ_CHANNEL" = "aeron:udp?endpoint=10.1.9.85:10010",
        "INJ_IDLE"    = "busyspin"
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