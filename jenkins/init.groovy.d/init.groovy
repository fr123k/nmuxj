#!groovy

/*
 * This script is designated for the init.groovy.d 
 * directory to be executed at startup time of the 
 * Jenkins instance. This script requires the jobDSL
 * Plugin. Tested with job-dsl:1.70
 */

import javaposse.jobdsl.dsl.DslScriptLoader
import javaposse.jobdsl.plugin.JenkinsJobManagement
import jenkins.model.Jenkins
import hudson.security.FullControlOnceLoggedInAuthorizationStrategy
import hudson.security.HudsonPrivateSecurityRealm
import hudson.model.*
import jenkins.security.*
import jenkins.security.apitoken.*

//TODO configure the disabling of csrf in the casc configuration
Jenkins.getInstance().setCrumbIssuer(null)
println("Disable CSRF protection")

// Create the configuration pipeline from a jobDSL script
def jobDslScript = new File('/var/jenkins_home/dsl/bootstrap.groovy')
def workspace = new File('.')
def jobManagement = new JenkinsJobManagement(System.out, [:], workspace)
new DslScriptLoader(jobManagement).runScript(jobDslScript.text)
// Schdule the Jenkins/Configure job
Jenkins.instance.getItemByFullName("Jenkins/Configure").scheduleBuild()

println(Jenkins.instance.getSecurityRealm().getClass().getSimpleName())
// Disable Wizards
if(Jenkins.instance.getSecurityRealm().getClass().getSimpleName() == 'None') {
    // def instance = Jenkins.getInstance()
    // def setupUser = "admin"
    // def setupPass = "admin"

    // def hudsonRealm = new HudsonPrivateSecurityRealm(false)
    // instance.setSecurityRealm(hudsonRealm)
    // def user = instance.getSecurityRealm().createAccount(setupUser, setupPass)
    // user.save()

    // def prop = user.getProperty(ApiTokenProperty.class)
    // def result = prop.tokenStore.generateNewToken("token-created-by-init-groovy")
    // user.save()

    // println("###################################################\n## Api-Token: " + result.plainValue + " ##\n###################################################")

    // def strategy = new FullControlOnceLoggedInAuthorizationStrategy()
    // strategy.setAllowAnonymousRead(true)
    // instance.setAuthorizationStrategy(strategy)

    // instance.save()

    println("SetupWizard Disabled")
}
