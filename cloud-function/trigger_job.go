package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

const (
	GrpcBuildTag    = "grpc-server-build"
	GrpcRenderTag   = "grpc-server-render"
	GrpcTagImageTag = "grpc-server-tag"
	HttpBuildTag    = "http-server-build"
	HttpRenderTag   = "http-server-render"
	HttpTagImageTag = "http-server-tag"

	ServiceDeployTag = "service-deploy"
	Status           = "SUCCESS"
	DevEnv           = "dev"
	ProdEnv          = "prod"
	ZoneSub          = "_ZONE"
	ClusterSub       = "_CLUSTER"
	ServiceSub       = "_SERVICE"
	EnvSub           = "_ENV"
	GrpcService      = "GRPC"
	HttpService      = "HTTP"
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
	name := j.Tag[0]
	env := j.Tag[1]
	var triggerId, repo string
	subs := make(map[string]string)
	projectId := "cloud-build-delivery-mad"

	if j.Status == Status {
		log.Println("ready to trigger next cloud build job ")
		switch name {
		case GrpcBuildTag:
			log.Printf("grpc %v render job  in queue", env)
			triggerId = "83babdf5-6538-47b3-b210-d2389c275907"
			repo = "mad-demo-grpc-service"
		case GrpcRenderTag:
			log.Printf("deploy %v in queue", env)
			triggerId = "0c9895dc-fa3b-4679-9c26-c770c9bb0924"
			repo = "mad-demo-config"
			subs[ServiceSub] = GrpcService
			if env == ProdEnv {
				subs[EnvSub] = ProdEnv
				subs[ZoneSub] = "us-east4-a"
				subs[ClusterSub] = "mad-demo-prod"
			}
		case GrpcTagImageTag:
			log.Printf("grpc %v render job  in queue", env)
			triggerId = "83babdf5-6538-47b3-b210-d2389c275907"
			repo = "mad-demo-grpc-service"
			subs[ServiceSub] = GrpcService
			subs[ZoneSub] = "us-east4-a"
			subs[ClusterSub] = "mad-demo-prod"
			subs[EnvSub] = ProdEnv
		case HttpBuildTag:
			log.Printf("http %v render job  in queue", env)
			triggerId = "1b9409fa-b90e-4c89-933a-111a4a319769"
			repo = "mad-demo-http-service"
		case HttpRenderTag:
			log.Println("deploy in queue ")
			triggerId = "0c9895dc-fa3b-4679-9c26-c770c9bb0924"
			repo = "mad-demo-config"
			subs[ServiceSub] = HttpService
			if env == ProdEnv {
				subs[EnvSub] = ProdEnv
				subs[ZoneSub] = "us-east4-a"
				subs[ClusterSub] = "mad-demo-prod"
			}
		case HttpTagImageTag:
			log.Printf("deploy %v in queue", env)
			triggerId = "1b9409fa-b90e-4c89-933a-111a4a319769"
			repo = "mad-demo-http-service"
			subs[ServiceSub] = HttpService
			subs[ZoneSub] = "us-east4-a"
			subs[ClusterSub] = "mad-demo-prod"
			subs[EnvSub] = ProdEnv
		case ServiceDeployTag:
			log.Printf("deploy %v completed", env)
			if env == DevEnv {
				if j.Tag[2] == GrpcService {
					triggerId = "6aac1818-b488-4875-be04-89c581c8c5fb"
					repo = "mad-demo-grpc-service"
				} else if j.Tag[2] == HttpService {
					triggerId = "29b3b817-f0b6-4597-9136-07b081c690e0"
					repo = "mad-demo-http-service"

				} else {
					return fmt.Errorf("unknow service")
				}
				subs[EnvSub] = ProdEnv
				break
			}
			return nil
		default:
			log.Println("nothing to enqueue ")
			return nil
		}

	} else {
		switch name {
		case GrpcBuildTag:
			log.Printf("grpc build job %v in %v", j.Status, env)
		case GrpcRenderTag:
			log.Printf("grpc render job %v in %v", j.Status, env)
		case GrpcTagImageTag:
			log.Printf("grpc tag image job %v in %v", j.Status, env)
		case HttpBuildTag:
			log.Printf("http build job  %v in %v", j.Status, env)
		case HttpRenderTag:
			log.Printf("http render job %v in %v", j.Status, env)
		case HttpTagImageTag:
			log.Printf("http tag image %v in %v", j.Status, env)
		case ServiceDeployTag:
			log.Printf("deploy job in %v %v ", env, j.Status)
		default:
			log.Printf("unknown job, j.Status: %q", j.Status)
		}
		return nil
	}

	return triggerCloudbuildJob(projectId, triggerId, repo, subs)
}

func triggerCloudbuildJob(projectId, triggerId, repo string, subs map[string]string) error {
	ctx := context.Background()
	c, err := cloudbuild.NewClient(ctx)
	if err != nil {
		return err
	}
	req := &cloudbuildpb.RunBuildTriggerRequest{
		ProjectId: projectId,
		TriggerId: triggerId,
		Source: &cloudbuildpb.RepoSource{
			RepoName:      repo,
			Substitutions: subs,
		},
	}
	_, err = c.RunBuildTrigger(ctx, req)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	//testData := []byte(`{"id":"6c6373ec-31bd-4eb7-8307-114584b77017","projectId":"cloud-build-delivery-mad","status":"WORKING","source":{"storageSource":{"bucket":"cloud-build-delivery-mad_cloudbuild","object":"source/1586207276.353547-ab1fff69c2704a829556a02a4aca64c0.tgz","generation":"1586207276771710"}},"steps":[{"name":"gcr.io/cloud-build-delivery-mad/skaffold:alpha","args":["build","-f=skaffold.yaml"]}],"createTime":"2020-04-06T21:07:57.077736329Z","startTime":"2020-04-06T21:07:58.356301014Z","timeout":"300s","logsBucket":"gs://416281961765.cloudbuild-logs.googleusercontent.com","sourceProvenance":{"resolvedStorageSource":{"bucket":"cloud-build-delivery-mad_cloudbuild","object":"source/1586207276.353547-ab1fff69c2704a829556a02a4aca64c0.tgz","generation":"1586207276771710"}},"options":{"substitutionOption":"ALLOW_LOOSE","logging":"LEGACY","env":["CLOUDSDK_COMPUTE_ZONE=us-east4-a","CLOUDSDK_CONTAINER_CLUSTER=cbd-mad-demo","ENV=dev"]},"logUrl":"https://console.cloud.google.com/cloud-build/builds/6c6373ec-31bd-4eb7-8307-114584b77017?project=416281961765","substitutions":{"_CLUSTER":"cbd-mad-demo","_ENV":"dev","_ZONE":"us-east4-a"},"tags":["grpc-server-build","dev"]}	:"1585785077233564"},"fileHashes":{"gs://416281961765.cloudbuild-source.googleusercontent.com/d73621656496ae854350b3233733dfe5fb275327-a45439dc-d083-4394-a64c-78f39d298e33.tar.gz#1585785077233564":{"fileHash":[{"type":"MD5","value":"JTpcim7hWtCilsWuSMgkuQ=="}]}}},"buildTriggerId":"9c60966a-d9ef-4915-910c-d5aa4583fd0c","options":{"substitutionOption":"ALLOW_LOOSE","logging":"LEGACY","env":["CLOUDSDK_COMPUTE_ZONE=us-east4-a","CLOUDSDK_CONTAINER_CLUSTER=cbd-mad-demo","ENV=dev"]},"logUrl":"https://console.cloud.google.com/cloud-build/builds/59337bcb-ce3a-4a49-a2e4-04b8b6c69d35?project=416281961765","substitutions":{"BRANCH_NAME":"master","COMMIT_SHA":"d73621656496ae854350b3233733dfe5fb275327","REPO_NAME":"mad-demo-grpc-service","REVISION_ID":"d73621656496ae854350b3233733dfe5fb275327","SHORT_SHA":"d736216","_ENV":"dev"},"tags":["grpc-server-build","trigger-9c60966a-d9ef-4915-910c-d5aa4583fd0c"],"timing":{"BUILD":{"startTime":"2020-04-01T23:51:23.177161570Z","endTime":"2020-04-01T23:52:39.461768573Z"},"FETCHSOURCE":{"startTime":"2020-04-01T23:51:19.433275125Z","endTime":"2020-04-01T23:51:23.177116759Z"}}}`)
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
