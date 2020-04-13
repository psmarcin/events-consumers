variable "get_jobs_name" {
  type = string
  default = "get_jobs"
}

resource "google_pubsub_topic" "get_jobs" {
  name = var.get_jobs_name
}

resource "google_storage_bucket" "get_jobs" {
  name = var.get_jobs_name
}

resource "google_storage_bucket_object" "get_jobs" {
  name   = "${var.get_jobs_name}_${uuid()}.zip"
  bucket = google_storage_bucket.get_jobs.name
  source = "./../dist/job.zip"

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}

resource "google_cloudfunctions_function" "get_jobs" {
  name        = var.get_jobs_name
  description = "Managed by Terraform"
  runtime     = "go111"

  available_memory_mb   = 128
  source_archive_bucket = google_storage_bucket.get_jobs.name
  source_archive_object = google_storage_bucket_object.get_jobs.name
  event_trigger {
    event_type = "google.pubsub.topic.publish"
    resource   = google_pubsub_topic.get_jobs.name
  }
  timeout               = 10
  entry_point           = "GetJobs"
  service_account_email = "root-481@events-consumer.iam.gserviceaccount.com"

  environment_variables = {
    GET_CONTENT_TOPIC_ID = var.get_content_name
  }

  lifecycle {
    ignore_changes = [
      source_archive_bucket,
    ]
  }
}

resource "google_cloud_scheduler_job" "get_jobs" {
  name        = var.get_jobs_name
  description = "Message by Terraform"
  schedule    = "*/10 * * * *"

  pubsub_target {
    # topic.id is the topic's full resource name.
    topic_name = google_pubsub_topic.get_jobs.id
    data       = base64encode("{}")
  }
}
