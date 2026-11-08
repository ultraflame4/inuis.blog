pipeline {
    agent any
    tools { go '1.26' }

    stages {
        stage('Fetch') {
            steps {
                // Get some code from a GitHub repository
                git branch: 'main', url: 'https://github.com/ultraflame4/inuis.blog.git'

            }
        }
       stage('Build Client') {
            steps {
                dir('web') {
                    nodejs('NodeJS') {
                        sh 'npm ci'
                        sh 'npm run build'
                        sh 'tar -czvf web-client.tar.gz ./build'
                        stash includes: 'web-client.tar.gz', name: 'web-client'
                        archiveArtifacts artifacts: 'web-client.tar.gz', followSymlinks: false
                    }
  
                }
            }
        }
       stage('Prepare Embed') {
            steps {
                unstash 'web-client'
                sh 'tar -xf web-client.tar.gz'
                sh '''
                    rm -rf ./embed/dist
                    mv ./build ./embed/dist
                '''
            }
        }
        stage('Build Go') {
            steps {
                sh 'mkdir -p ./bin'
                sh 'go build -o ./bin/app .'
                archiveArtifacts artifacts: 'bin/app', followSymlinks: false
            }
        }
    }
}
