FROM debian:buster

ADD setup.sh .
RUN ./setup.sh

ADD server.conf vernemq.conf.template generate_config.sh ./
RUN ./generate_config.sh -c server.conf -o /usr/bin/vernemq/etc/vernemq.conf -u /usr/bin/vernemq/etc/vmq.passwd

ADD vmq.acl /usr/bin/vernemq/etc/vmq.acl

ENV VERNEMQ_PATH /usr/bin/vernemq
ADD run_vernemq.sh .
CMD [ "./run_vernemq.sh" ]
# CMD ["bash"]
