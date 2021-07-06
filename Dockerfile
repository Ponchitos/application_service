FROM docker-registry.x5.ru/ubuntu
ARG GO_VER=1.16.2
ARG SDK_VERSION=commandlinetools-linux-6858069_latest.zip

ENV DEBIAN_FRONTEND=noninteractive
ENV ANDROID_SDK_ROOT=/android-sdk


ENV PATH=$ANDROID_SDK_ROOT/cmdline-tools/tools/latest/bin:$ANDROID_SDK_ROOT/cmdline-tools/tools/bin:$ANDROID_SDK_ROOT/build-tools/30.0.3:/usr/local/go/bin:$PATH

ENV JAVA_HOME=/usr/lib/jvm/java-8-openjdk-amd64/

ENV PATH=/usr/bin/python3.9:$PATH

RUN apt update -qq && apt install -qqy openssl openjdk-8-jdk python3 python3-pip wget curl unzip

RUN ln -s /usr/bin/python3.9 /usr/bin/python

RUN wget https://golang.org/dl/go${GO_VER}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VER}.linux-amd64.tar.gz && \
    rm go${GO_VER}.linux-amd64.tar.gz


RUN curl -sS https://dl.google.com/android/repository/${SDK_VERSION} -o /tmp/sdk.zip && \
    mkdir -p ${ANDROID_SDK_ROOT}/cmdline-tools/tools && \
    unzip -q -d /tmp/ /tmp/sdk.zip && rm /tmp/sdk.zip && \
    mv /tmp/cmdline-tools/* ${ANDROID_SDK_ROOT}/cmdline-tools/tools/ && rm -rf /tmp/cmdline-tools/ && \
    yes | sdkmanager --licenses && \
    sdkmanager --update && sdkmanager "platforms;android-30" "build-tools;30.0.3"

WORKDIR /opt/application-manager
COPY . .
RUN CGO_ENABLED=0 go build -o /opt/application-manager/application-manager server/cmd/seam/*.go

COPY config /opt/application-service/config
COPY migrations /opt/application-service/migrations
ENV TZ=Europe/Moscow
ENTRYPOINT ["/opt/application-service/application-service"]