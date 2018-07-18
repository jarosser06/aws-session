def prepareGoWorkspace(gitHost, projectOwner, project) {
  def goProject = '${pwd()}/${gitHost}/${projectOwner}/${project}'
  sh 'mkdir -p ${goProject}'
  return goProject
}

node {
  stage('Main') {
    docker.image('us.gcr.io/rosser-cloud/go-test').inside {
      stage('Checkout Changes') {
        def projectDir = prepareGoWorkspace('github.com', 'jarosser06', 'aws-session')
        dir(projectDir) {
          sh 'ls -lh'
        }
      }
    }
  }
}
