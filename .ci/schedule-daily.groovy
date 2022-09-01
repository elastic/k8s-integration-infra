@Library('apm@current') _

pipeline {
  agent none
  environment {
    NOTIFY_TO = credentials('notify-to')
    PIPELINE_LOG_LEVEL = 'INFO'
  }
  options {
    timeout(time: 1, unit: 'HOURS')
    buildDiscarder(logRotator(numToKeepStr: '20', artifactNumToKeepStr: '20'))
    timestamps()
    ansiColor('xterm')
    disableResume()
    durabilityHint('PERFORMANCE_OPTIMIZED')
  }
  triggers {
    cron('H H(2-3) * * 1-5')
  }
  stages {
    stage('Nighly beats builds') {
      steps {
        runBuilds(quietPeriodFactor: 2000, branches: ['main'])
      }
    }
  }
  post {
    cleanup {
      notifyBuildResult(prComment: false)
    }
  }
}

def runBuilds(Map args = [:]) {
  def branches = getBranchesFromAliases(aliases: args.branches)
  def quietPeriod = 0
  branches.each { branch ->
    build(quietPeriod: quietPeriod,
          job: "cloudnative/k8s-integration-infra-mbp/${branch}",
          parameters: [
            booleanParam(name: 'rally_tests_ci', value: true),
          ],
          wait: false, propagate: false)
    quietPeriod += args.quietPeriodFactor
  }
}