#!groovy

folder('Jenkins') {
    description('Folder containing configuration and seed jobs')
}

pipelineJob("Jenkins/Configure") {
    logRotator {
        numToKeep(50)
    }
    definition {
        cps {
            script("""
node ("master") {
    sh("echo seed job")
}
            """)
        }
    }
}

pipelineJob("provision-redis") {
    parameters {
        stringParam('SERVER', '', '')
    }
    definition {
        cps {
            script("""
node ("master") {
    sh("echo provision redis with " + params.SERVER)
}
            """)
        }
    }
}

pipelineJob("provision-elasticsearch") {
    definition {
        cps {
            script("""
node ("master") {
    sh("echo provision elasticsearch")
}
            """)
        }
    }
}
