FROM busybox

EXPOSE 3306

COPY ./majiang /
COPY ./conf.yaml /etc/


ENTRYPOINT ["/majiang"]