package jobs

import (
	"context"
	"errors"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"

	"events-consumers/admin/pkg/config"
)

type Job struct {
	FirestoreClient *firestore.Client
}

func New() (Job, error) {
	job := Job{}
	bgCtx := context.Background()
	client, err := firestore.NewClient(bgCtx, config.C.CloudProject)
	if err != nil {
		return job, err
	}
	job.FirestoreClient = client

	return job, nil
}

type JobSerialized struct {
	Id        string    `json:"id"`
	Command   string    `json:"command"`
	Name      string    `json:"name"`
	Selector  string    `json:"selector"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (j *Job) List() ([]JobSerialized, error) {
	var jobs []JobSerialized

	ctx := context.Background()

	query := j.FirestoreClient.Collection(config.C.FirestoreCollectionId).OrderBy("updatedAt", firestore.Desc)

	iter := query.Documents(ctx)
	for {
		var job JobSerialized
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return jobs, err
		}
		doc.DataTo(&job)
		job.Id = doc.Ref.ID

		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (j *Job) Get(id string) (JobSerialized, error) {
	var job JobSerialized
	ctx := context.Background()

	result, err := j.FirestoreClient.Collection(config.C.FirestoreCollectionId).Doc(id).Get(ctx)
	log.Printf("result %+v", result.Exists())
	if err != nil {
		return job, err
	}
	if !result.Exists() {
		return job, errors.New("Job not found!")
	}
	result.DataTo(&job)
	job.Id = result.Ref.ID

	return job, nil
}

func (j *Job) Create(command, name, selector string) (JobSerialized, error) {
	var job JobSerialized
	ctx := context.Background()

	ref := j.FirestoreClient.Collection(config.C.FirestoreCollectionId).NewDoc()
	_, err := ref.Set(ctx, map[string]interface{}{
		"command":   command,
		"name":      name,
		"selector":  selector,
		"createdAt": time.Now(),
		"updatedAt": time.Now(),
	})
	if err != nil {
		return job, err
	}

	document, err := ref.Get(ctx)
	if err != nil {
		return job, err
	}

	document.DataTo(&job)
	job.Id = document.Ref.ID

	return job, nil
}

func (j *Job) Update(id, command, name, selector string) (JobSerialized, error) {
	var job JobSerialized
	ctx := context.Background()
	_, err := j.FirestoreClient.Collection(config.C.FirestoreCollectionId).Doc(id).Set(ctx, map[string]interface{}{
		"command":   command,
		"name":      name,
		"selector":  selector,
		"updatedAt": time.Now(),
	}, firestore.MergeAll)
	if err != nil {
		return job, err
	}

	job, err = j.Get(id)
	if err != nil {
		return job, err
	}
	log.Printf("after update %+v", job)
	return job, nil
}

func (j *Job) Delete(id string) error {
	ctx := context.Background()

	_, err := j.FirestoreClient.Collection(config.C.FirestoreCollectionId).Doc(id).Delete(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (j *Job) Close() error {
	return j.FirestoreClient.Close()
}
