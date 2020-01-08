resource "google_storage_bucket" "bucket" {
  name = "functions-message"
}

resource "google_storage_bucket_object" "archive" {
  name   = "index.zip"
  bucket = google_storage_bucket.bucket.name
  source = "./../dist/message.zip"
}

resource "google_cloudfunctions_function" "function" {
  name        = "message"
  description = "Message by Terraform"
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.archive.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.notification_topic.name
  }
  timeout               = 10
  entry_point           = "Consume"
  service_account_email = "root-481@events-consumer.iam.gserviceaccount.com"
}
