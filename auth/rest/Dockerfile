FROM phusion/baseimage:latest

CMD ["/sbin/my_init"]

#+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
RUN echo 'PS1="\[\033[36m\]\u\[\033[m\]@\[\033[95;1m\]psd/golang-auth-test:\[\033[34m\]\w\[\033[m\]\$ "' >> ~/.bashrc
RUN TERM=xterm
RUN DEBIAN_FRONTEND=noninteractive
#+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

RUN DEBIAN_FRONTEND=noninteractive \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        software-properties-common \
        apt-utils \
        git \
        wget \
        vim \
        curl \
		g++ \
		gcc \
		libc6-dev \
		make \
		sqlite3 \
        libsqlite3-dev \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.6.2
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 e40c36ae71756198478624ed1bb4ce17597b3c19d243f3f0899bb5740d56212a

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /var/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ENV ROOT_SRC /var/www/src/github.com/madmis
ENV APP_SRC "$ROOT_SRC/rest"
# Copy the local package files to the container's workspace.
ADD ./ "$APP_SRC"

RUN cd "$APP_SRC" && go get ./
RUN cd "$APP_SRC" && go install

RUN mkdir /etc/service/application
ADD ./build/app/run.sh /etc/service/application/run
RUN chmod +x /etc/service/application/run

COPY ./auth.db "$GOPATH/bin"

WORKDIR $GOPATH