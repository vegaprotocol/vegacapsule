/* groovylint-disable DuplicateStringLiteral, LineLength, NestedBlockDepth */
@Library('vega-shared-library') _

def commitHash = 'UNKNOWN'

pipeline {
    agent any
    options {
        skipDefaultCheckout true
        timestamps()
        timeout(time: 45, unit: 'MINUTES')
    }
    parameters {
        string( name: 'VEGA_BRANCH', defaultValue: '',
                description: '''Git branch, tag or hash of the vegaprotocol/vega repository.
                    e.g. "develop", "v0.44.0" or commit hash. Default empty: use latests published version.''')
        string( name: 'SYSTEM_TESTS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/system-tests repository')
        string( name: 'VEGATOOLS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/vegatools repository')
        string( name: 'DEVOPS_INFRA_BRANCH', defaultValue: 'master',
                description: 'Git branch, tag or hash of the vegaprotocol/devops-infra repository')
        string( name: 'DEVOPSSCRIPTS_BRANCH', defaultValue: 'main',
                description: 'Git branch, tag or hash of the vegaprotocol/devopsscripts repository')
        string(name: 'JENKINS_SHARED_LIB_BRANCH', defaultValue: 'main',
                description: 'Git branch name of the vegaprotocol/jenkins-shared-library repository')
    }

    stages {
        stage('Config') {
            steps {
                cleanWs()
                sh 'printenv'
                echo "params=${params}"
                echo "isPRBuild=${isPRBuild()}"
                script {
                    params = pr.injectPRParams()
                }
                echo "params (after injection)=${params}"
            }
        }

        stage('Git clone') {
            options { retry(3) }
            steps {
                checkout scm
                script {
                    commitHash = getCommitHash()
                }
                echo "commitHash=${commitHash}"
            }
        }

        stage('Tests') {
            parallel {
                stage('System Tests Network Smoke') {
                    steps {
                        script {
                            systemTestsCapsule ignoreFailure: !isPRBuild(),
                                timeout: 30,
                                vegacapsule: commitHash,
                                vegaVersion: params.VEGA_BRANCH,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                devopsScripts: params.DEVOPSSCRIPTS_BRANCH,
                                jenkinsSharedLib: params.JENKINS_SHARED_LIB_BRANCH,
                                testMark: "network_infra_smoke"
                        }
                    }
                }
                stage('System Tests') {
                    steps {
                        script {
                            systemTestsCapsule ignoreFailure: !isPRBuild(),
                                timeout: 30,
                                vegacapsule: commitHash,
                                vegaVersion: params.VEGA_BRANCH,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                devopsScripts: params.DEVOPSSCRIPTS_BRANCH,
                                jenkinsSharedLib: params.JENKINS_SHARED_LIB_BRANCH
                        }
                    }
                }
            }
        }
    }
}
