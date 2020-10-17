resource "google_container_cluster" "cluster" {
  name     = "test-cluster"
  location = "asia-northeast1-a"

  remove_default_node_pool = true
  initial_node_count       = 1
}

resource "google_container_node_pool" "normal" {
  name       = "normal"
  cluster    = google_container_cluster.cluster.name
  location   = "asia-northeast1-a"
  node_count = 1

  node_config {
    preemptible  = true
    machine_type = "e2-standard-2"

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}

resource "google_container_node_pool" "preemptible" {
  name       = "preemptible"
  location   = "asia-northeast1-a"
  cluster    = google_container_cluster.cluster.name
  node_count = 1

  node_config {
    preemptible = true

    machine_type = "e2-standard-2"

    taint = [
      {
        key    = "dedicated"
        value  = "preemptible"
        effect = "NO_SCHEDULE"
      }
    ]

    metadata = {
      disable-legacy-endpoints = "true"
    }

    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring",
    ]
  }
}
