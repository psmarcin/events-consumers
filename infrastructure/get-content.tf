resource "google_pubsub_topic" "get_content" {
  name = "get_content"
}

resource "google_storage_bucket" "get_content" {
  name = "get_content"
}

resource "google_storage_bucket_object" "get_content" {
  name   = "get_content_${uuid()}.zip"
  bucket = google_storage_bucket.get_content.name
  source = "./../dist/content.zip"

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}

resource "google_cloudfunctions_function" "get_content" {
  name        = "get_content"
  description = "Message by Terraform"
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.get_content.name
  source_archive_object = google_storage_bucket_object.get_content.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.get_content.name
  }
  timeout               = 10
  entry_point           = "Get"
  service_account_email = "root-481@events-consumer.iam.gserviceaccount.com"

  lifecycle {
    ignore_changes = [
      source_archive_bucket,
    ]
  }
}

resource "google_cloud_scheduler_job" "get_content" {
  name        = "get_content"
  description = "Message by Terraform"
  schedule    = "*/5 * * * *"

  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.get_content.id
    data       = base64encode("test")
  }
}
