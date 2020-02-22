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
