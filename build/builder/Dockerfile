FROM ubuntu:latest

# Install deps
RUN set -x; \
    apt-get update                                     \
 && apt-get install -y -q                              \
	clang \
	llvm-dev \
	libxml2-dev \
	uuid-dev \
  	libssl-dev \
	bash \
	patch \
	make \
	tar \
	xz-utils \
	bzip2 \
	gzip \
	sed \
	cpio \
	curl \
	git \
	ssh \
	rsync \
 && apt-get clean

# Instal mingw for windows cross compilation
RUN apt-get install -y -q mingw-w64

# Install golang
ARG go_version="1.11.3"
RUN  curl https://dl.google.com/go/go${go_version}.linux-amd64.tar.gz > go.tar.gz \
 && tar -xf go.tar.gz \
 && rm go.tar.gz \
 && rm -Rf /usr/local/go \
 && mv go /usr/local

ENV GOPATH "/go/"

ENV PATH "/usr/local/go/bin:${PATH}"

# Install OSx cross-tools

#Build arguments
ARG osxcross_repo="tpoechtrager/osxcross"
ARG osxcross_revision="a845375e028d29b447439b0c65dea4a9b4d2b2f6"
ARG darwin_sdk_version="10.10"
ARG darwin_osx_version_min="10.6"
ARG darwin_version="14"
ARG darwin_sdk_url="https://www.dropbox.com/s/yfbesd249w10lpc/MacOSX${darwin_sdk_version}.sdk.tar.xz"

# ENV available in docker image
ENV OSXCROSS_REPO="${osxcross_repo}"                   \
    OSXCROSS_REVISION="${osxcross_revision}"           \
    DARWIN_SDK_VERSION="${darwin_sdk_version}"         \
    DARWIN_VERSION="${darwin_version}"                 \
    DARWIN_OSX_VERSION_MIN="${darwin_osx_version_min}" \
    DARWIN_SDK_URL="${darwin_sdk_url}"

RUN mkdir -p "/tmp/osxcross"                                                                                   \
 && cd "/tmp/osxcross"                                                                                         \
 && curl -sLo osxcross.tar.gz "https://codeload.github.com/${OSXCROSS_REPO}/tar.gz/${OSXCROSS_REVISION}"  \
 && tar --strip=1 -xzf osxcross.tar.gz                                                                         \
 && rm -f osxcross.tar.gz                                                                                      \
 && curl -sLo tarballs/MacOSX${DARWIN_SDK_VERSION}.sdk.tar.xz                                                  \
             "${DARWIN_SDK_URL}"                \
 && yes "" | SDK_VERSION="${DARWIN_SDK_VERSION}" OSX_VERSION_MIN="${DARWIN_OSX_VERSION_MIN}"

RUN cd /tmp/osxcross \
 && UNATTENDED=1 ./build.sh
 
ENV PATH "/tmp/osxcross/target/bin:${PATH}"
ENV PATH "/go/bin:${PATH}"
