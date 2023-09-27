docker run -d --restart=always --privileged=true \
--name etcd -p 2379:2379 \
-p 2380:2380 -v /opt/etcd_data:/bitnami/etcd --env ALLOW_NONE_AUTHENTICATION=yes --env ETCD_ADVERTISE_CLIENT_URLS=http://0.0.0.0:2379 --log-opt max-size=100m bitnami/etcd:3.4.15

docker run -d --name=myredis -p 6379:6379 --restart=always redis:latest

