FROM mcr.microsoft.com/devcontainers/go

# Install go dependencies
Run go mod download
# Install goyang YANG parser
RUN go install github.com/openconfig/goyang@latest
# Install ygot code generator
RUN go install github.com/openconfig/ygot/generator@latest