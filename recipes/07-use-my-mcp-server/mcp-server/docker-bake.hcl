variable "REPO" {
  default = "k33g"
}

variable "TAG" {
  default = "with-agents"
}

group "default" {
  targets = ["mcp-demo"]
}

target "mcp-demo" {
  context = "."
  dockerfile = "Dockerfile"
  args = {}
  platforms = [
    //"linux/amd64",
    "linux/arm64"
  ]
  tags = ["${REPO}/mcp-demo:${TAG}"]
}

# docker buildx bake --push --file docker-bake.hcl
