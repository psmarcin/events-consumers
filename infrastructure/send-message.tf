variable "send_message_name" {
  type = string
  default = "send_message"
}

variable "telegram_api_key" {
  type = string
  default = ""
}

variable "telegram_channel_id" {
  type = string
  default = ""
}

resource "google_pubsub_topic" "send_message" {
  name = var.send_message_name
}

resource "google_storage_bucket" "send_message" {
  name = var.send_message_name
}

resource "google_storage_bucket_object" "send_message" {
  name   = "${var.send_message_name}_${uuid()}.zip"
  bucket = google_storage_bucket.send_message.name
  source = "./../dist/message.zip"

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}

resource "google_cloudfunctions_function" "send_message" {
  name        = var.send_message_name
  description = "Message by Terraform"
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.send_message.name
  source_archive_object = google_storage_bucket_object.send_message.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.send_message.name
  }
  timeout               = 10
  entry_point           = "Send"
  service_account_email = "root-481@events-consumer.iam.gserviceaccount.com"

  environment_variables = {
    TELEGRAM_API_KEY = var.telegram_api_key,
    TELEGRAM_CHANNEL_ID = var.telegram_channel_id,
  }

  lifecycle {
    ignore_changes = [
      source_archive_bucket,
    ]
  }
}
