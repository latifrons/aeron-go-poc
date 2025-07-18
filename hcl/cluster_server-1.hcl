variable "image" {
  type    = string
  default = "arr:v1.1"
}
variable "nodeId" {
  type    = number
  default = 1
}
variable "clusterHostNames" {
  type    = string
  default = "10.2.3.5,10.2.3.5,10.2.3.5"
}

job "cluster-server-node1" {
  datacenters = ["aws-main"]

  group "cluster-server-node" {
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
      port "gpp-archive-control" {
        static = 9001 + var.nodeId * 100
        to     = 9001 + var.nodeId * 100
      }
      port "gpp-ingress" {
        static = 9002 + var.nodeId * 100
        to     = 9002 + var.nodeId * 100
      }
      port "gpp-consensus" {
        static = 9003 + var.nodeId * 100
        to     = 9003 + var.nodeId * 100
      }
      port "gpp-log" {
        static = 9004 + var.nodeId * 100
        to     = 9004 + var.nodeId * 100
      }
      port "gpp-transfer" {
        static = 9005 + var.nodeId * 100
        to     = 9005 + var.nodeId * 100
      }
    }

    restart {
      attempts = 20
      interval = "5m"
      delay    = "15s"
      mode     = "delay"
    }
    shutdown_delay = "15s"

    task "cluster-server-node" {
      logs {
        max_files     = 4
        max_file_size = 10
      }

      driver = "docker"

      config {
        network_mode = "host"
        shm_size     = "268435456"
        sysctl = {
          # "net.ipv4.ip_local_port_range" = "10000 65535"
        }
        image = "${var.image}"
        ports = ["gpp-archive-control", "gpp-ingress", "gpp-consensus", "gpp-log", "gpp-transfer"]
      }

      env = {
        "aeron.idle"                    = "busyspin", # client
        "aeron.driver.lowLatency"       = "1", # server
        "aeron.cluster.nodeId"          = "${var.nodeId}",
        "aeron.cluster.hostnames"       = "${var.clusterHostNames}",
        "aeron.cluster.clusterDir"      = "/data/cluster-${var.nodeId}",
        "aeron.cluster.aeronDir"        = "/dev/shm/aeron-md",
        "aeron.cluster.ingressStreamId" = "110",
        "command"                       = "clusterServerEcho"
      }
      resources {
        cpu        = 100
        memory     = 200
        memory_max = 1000
      }
      kill_timeout = "15s"
    }
  }
}




