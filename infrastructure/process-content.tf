resource "google_pubsub_topic" "process_content" {
  name = "process_content"
}

resource "google_storage_bucket" "process_content" {
  name = "process_content"
}

resource "google_storage_bucket_object" "process_content" {
  name   = "process_content_${uuid()}.zip"
  bucket = google_storage_bucket.process_content.name
  source = "./../dist/content.zip"

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}

resource "google_cloudfunctions_function" "process_content" {
  name        = "process_content"
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

  lifecycle {
    ignore_changes = [
      source_archive_bucket,
    ]
  }
}
