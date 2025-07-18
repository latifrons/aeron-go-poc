variable "image" {
  type    = string
  default = "arr:v1.1"
}
variable "clusterIngresses" {
  type    = string
  default = "0=10.2.3.5:9002,1=10.2.3.5:9102,2=10.2.3.5:9202"
}
variable "egress" {
  type    = string
  default = "10.2.3.5:10000"
}

job "cluster-server-node0" {
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
        shm_size     = "536870912"
        sysctl = {
          # "net.ipv4.ip_local_port_range" = "10000 65535"
        }
        image = "${var.image}"
        ports = ["gpp-archive-control", "gpp-ingress", "gpp-consensus", "gpp-log", "gpp-transfer"]
      }

      env = {
        "command"                 = "clusterLatencyCheckClient"
        "aeron.driver.dir"        = "/dev/shm/aeron-md"
        "aeron.driver.idle"       = "busyspin"
        "aeron.driver.lowLatency" = "1",

        "aeron.egressChannel"    = "aeron:udp?endpoint=${var.egress}"
        "aeron.egressStreamId"   = "111"
        "aeron.ingressEndpoints" = "${var.clusterIngresses}"
        "aeron.ingressStreamId"  = "110"
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




