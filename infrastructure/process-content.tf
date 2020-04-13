variable "process_content_name" {
  type = string
  default = "process_content"
}

resource "google_pubsub_topic" "process_content" {
  name = var.process_content_name
}

resource "google_storage_bucket" "process_content" {
  name = var.process_content_name
}

resource "google_storage_bucket_object" "process_content" {
  name   = "${var.process_content_name}_${uuid()}.zip"
  bucket = google_storage_bucket.process_content.name
  source = "./../dist/content.zip"

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}

resource "google_cloudfunctions_function" "process_content" {
  name        = var.process_content_name
  description = "Message by Terraform"
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.process_content.name
  source_archive_object = google_storage_bucket_object.process_content.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.process_content.name
  }
  timeout               = 10
  entry_point           = "Process"
  service_account_email = "root-481@events-consumer.iam.gserviceaccount.com"

  environment_variables = {
    SEND_MESSAGE_TOPIC_ID = var.send_message_name
  }

  lifecycle {
    ignore_changes = [
      source_archive_bucket,
    ]
  }
}
