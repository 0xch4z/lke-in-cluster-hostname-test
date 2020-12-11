resource "linode_lke_cluster" "cluster" {
  label       = "example-cluster"
  k8s_version = "1.17"
  region      = "us-east"

  pool {
    type  = "g6-standard-2"
    count = 3
  }
}

resource "null_resource" "kubeconfig_fetch" {
  triggers = {
    kubeconfig = local.kubeconfig_string
  }

  provisioner "local-exec" {
    command = "echo '${local.kubeconfig_string}' > kubeconfig.yaml"
  }
}

locals {
  kubeconfig_string = base64decode(linode_lke_cluster.cluster.kubeconfig)
  kubeconfig        = yamldecode(local.kubeconfig_string)
}
