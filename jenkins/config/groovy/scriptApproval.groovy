#! groovy
import static org.jenkinsci.plugins.scriptsecurity.scripts.ScriptApproval.get

// approve the two pipeline jobs redis and elasticsearch
get().approveScript('1b1c1589febf4a8382cedc078904e5fa2a6e2254')
get().approveScript('4db0d5783e08301a1ce66f8ada6c9fc8d0d30552')
get().approveScript('a4fd50f0811d749ddea1971702e5f488e351f33f')
get().approveScript('b1d2101d61b9bf0fe45924a955ed8edc3a983474')
