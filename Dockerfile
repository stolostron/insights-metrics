FROM registry.ci.openshift.org/stolostron/builder:go1.24-linux AS builder

WORKDIR /go/src/github.com/stolostron/insights-metrics
COPY . .
RUN CGO_ENABLED=1 go build -trimpath -o insights-metrics main.go

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest

ENV VCS_REF="$VCS_REF" \
    USER_UID=1001

COPY --from=builder /go/src/github.com/stolostron/insights-metrics/insights-metrics /bin

EXPOSE 3031
USER ${USER_UID}
ENTRYPOINT ["/bin/insights-metrics"]
