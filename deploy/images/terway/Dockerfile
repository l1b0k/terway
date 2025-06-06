# syntax=docker/dockerfile:1-labs
ARG TERWAY_POLICY_IMAGE=registry-cn-zhangjiakou.ack.aliyuncs.com/acs/terway:policy-f992d305@sha256:d20c35c41af4f906ce0b169ac0dd06528172344c02411cd40bf6b35f77701e1f
ARG UBUNTU_IMAGE=registry.cn-hangzhou.aliyuncs.com/acs/ubuntu:22.04-update
ARG CILIUM_IPROUTE2_IMAGE=quay.io/cilium/cilium-iproute2:3570d58349efb2d6b0342369a836998c93afd291@sha256:1abcd7a5d2117190ab2690a163ee9cd135bc9e4cf8a4df662a8f993044c79342
ARG CILIUM_LLVM_IMAGE=quay.io/cilium/cilium-llvm:9f1bfe736009afb1fbb562718bbc42ea07d37d8e@sha256:a666a7a01a2dc610c3ab6e32f25ca5e294201f3cbbc01f233320c527955deee3
ARG CILIUM_BPFTOOL_IMAGE=quay.io/cilium/cilium-bpftool:0db3a73729ceb42e947d826bb96a655be79e5317@sha256:de23c9546c4eafab33f75d6f5d129947bbbafc132dbd113c0cecc9a61929e6b0
ARG CILIUM_IPTABLES_IMAGE=quay.io/cilium/iptables:67f517af50e18f64cd12625021f1c39246bb4f92@sha256:d075f03e89aacf51908346ec8ed5d251b8d3ad528ce30a710fcd074cdf91f11d

FROM --platform=$TARGETPLATFORM ${TERWAY_POLICY_IMAGE} AS policy-dist
FROM --platform=$TARGETPLATFORM ${CILIUM_LLVM_IMAGE} AS llvm-dist
FROM --platform=$TARGETPLATFORM ${CILIUM_BPFTOOL_IMAGE} AS bpftool-dist
FROM --platform=$TARGETPLATFORM ${CILIUM_IPROUTE2_IMAGE} AS iproute2-dist
FROM --platform=$TARGETPLATFORM ${CILIUM_IPTABLES_IMAGE} AS iptables-dist

FROM --platform=$BUILDPLATFORM golang:1.24.2 AS builder
ARG GOPROXY
ARG TARGETOS
ARG TARGETARCH
ENV GOPROXY=$GOPROXY
WORKDIR /go/src/github.com/AliyunContainerService/terway/
COPY --parents go.mod go.sum .git cmd daemon deploy deviceplugin pkg plugin rpc tests types ./
RUN cd cmd/terway && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags default_build \
    -ldflags \
    "-s -w -X \"github.com/AliyunContainerService/terway/pkg/version.gitCommit=`git rev-parse HEAD`\" \
    -X \"github.com/AliyunContainerService/terway/pkg/version.buildDate=`date -u +'%Y-%m-%dT%H:%M:%SZ'`\" \
    -X \"github.com/AliyunContainerService/terway/pkg/version.gitVersion=`git describe --tags --match='v*' --abbrev=14`\" \
    -X \"github.com/AliyunContainerService/terway/pkg/aliyun/credential.kubernetesAlicloudIdentity=terwayd/`git rev-parse --short HEAD 2>/dev/null`\"" -o terwayd .
RUN cd plugin/terway && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags default_build -ldflags "-s -w" -o terway .
RUN cd cmd/terway-cli && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -tags default_build -ldflags "-s -w" -o terway-cli .

FROM --platform=$TARGETPLATFORM ${UBUNTU_IMAGE}
RUN apt-get update && apt-get install -y kmod libelf1 libmnl0 iptables nftables kmod curl ipset bash ethtool bridge-utils socat grep findutils jq conntrack iputils-ping && \
    apt-get purge --auto-remove && apt-get clean && rm -rf /var/lib/apt/lists/*

COPY --from=llvm-dist /usr/local/bin/clang /usr/local/bin/llc /bin/
COPY --from=bpftool-dist /usr/local /usr/local
COPY --from=iproute2-dist /usr/local /usr/local
COPY --from=iproute2-dist /usr/lib/libbpf* /usr/lib/
COPY --from=policy-dist /bin/calico-felix /bin/calico-felix
COPY --from=iptables-dist /iptables /iptables
RUN dpkg -i /iptables/*\.deb && rm -rf /iptables

COPY ../../../policy/policyinit.sh ./../../policy/uninstall_policy.sh ../../../hack/init.sh /bin/
COPY --from=policy-dist /tmp/install/ /
COPY --from=builder /go/src/github.com/AliyunContainerService/terway/cmd/terway/terwayd /go/src/github.com/AliyunContainerService/terway/plugin/terway/terway /go/src/github.com/AliyunContainerService/terway/cmd/terway-cli/terway-cli /usr/bin/
COPY ../../../hack/iptables-wrapper-installer.sh /iptables-wrapper-installer.sh
RUN /iptables-wrapper-installer.sh --no-sanity-check
RUN rm -f /usr/bin/cilium-operator-generic /usr/bin/cilium-health*
ENTRYPOINT ["/usr/bin/terwayd"]
