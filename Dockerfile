FROM scratch
COPY agbridge /usr/bin/agbridge
ENTRYPOINT [ "/usr/bin/agbridge" ]