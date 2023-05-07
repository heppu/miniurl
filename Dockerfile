FROM scratch
COPY target/miniurl /
ENTRYPOINT [ "/miniurl" ]
