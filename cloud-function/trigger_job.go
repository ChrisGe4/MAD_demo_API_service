package main

import (
	"context"
	"encoding/json"
	"log"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

const (
	GrpcBuildTag     = "grpc-server-build"
	GrpcRenderTag    = "grpc-server-render"
	HttpBuildTag     = "http-server-build"
	HttpRenderTag    = "http-server-render"
	ServiceDeployTag = "service-deploy"
	Status           = "SUCCESS"
)

type PubSubMessage struct {
	Data []byte `json:"data"`
}
type Job struct {
	Tag    []string `json:"tags"`
	Status string   `json:"status"`
}

// HelloPubSub consumes a Pub/Sub message.
func TriggerJob(ctx context.Context, m PubSubMessage) error {
	var j Job
	if err := json.Unmarshal(m.Data, &j); err != nil {
		return err
	}
	tag := j.Tag[0]
	var triggerId, repo string
	projectId := "cloud-build-delivery-mad"

	if j.Status == Status {
		log.Println("ready to trigger next cloud build job ")
		switch tag {
		case GrpcBuildTag:
			log.Println("grpc render job in queue ")
			triggerId = "83babdf5-6538-47b3-b210-d2389c275907"
			repo = "mad-demo-grpc-service"
		case GrpcRenderTag:
			log.Println("deploy in queue ")
			triggerId = "0c9895dc-fa3b-4679-9c26-c770c9bb0924"
			repo = "mad-demo-config"
		case HttpBuildTag:
			log.Println("http render job in queue ")
			triggerId = "1b9409fa-b90e-4c89-933a-111a4a319769"
			repo = "mad-demo-http-service"
		case HttpRenderTag:
			log.Println("deploy in queue ")
			triggerId = "0c9895dc-fa3b-4679-9c26-c770c9bb0924"
			repo = "mad-demo-config"
		case ServiceDeployTag:
			log.Println("deploy completed")
			return nil
		default:
			log.Println("nothing to enqueue ")
			return nil
		}

	} else {
		switch tag {
		case GrpcBuildTag:
			log.Printf("grpc build job %v ", j.Status)
		case GrpcRenderTag:
			log.Printf("grpc render job %v", j.Status)
		case HttpBuildTag:
			log.Printf("http build job  %v", j.Status)
		case HttpRenderTag:
			log.Printf("http render job %v", j.Status)
		case ServiceDeployTag:
			log.Printf("deploy job %v", j.Status)
		default:
			log.Printf("unknown job, tag: %q", tag)
		}
		return nil
	}

	return triggerCloudbuildJob(projectId, triggerId, repo)
}

func triggerCloudbuildJob(projectId, triggerId, repo string) error {
	ctx := context.Background()
	c, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return err
	}
	req := &cloudbuildpb.RunBuildTriggerRequest{
		ProjectId: projectId,
		TriggerId: triggerId,
		Source: &cloudbuildpb.RepoSource{
			RepoName: repo,
		},
	}
	_, err = c.RunBuildTrigger(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//testData := []byte(`{"id":"59337bcb-ce3a-4a49-a2e4-04b8b6c69d35","projectId":"cloud-build-delivery-mad","status":"SUCCESS","source":{"storageSource":{"bucket":"416281961765.cloudbuild-source.googleusercontent.com","object":"d73621656496ae854350b3233733dfe5fb275327-a45439dc-d083-4394-a64c-78f39d298e33.tar.gz"}},"steps":[{"name":"gcr.io/cloud-build-delivery-mad/skaffold:alpha","args":["build","-f=skaffold.yaml"],"id":"build-grpc-server","timing":{"startTime":"2020-04-01T23:51:23.993610304Z","endTime":"2020-04-01T23:52:39.461735024Z"},"pullTiming":{"startTime":"2020-04-01T23:51:23.993610304Z","endTime":"2020-04-01T23:51:27.331610352Z"},"status":"SUCCESS"}],"results":{"buildStepImages":[""],"buildStepOutputs":[""]},"createTime":"2020-04-01T23:51:17.420407503Z","startTime":"2020-04-01T23:51:18.802219616Z","finishTime":"2020-04-01T23:52:40.442709Z","timeout":"300s","logsBucket":"gs://416281961765.cloudbuild-logs.googleusercontent.com","sourceProvenance":{"resolvedStorageSource":{"bucket":"416281961765.cloudbuild-source.googleusercontent.com","object":"d73621656496ae854350b3233733dfe5fb275327-a45439dc-d083-4394-a64c-78f39d298e33.tar.gz","generation":"1585785077233564"},"fileHashes":{"gs://416281961765.cloudbuild-source.googleusercontent.com/d73621656496ae854350b3233733dfe5fb275327-a45439dc-d083-4394-a64c-78f39d298e33.tar.gz#1585785077233564":{"fileHash":[{"type":"MD5","value":"JTpcim7hWtCilsWuSMgkuQ=="}]}}},"buildTriggerId":"9c60966a-d9ef-4915-910c-d5aa4583fd0c","options":{"substitutionOption":"ALLOW_LOOSE","logging":"LEGACY","env":["CLOUDSDK_COMPUTE_ZONE=us-east4-a","CLOUDSDK_CONTAINER_CLUSTER=cbd-mad-demo","ENV=dev"]},"logUrl":"https://console.cloud.google.com/cloud-build/builds/59337bcb-ce3a-4a49-a2e4-04b8b6c69d35?project=416281961765","substitutions":{"BRANCH_NAME":"master","COMMIT_SHA":"d73621656496ae854350b3233733dfe5fb275327","REPO_NAME":"mad-demo-grpc-service","REVISION_ID":"d73621656496ae854350b3233733dfe5fb275327","SHORT_SHA":"d736216","_ENV":"dev"},"tags":["grpc-server-build","trigger-9c60966a-d9ef-4915-910c-d5aa4583fd0c"],"timing":{"BUILD":{"startTime":"2020-04-01T23:51:23.177161570Z","endTime":"2020-04-01T23:52:39.461768573Z"},"FETCHSOURCE":{"startTime":"2020-04-01T23:51:19.433275125Z","endTime":"2020-04-01T23:51:23.177116759Z"}}}`)
	testData := []byte(`{"status":"SUCCESS","tags":["grpc-server-build","trigger-9c60966a-d9ef-4915-910c-d5aa4583fd0c"]}`)
	var j Job
	if err := json.Unmarshal(testData, &j); err != nil {
		log.Fatal(err)
	}
	log.Println(j)
	log.Println(j.Status)
	log.Println(j.Tag[0])

	ctx := context.Background()
	c, err := cloudbuild.NewClient(ctx)
	if err != nil {
		// return err
	}
	req := &cloudbuildpb.RunBuildTriggerRequest{
		ProjectId: "cloud-build-delivery-mad",
		TriggerId: "83babdf5-6538-47b3-b210-d2389c275907",
		Source: &cloudbuildpb.RepoSource{
			RepoName: "mad-demo-grpc-service",
		},
	}
	op, err := c.RunBuildTrigger(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := op.Wait(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: Use resp.
	log.Println(resp)
}
