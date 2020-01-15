FROM alpine

COPY build/k8s-kms-plugin /k8s-kms-plugin

VOLUME /etc/kubernetes/jdcloud-kms-plugin.json

CMD [ "/k8s-kms-plugin" ]
