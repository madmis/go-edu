FROM phusion/baseimage:latest

CMD ["/sbin/my_init"]

#+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++
RUN echo 'PS1="\[\033[36m\]\u\[\033[m\]@\[\033[95;1m\]psd/golang-auth-test:\[\033[34m\]\w\[\033[m\]\$ "' >> ~/.bashrc
RUN TERM=xterm
RUN DEBIAN_FRONTEND=noninteractive
#+++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

RUN DEBIAN_FRONTEND=noninteractive \
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 2EA8F35793D8809A \
    && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys E5267A6C \
    && add-apt-repository ppa:nginx/stable \
    && apt-key adv --keyserver keyserver.ubuntu.com --verbose --refresh-keys \
    && sudo apt-key update \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
        software-properties-common \
        apt-utils \
        git \
        wget \
        vim \
        curl \
        nginx \
	&& rm -rf /var/lib/apt/lists/*

ENV ROOT_SRC /var/www/src/github.com/madmis
ENV APP_SRC "$ROOT_SRC/rest-docs"
# Copy the local package files to the container's workspace.
RUN mkdir -p /var/www
ADD ./ "$APP_SRC"
RUN ln -s "$APP_SRC" /var/www/rest-docs

# Configure NGINX
RUN rm -v /etc/nginx/nginx.conf
ADD ./build/nginx/nginx.conf /etc/nginx/nginx.conf

RUN rm -v /etc/nginx/sites-available/default
ADD ./build/nginx/app.conf /etc/nginx/sites-available/default

RUN mkdir /etc/service/nginx
ADD ./build/nginx/nginx.sh /etc/service/nginx/run
RUN chmod +x /etc/service/nginx/run

RUN usermod -u 1000 www-data
RUN usermod -a -G root www-data

WORKDIR /var/www/rest-docs