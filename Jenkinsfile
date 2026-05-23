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
                script {
                    def tag = sh(returnStdout: true, script: 'git tag --points-at HEAD').trim()
                    sh "GORELEASER_CURRENT_TAG=${tag} goreleaser release --clean"
                }
            }
        }
    }

    post {
        always {
            cleanWs()
        }
        success {
            sh """
                curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
                    -d chat_id="${TELEGRAM_CHAT_ID}" \
                    -d parse_mode="Markdown" \
                    -d text="✅ *expense-tracker deployed*%0ARepo: ${env.JOB_NAME}%0ABranch: ${env.GIT_BRANCH}%0ACommit: ${env.GIT_COMMIT.take(7)}%0AView run: ${env.BUILD_URL}"
            """
        }
        failure {
            sh """
                curl -s -X POST "https://api.telegram.org/bot${TELEGRAM_BOT_TOKEN}/sendMessage" \
                    -d chat_id="${TELEGRAM_CHAT_ID}" \
                    -d parse_mode="Markdown" \
                    -d text="❌ *expense-tracker deploy failed*%0ARepo: ${env.JOB_NAME}%0ABranch: ${env.GIT_BRANCH}%0ACommit: ${env.GIT_COMMIT.take(7)}%0AView run: ${env.BUILD_URL}"
            """
        }
    }
}
