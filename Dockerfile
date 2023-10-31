# Use a lightweight base image
FROM alpine:3.14

WORKDIR /root/

# Copy the pre-built binary from your host into the container
COPY .build/refineMetric_exporter .

COPY ./metric_configuration.yaml .

# Ensure the binary is executable
RUN chmod +x ./refineMetric_exporter && apk add gcompat

# When mounting the volume, the socket at /var/run/libvirt/libvirt-sock on the host 
# will be available at the same path in the container. 
# Ensure your app has permissions to access the socket.

EXPOSE 9188

USER root

# Command to run the binary
ENTRYPOINT ["./refineMetric_exporter"]
