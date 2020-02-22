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
