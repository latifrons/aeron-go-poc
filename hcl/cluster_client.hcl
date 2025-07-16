variable "image" {
  type    = string
  default = "654654541151.dkr.ecr.ap-east-1.amazonaws.com/exchange-api:v1.14.1"
}

job "cluster-server-0" {
  datacenters = ["aws-main"]

  group "cluster-server-0" {
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
      port "gpp-1" {
        static = 9001
        to     = 9001
      }
      port "gpp-2" {
        static = 9002
        to     = 9002
      }
      port "gpp-3" {
        static = 9003
        to     = 9003
      }
      port "gpp-4" {
        static = 9004
        to     = 9004
      }
      port "gpp-5" {
        static = 9005
        to     = 9005
      }
    }

    restart {
      attempts = 20
      interval = "5m"
      delay    = "15s"
      mode     = "delay"
    }
    shutdown_delay = "15s"

    task "cluster-server-0" {
      logs {
        max_files     = 4
        max_file_size = 10
      }

      driver = "docker"

      config {
        sysctl = {
          # "net.ipv4.ip_local_port_range" = "10000 65535"
        }
        image = "${var.image}"
        ports = ["gpp-1", "gpp-2", "gpp-3", "gpp-4", "gpp-5"]
      }

      env = {
        "aeron.cluster.nodeId"          = "0",
        "aeron.cluster.hostnames"       = "127.0.0.1",
        "aeron.cluster.clusterDir"      = "/data/cluster-0",
        "aeron.cluster.aeronDir"        = "/dev/shm/aeron-md",
        "aeron.cluster.ingressStreamId" = "102",
        "command"                       = "clusterServerEcho"
      }
      resources {
        cpu        = 101
        memory     = 200
        memory_max = 1000
      }
      kill_timeout = "15s"
    }
  }
}




