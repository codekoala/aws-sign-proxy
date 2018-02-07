FROM debian:8-slim
MAINTAINER Josh VanderLinden <codekoala@gmail.com>

ADD ./bin/aws-sign-proxy /bin/

ENTRYPOINT /bin/aws-sign-proxy
