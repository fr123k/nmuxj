node('master') {
    stage('Checkout') {
        cleanWs()
        checkout scm
    }

    stage('Configuration') {
        // set config files in master
        sh('rsync -r ${HOME}/workspace/Admin/Configure/resources/config/configuration-as-code-plugin/ ${HOME}/casc-config/')

        // set userContent
        sh('rsync -r ${HOME}/workspace/Admin/Configure/resources/userContent/ ${HOME}/userContent/')

        // disable csrf for easier jenkins api calls
        load('resources/config/groovy/csrf.groovy')
    }

    stage('Seed') {
        // https://issues.jenkins-ci.org/browse/JENKINS-44142
        // --> Note: when using multiple Job DSL build steps in a single job, set this to "Delete" only for the last Job DSL build step. 
        // Otherwise views may be deleted and re-created. See JENKINS-44142 for details.
        jobDsl(targets: 'resources/jobDSL/sre_dirs.groovy', sandbox: false, removedJobAction: 'IGNORE')
        jobDsl(targets: 'resources/jobDSL/*.groovy', sandbox: false, removedJobAction: 'DELETE')
    }
}
