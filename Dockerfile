FROM plugins/base:linux-amd64

LABEL maintainer="Vickxxx <lianghuiming@live.com>" \
  org.label-schema.name="Drone Git Push Ext" \
  org.label-schema.vendor="vickxxx" \
  org.label-schema.schema-version="1.1e"

RUN sed -i "s@https://dl-cdn.alpinelinux.org/@https://mirrors.aliyun.com/@g" /etc/apk/repositories && \
  sed -i "s@http://nl.alpinelinux.org/@https://mirrors.aliyun.com/@g" /etc/apk/repositories && \
  apk update && \
  apk add --no-cache ca-certificates git openssh curl perl && \
  rm -rf /var/cache/apk/* 


COPY ./drone-git-push /bin/

ENTRYPOINT ["/bin/drone-git-push"]
