FROM node:22.11.0

WORKDIR /root

RUN mkdir -p /output

RUN npm install aws-ses-v2-local -g

RUN cat <<'EOF' > /bin/aws-ses-v2-local.sh
#!/bin/bash
#
# Wrapper for SES Local emulator which monitors kill signals. Without this,
# the `aws-ses-v2-local` CLI does not respond to Ctrl-C making shutdown times
# longer when running in Docker.

( exec aws-ses-v2-local $@ ) &
pid=$!

trap "kill -9 $pid" SIGINT SIGTERM

wait $pid
EOF
RUN chmod +x /bin/aws-ses-v2-local.sh

EXPOSE 8005

ENTRYPOINT ["/bin/aws-ses-v2-local.sh", "--port", "8005", "--host", "0.0.0.0"]