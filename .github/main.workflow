workflow "Test" {
  on = "push"
  resolves = ["CodeCov"]
}

action "Setup Go" {
  uses = "actions/setup-go@632d18fc920ce2926be9c976a5465e1854adc7bd"
}

action "CodeCov" {
  uses = "actions/docker/cli@fe7ed3ce992160973df86480b83a2f8ed581cd50"
  needs = ["Setup Go"]
  args = "codecov"
  secrets = ["CODECOV_TOKEN"]
}
