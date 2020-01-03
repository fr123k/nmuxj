#!groovy

folder('Admin') {
    description('Folder containing configuration and seed jobs')
}

pipelineJob("Admin/Configure") {
    parameters {
        gitParam('revision') {
            type('BRANCH_TAG')
            sortMode('ASCENDING_SMART')
            defaultValue('origin/master')
        }
    }

    triggers {
        githubPush()
    }

    logRotator {
        numToKeep(50)
    }

    definition {
        cpsScm {
            scm {
                git {
                    remote {
                        github("", "ssh")
                        credentials("ssh-git")
                    }

                    branch('$revision')
                }
            }
            
            scriptPath('pipeline.groovy')
        }
    }
}
