pipeline {
    agent any

    environment {
        GITHUB_TOKEN = credentials('tuigram-github-token')
        PATH = "/usr/local/go/bin:${env.PATH}"
    }

    stages {
        stage('Checkout') {
            steps {
                checkout scm
                sh 'git fetch --tags'
            }
        }

        stage('Verify') {
            steps {
                sh 'go version'
                sh 'goreleaser --version'
            }
        }

        stage('Vet') {
            steps {
                sh 'go vet ./...'
            }
        }

        stage('Build') {
            steps {
                sh 'go build ./...'
            }
        }

        stage('Release') {
            when {
                expression {
                    def tag = sh(returnStdout: true, script: 'git tag --points-at HEAD').trim()
                    return tag ==~ /v\d+\.\d+\.\d+.*/
                }
            }
            steps {
                sh 'goreleaser release --clean'
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            echo 'Pipeline succeeded.'
        }
        failure {
            echo 'Pipeline failed.'
        }
    }
}
