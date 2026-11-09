pipeline {
    agent any
    tools { go '1.26' }

    
    stages {
        stage('Check') {
            steps {
                script {
                    env.GIT_COMMIT_HASH = sh (script: "git log -n 1 --pretty=format:'%H'", returnStdout: true).trim
                    env.GIT_TAG = sh (script: "git describe --exact-match --tags'", returnStdout: true).trim
                    env.BUILD_DATE = sh (script: "date'", returnStdout: true).trim
                }
                echo "GIT_COMMIT: ${GIT_COMMIT_HASH}"
                echo "GIT_TAG: ${GIT_TAG}"
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
                    env.LD_FLAGS = "-X main.version=${env.GIT_TAG} -X main.build=${env.BUILD_DATE}"
                    env.BIN_NAME = "inuis-${env.GIT_TAG}"
                }
            
                sh 'mkdir -p ./bin'
           
                echo "LD_FLAGS: ${LD_FLAGS}" 
                sh "go build -o ./bin/${BIN_NAME} .   -ldflags \"${LD_FLAGS}\""
                archiveArtifacts artifacts: "bin/${BIN_NAME}", followSymlinks: false
            }
        }
    }
}
