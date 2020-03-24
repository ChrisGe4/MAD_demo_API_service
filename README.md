# MAD_demo_API_service

This service is being used for preliminary research for Managed Application Delivery. 

**Complete**
  * Cloud Build setup
  * Helm Chart
  * Rest API    
  * gRPC

**In progress**
  * redis
  * helm for gRPC server 
  * skaffold


**Deploy GRPC service**
  * cd build/grpc
  * gcloud builds submit ../../
  * helm install grpc-server --debug --dry-run helm/grpc

    