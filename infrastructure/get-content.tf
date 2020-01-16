variable "get_content_name" {
  type = string
  default = "get_content"
}

resource "google_pubsub_topic" "get_content" {
  name = var.get_content_name
}

resource "google_storage_bucket" "get_content" {
  name = var.get_content_name
}

resource "google_storage_bucket_object" "get_content" {
  name   = "${var.get_content_name}_${uuid()}.zip"
  bucket = google_storage_bucket.get_content.name
  source = "./../dist/get-content.zip"

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}

resource "google_cloudfunctions_function" "get_content" {
  name        = var.get_content_name
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

  environment_variables = {
    PROCESS_CONTENT_TOPIC_ID = var.process_content_name
  }

  lifecycle {
    ignore_changes = [
      source_archive_bucket,
    ]
  }
}

resource "google_cloud_scheduler_job" "get_content" {
  name        = var.get_content_name
  description = "Message by Terraform"
  schedule    = "*/30 * * * *"

  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.get_content.id
    data       = base64encode("starting get_content...")
  }
}
