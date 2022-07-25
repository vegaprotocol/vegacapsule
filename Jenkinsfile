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
        string( name: 'VEGA_CORE_BRANCH', defaultValue: '',
                description: '''Git branch, tag or hash of the vegaprotocol/vega repository.
                    e.g. "develop", "v0.44.0" or commit hash. Default empty: use latests published version.''')
        string( name: 'DATA_NODE_BRANCH', defaultValue: '',
                description: '''Git branch, tag or hash of the vegaprotocol/data-node repository.
                    e.g. "develop", "v0.44.0" or commit hash. Default empty: use latests published version.''')
        string( name: 'PROTOS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/protos repository')
        string( name: 'VEGATOOLS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/vegatools repository')
        string( name: 'SYSTEM_TESTS_BRANCH', defaultValue: 'develop',
                description: 'Git branch, tag or hash of the vegaprotocol/system-tests repository')
        string( name: 'DEVOPS_INFRA_BRANCH', defaultValue: 'master',
                description: 'Git branch, tag or hash of the vegaprotocol/devops-infra repository')
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
                                vegacapsule: commitHash,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                protos: params.PROTOS_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH,
                                testMark: "network_infra_smoke"
                        }
                    }
                }
                stage('System Tests') {
                    steps {
                        script {
                            systemTestsCapsule ignoreFailure: !isPRBuild(),
                                vegacapsule: commitHash,
                                systemTests: params.SYSTEM_TESTS_BRANCH,
                                vegaCore: params.VEGA_CORE_BRANCH,
                                dataNode: params.DATA_NODE_BRANCH,
                                protos: params.PROTOS_BRANCH,
                                vegatools: params.VEGATOOLS_BRANCH,
                                devopsInfra: params.DEVOPS_INFRA_BRANCH
                        }
                    }
                }
            }
        }
    }
}
