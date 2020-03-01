node('master') {
    stage('Checkout') {
        cleanWs()
        checkout scm
    }

    stage('Configuration') {
        // add programmacticly needed script approvals
        load('jenkins/config/groovy/scriptApproval.groovy')
    }

    stage('Seed') {
        // https://issues.jenkins-ci.org/browse/JENKINS-44142
        // --> Note: when using multiple Job DSL build steps in a single job, set this to "Delete" only for the last Job DSL build step. 
        // Otherwise views may be deleted and re-created. See JENKINS-44142 for details.
        jobDsl(targets: 'jenkins/jobDSL/*.groovy', sandbox: false, removedJobAction: 'DELETE')
    }
}
