node {
        ws('${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}/src') {
            withEnv(['GOPATH=${JENKINS_HOME}/jobs/${JOB_NAME}/builds/${BUILD_ID}']) {
                env.PATH='${GOPATH}/bin:$PATH'

                stage ('Checkking out Git files'){
                        echo 'Checking out SCM'
                        checkout scm
                }
                stage ('Prepare the Environment') {
                                script {
                                        def tfHome = tool 'Terraform'
                                        def jdk = tool 'jdk'
                                        env.PATH = "${tfHome}:${env.PATH}"
                                }
                                sh 'terraform --version'
                                sh 'java -version'
                }

                stage('Testing') {
                                script {
                                        def root = tool 'Go'
                                        sh 'go version'
                                                dir("terratest-tutorial/test/"){         // test dir
                                                        sh 'pwd'
                                                        sh 'sudo /root/go/bin/dep ensure'
                                                        sh 'go test -v'  // put test  here
                                                }

                                        }
                        }
                stage ('Provisioning Infrastructure') {
                                dir ('terratest-tutorial/tf') {
                                        withCredentials([azureServicePrincipal('azurelogin')]){
                                                sh 'terraform init'
                                                sh 'terraform plan -out "plan.out"'
                                                sh 'terraform apply "plan.out"'
                                        }
                                }
                }

        }
    }
}

