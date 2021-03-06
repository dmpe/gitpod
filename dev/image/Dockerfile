# Copyright (c) 2020 TypeFox GmbH. All rights reserved.
# Licensed under the GNU Affero General Public License (AGPL).
# See License-AGPL.txt in the project root for license information.

FROM gitpod/workspace-full-vnc:latest

ENV TRIGGER_REBUILD 1

USER root

### cloud_sql_proxy ###
ARG CLOUD_SQL_PROXY=/usr/local/bin/cloud_sql_proxy
RUN curl -fsSL https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 > $CLOUD_SQL_PROXY \
    && chmod +x $CLOUD_SQL_PROXY

### Docker client ###
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - \
    # 'cosmic' not supported
    && add-apt-repository -yu "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable" \
    && apt-get install -yq docker-ce-cli=5:18.09.0~3-0~ubuntu-bionic \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

### Helm3 ###
RUN mkdir -p /tmp/helm/ \
    && curl -fsSL https://get.helm.sh/helm-v3.0.3-linux-amd64.tar.gz | tar -xzvC /tmp/helm/ --strip-components=1 \
    && cp /tmp/helm/helm /usr/local/bin/helm \
    && cp /tmp/helm/helm /usr/local/bin/helm3 \
    && rm -rf /tmp/helm/

### kubectl ###
RUN curl -fsSL https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add - \
    # really 'xenial'
    && add-apt-repository -yu "deb https://apt.kubernetes.io/ kubernetes-xenial main" \
    && apt-get install -yq kubectl=1.13.0-00 \
    && kubectl completion bash > /usr/share/bash-completion/completions/kubectl \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

RUN curl -o /usr/bin/kubectx https://raw.githubusercontent.com/ahmetb/kubectx/master/kubectx && chmod +x /usr/bin/kubectx \
 && curl -o /usr/bin/kubens  https://raw.githubusercontent.com/ahmetb/kubectx/master/kubens  && chmod +x /usr/bin/kubens

### MySQL client ###
RUN apt-get update && apt-get install -yq \
    mysql-client  \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

# golangci-lint
RUN cd /usr/local && curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.16.0

# leeway
ENV LEEWAY_NESTED_WORKSPACE=true
RUN cd /usr/bin && curl -L https://github.com/TypeFox/leeway/releases/download/v0.1.0/leeway_0.1.0_Linux_x86_64.tar.gz | tar xz

# dazzle
RUN cd /usr/bin && curl -L https://github.com/32leaves/dazzle/releases/download/v0.0.3/dazzle_0.0.3_Linux_x86_64.tar.gz | tar xz

# werft CLI
RUN cd /usr/bin && curl -L https://github.com/csweichel/werft/releases/download/v0.0.4/werft-client-linux-amd64.tar.gz | tar xz && mv werft-client-linux-amd64 werft

# yq - jq for YAML files
RUN cd /usr/bin && curl -L https://github.com/mikefarah/yq/releases/download/2.4.0/yq_linux_amd64 > yq && chmod +x yq

### Protobuf
RUN set -ex && \
    tmpdir=$(mktemp -d) && \
    curl -L -o $tmpdir/protoc.zip https://github.com/protocolbuffers/protobuf/releases/download/v3.7.1/protoc-3.7.1-linux-x86_64.zip && \
    mkdir -p /usr/lib/protoc && cd /usr/lib/protoc && unzip $tmpdir/protoc.zip && \
    chmod -R 755 /usr/lib/protoc/include/google && \
    ln -s /usr/lib/protoc/bin/* /usr/bin && \
    rm -rf $tmpdir

### Telepresence ###
RUN curl -fsSL https://packagecloud.io/datawireio/telepresence/gpgkey | apt-key add - \
    # 'cosmic' not supported
    && add-apt-repository -yu "deb https://packagecloud.io/datawireio/telepresence/ubuntu/ bionic main" \
    # 0.95 (current at the time of this commit) is broken
    && apt-get install -yq \
    iproute2 \
    iptables \
    net-tools \
    telepresence=0.94 \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/*

### Toxiproxy CLI
RUN curl -L -o /usr/bin/toxiproxy https://github.com/Shopify/toxiproxy/releases/download/v2.1.4/toxiproxy-cli-linux-amd64 \
    && chmod +x /usr/bin/toxiproxy

USER gitpod

# Fix node version we develop against
ARG GITPOD_NODE_VERSION=12.18.3
RUN bash -c ". .nvm/nvm.sh \
        && OLD_NODE=\$(node --version) \
        && nvm install $GITPOD_NODE_VERSION \
        && nvm uninstall \$OLD_NODE || : \
        && npm install -g typescript yarn"
ENV PATH=/home/gitpod/.nvm/versions/node/v${GITPOD_NODE_VERSION}/bin:$PATH

# Fix the Go version
ENV GOPATH=$HOME/go-packages
ENV GOROOT=$HOME/go
RUN go get golang.org/dl/go1.15 && \
    go1.15 download && \
    mv $(which go1.15) $(which go)
ENV GOPATH=/workspace/go \
    PATH=/workspace/go/bin:$PATH

## Register leeway autocompletion in bashrc
RUN bash -c "echo . \<\(leeway bash-completion\) >> ~/.bashrc"

# Install pulumi
RUN wget -qO- https://get.pulumi.com/releases/sdk/pulumi-v1.9.1-linux-x64.tar.gz | tar xvzC . \
    && mkdir -p /home/gitpod/.pulumi/bin \
    && mv pulumi/* /home/gitpod/.pulumi/bin \
    && echo "export PATH=\$PATH:/home/gitpod/.pulumi/bin" >> /home/gitpod/.bashrc \
    && rm -Rf ./pulumi

### Google Cloud ###
# not installed via repository as then 'docker-credential-gcr' is not available
ARG GCS_DIR=/opt/google-cloud-sdk
ENV PATH=$GCS_DIR/bin:$PATH
RUN sudo chown gitpod: /opt \
    && mkdir $GCS_DIR \
    && curl -fsSL https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-245.0.0-linux-x86_64.tar.gz \
    | tar -xzvC /opt \
    && /opt/google-cloud-sdk/install.sh --quiet --usage-reporting=false --bash-completion=true \
    --additional-components docker-credential-gcr alpha beta \
    # needed for access to our private registries
    && docker-credential-gcr configure-docker

# Install tools for gsutil
RUN sudo apt-get update -qq \
    && sudo apt-get install -qq -y \
        gcc \
        python-dev \
        python-setuptools \
    && sudo apt-get clean && sudo rm -rf /var/lib/apt/lists/* /tmp/*
RUN bash -c "pip uninstall crcmod; pip install --no-cache-dir -U crcmod"

### gitpod-core specific gcloud/kubectl config
# Copy GCloud default config that points to gitpod-dev
ARG GCLOUD_CONFIG_DIR=/home/gitpod/.config/gcloud
COPY --chown=gitpod gcloud-default-config $GCLOUD_CONFIG_DIR/configurations/config_default

# Set kubeconfig file for dev cluster, using GCloud Application Default Credentials (ADC) as auth provider
ARG KUBE_CONFIG_PATH=/home/gitpod/.kube/config
COPY --chown=gitpod kubeconfig.yaml $KUBE_CONFIG_PATH

# Set Application Default Credentials (ADC) based on user-provided env var
RUN echo ". /workspace/gitpod/scripts/setup-google-adc.sh" >> ~/.bashrc

## install Terraform kubectl provider
# https://gavinbunney.github.io/terraform-provider-kubectl/docs/provider.html
RUN mkdir -p ~/.terraform.d/plugins && \
    curl -Lo ~/.terraform.d/plugins/terraform-provider-kubectl https://github.com/gavinbunney/terraform-provider-kubectl/releases/download/v1.5.1/terraform-provider-kubectl-linux-amd64 && \
    chmod +x ~/.terraform.d/plugins/terraform-provider-kubectl

ENV DB_HOST=localhost

ENV LEEWAY_WORKSPACE_ROOT=/workspace/gitpod
ENV LEEWAY_REMOTE_CACHE_BUCKET=gitpod-core-leeway-cache-branch

### AWS Cli ###
RUN pip install awscli

# Install aws-iam-authenticator
RUN sudo curl -o aws-iam-authenticator "https://amazon-eks.s3-us-west-2.amazonaws.com/1.13.7/2019-06-11/bin/linux/amd64/aws-iam-authenticator" \
    && sudo chmod +x ./aws-iam-authenticator \
    && sudo mkdir -p $HOME/.aws-iam \
    && sudo cp ./aws-iam-authenticator $HOME/.aws-iam/aws-iam-authenticator

# Install Terraform
ARG RELEASE_URL="https://releases.hashicorp.com/terraform/0.13.0/terraform_0.13.0_linux_amd64.zip"
RUN mkdir -p ~/.terraform \
    && cd ~/.terraform \
    && wget ${RELEASE_URL} \
    && unzip *.zip \
    && rm -f *.zip \
    && printf "terraform -install-autocomplete\n" >>~/.bashrc

# Install GraphViz to help debug terraform scripts
RUN sudo apt-get update -qq \
    && sudo apt-get install -yq \
        graphviz

ENV PATH=$PATH:$HOME/.aws-iam:$HOME/.terraform
