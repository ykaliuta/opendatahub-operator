FROM ubi8:latest
WORKDIR /
ADD webhook .
USER 1001

ENTRYPOINT ["/webhook"]
