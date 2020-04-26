pipeline {
	agent any
	stages {
		stage ('Checkking out Git files'){
			steps {
				checkout scm
			}
		}
		stage ('Prepare the Environment') {
			steps {
				script {
					def tfHome = tool 'Terraform'
					def jdk = tool 'jdk8'
					env.PATH = "${tfHome}:${env.PATH}"
				}
				sh 'terraform --version'
			}
		}

		stage('Testing') { 
			steps {
				script { 
					def root = tool name: 'Go'
					withEnv(["GOROOT=${root}", "PATH+GO=${root}/bin"]) {
					sh 'go version'
						dir("test"){         // test dir
							sh 'pwd'
							sh 'go test -v'  // put test  here
						}

					}
				}
			}
		}
		stage ('Provisioning Infrastructure') {
			steps {
				dir ('Azure') {
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
