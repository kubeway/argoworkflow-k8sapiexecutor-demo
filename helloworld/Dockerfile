FROM registry.cn-huhehaote.aliyuncs.com/kubeway/alpine:3.14 as builder
WORKDIR /go/src/github.com/kubeway/argoworkflow-k8sapiexecutor-demo/helloworld
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -o helloworld

FROM registry.cn-huhehaote.aliyuncs.com/kubeway/alpine:3.14
COPY --from=builder /go/src/github.com/kubeway/argoworkflow-k8sapiexecutor-demo/helloworld/helloworld /helloworld
# Run the web service on container startup.
CMD ["/helloworld"]
