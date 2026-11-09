pipeline {
    agent any
    tools { go '1.26' }

    
    stages {
        stage('Check') {
            steps {
                script {
                    env.GIT_COMMIT_HASH = sh (script: "git log -n 1 --pretty=format:'%H'", returnStdout: true).trim()
                    env.APP_VERSION = sh (script: "cat ./VERSION", returnStdout: true).trim()
                    env.BUILD_DATE = sh (script: "date", returnStdout: true).trim()
                }
                echo "GIT_COMMIT: ${GIT_COMMIT_HASH}"
                echo "VERSION: ${APP_VERSION}"
                echo "BUILD_DATE: ${BUILD_DATE}"
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
                
                script{
                    env.LD_FLAGS = "-X 'main.version=${env.APP_VERSION}' -X 'main.build=${env.BUILD_DATE}' -X 'main.commit=${env.GIT_COMMIT_HASH}'"
                    env.BIN_NAME = "inuis-app"
                }
            
                sh 'mkdir -p ./bin'
           
                echo "LD_FLAGS: ${LD_FLAGS}" 
                sh 'go build -ldflags=\"${LD_FLAGS}\" -o ./bin/${BIN_NAME} . '
                archiveArtifacts artifacts: "bin/${BIN_NAME}", followSymlinks: false
            }
        }
    }
}
