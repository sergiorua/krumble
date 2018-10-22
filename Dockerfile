FROM golang:alpine as builder

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go

RUN apk add --no-cache \
	ca-certificates \
	--virtual .build-deps \
	git \
	gcc \
	libc-dev \
	libgcc \
  make \
  wget

COPY . /go/src/github.com/sergiorua/krumble
WORKDIR /go/src/github.com/sergiorua/krumble
ARG BUILD_TAGS="netgo osusergo"
RUN make VK_BUILD_TAGS="${BUILD_TAGS}" build
RUN cp bin/krumble /usr/bin/krumble

RUN wget https://github.com/kubernetes/kops/releases/download/1.10.0/kops-linux-amd64 -O kops-linux-amd64 -o /dev/null && \
    wget https://releases.hashicorp.com/terraform/0.11.8/terraform_0.11.8_linux_amd64.zip -O terraform_0.11.8_linux_amd64.zip -o /dev/null && \
    wget https://storage.googleapis.com/kubernetes-helm/helm-v2.11.0-linux-amd64.tar.gz -O helm-v2.11.0-linux-amd64.tar.gz -o /dev/null && \
    wget https://github.com/kubeless/kubeless/releases/download/v1.0.0-alpha.8/kubeless_linux-amd64.zip -O kubeless_linux-amd64.zip -o /dev/null && \
    wget https://github.com/roboll/helmfile/releases/download/v0.40.1/helmfile_linux_amd64 -O /usr/bin/helmfile -o /dev/null && \
    wget https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v0.3.0/heptio-authenticator-aws_0.3.0_linux_amd64 -O /usr/bin/heptio-authenticator-aws -o /dev/null && \
    ls -l && tar zxvf helm-v2.11.0-linux-amd64.tar.gz && \
    cp linux-amd64/helm linux-amd64/tiller /usr/bin && \
    unzip terraform_0.11.8_linux_amd64.zip -d /usr/bin && \
    unzip kubeless_linux-amd64.zip && \
    cp bundles/kubeless_linux-amd64/kubeless /usr/bin/kubeless && \
    mv kops-linux-amd64 /usr/bin/kops && \
    chmod 755 /usr/bin/kops /usr/bin/terraform /usr/bin/helm /usr/bin/tiller /usr/bin/kubeless /usr/bin/helmfile /usr/bin/heptio-authenticator-aws

FROM scratch
COPY --from=builder /usr/bin/krumble /usr/bin/krumble
COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs
COPY --from=builder /usr/bin/kops /usr/bin/kops
COPY --from=builder /usr/bin/terraform /usr/bin/terraform
COPY --from=builder /usr/bin/helm /usr/bin/helm
COPY --from=builder /usr/bin/tiller /usr/bin/tiller
COPY --from=builder /usr/bin/kubeless /usr/bin/kubeless
COPY --from=builder /usr/bin/helmfile /usr/bin/helmfile
COPY --from=builder /usr/bin/heptio-authenticator-aws /usr/bin/heptio-authenticator-aws
ENTRYPOINT [ "/usr/bin/krumble" ]
CMD [ "--help" ]
