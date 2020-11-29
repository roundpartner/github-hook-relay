FROM scratch

ARG commit_id=master
LABEL maintainer="tom@roundpartner.co.uk"
LABEL org.label-schema.description="GitHub Hook Relay"
LABEL org.label-schema.name="github-hook-relay"
LABEL org.label-schema.schema-version="1.0"
LABEL org.label-schema.vcs-url="https://github.com/roundpartner/github-hook-relay"
LABEL org.label-schema.vcs-ref="${commit_id}"
LABEL org.label-schema.vendor="RoundPartner"

ARG build_number=unknown
ENV VERSION=${build_number}
ENV PATH=/

WORKDIR /
COPY github-hook-relay github-hook-relay

ENTRYPOINT ["github-hook-relay"]
