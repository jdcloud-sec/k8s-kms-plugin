FROM alpine

COPY build/ /k8s-kms-plugin

VOLUME /etc/k8s-kms-plugin.json

CMD [ "/k8s-kms-plugin" ]
