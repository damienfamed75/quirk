workflow "Test" {
  resolves = [
    "CodeCov",
    "Ilshidur/action-discord@master",
  ]
  on = "push"
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

action "Ilshidur/action-discord@master" {
  uses = "Ilshidur/action-discord@master"
  secrets = ["DISCORD_WEBHOOK"]
}
