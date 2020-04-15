FROM alpine:3.6

RUN apk add --no-cache ca-certificates
ADD cloudlycke-cloud-controller-manager /bin/
CMD ["/bin/cloudlycke-cloud-controller-manager"]