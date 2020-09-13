this 是在整个class的父级

this.所有对象都可以传递到library里

```groovy
class DockerBuildx {
    DockerBuildx {
        this.xx = steps
    }
   private final script

  DockerBuildx(final script) {
    this.xx = script
  }
}
```





No such DSL method '$' found among steps [acceptGitLabMR, addGitLabMRComment, archive, bat, build, catchError, checkout, compareVersions, container, containerLog, deleteDir, dir, dockerFingerprintFrom, dockerFingerprintRun, echo, emailext, emailextrecipients, envVarsForTool, error, fileExists, findBuildScans, findFiles, getContext, git, gitlabBuilds, gitlabCommitStatus, input, isUnix, jiraComment, jiraIssueSelector, jiraSearch, junit, library, libraryResource, load, lock, mail, milestone, node, nodesByLabel, parallel, podTemplate, powershell, properties, publishHTML, pwd, pwsh, readCSV, readFile, readJSON, readManifest, readMavenPom, readProperties, readTrusted, readYaml, resolveScm, retry, script, sh, sha1, sleep, sshCommand, sshGet, sshPut, sshRemove, sshScript, stage, stash, step, svn, tee, timeout, timestamps, tm, tool, touch, unarchive, unstable, unstash, unzip, updateGitlabCommitStatus, validateDeclarativePipeline, waitUntil, warnError, withContext, withCredentials, withDockerContainer, withDockerRegistry, withDockerServer, withEnv, withGradle, wrap, writeCSV, writeFile, writeJSON, writeMavenPom, writeYaml, ws, zip] or symbols [all, allOf, always, ant, antFromApache, antOutcome, antTarget, any, anyOf, apiToken, architecture, archiveArtifacts, artifactManager, authorizationMatrix, batchFile, bitbucket, bitbucketBranchDiscovery, bitbucketForkDiscovery, bitbucketPublicRepoPullRequestFilter, bitbucketPullRequestDiscovery, bitbucketSshCheckout, bitbucketTagDiscovery, bitbucketTrustEveryone, bitbucketTrustNobody, bitbucketTrustProject, bitbucketTrustTeam, bitbucketWebhookConfiguration, bitbucketWebhookRegistration, booleanParam, branch, brokenBuildSuspects, brokenTestsSuspects, buildButton, buildDiscarder, buildDiscarders, buildParameter, buildSelector, buildTimestamp, buildTimestampExtraProperties, buildingTag, caseInsensitive, caseSensitive, certificate, changeRequest, changelog, changeset, checkoutToSubdirectory, choice, choiceParam, cleanWs, clock, command, configMapVolume, containerEnvVar, containerLivenessProbe, containerTemplate, copyArtifactPermission, copyArtifacts, credentials, cron, crumb, culprits, default, defaultFolderConfiguration, defaultView, demand, developers, disableConcurrentBuilds, disableResume, docker, dockerCert, dockerfile, downstream, dumb, durabilityHint, dynamicPVC, emptyDirVolume, emptyDirWorkspaceVolume, envInject, envVar, envVars, environment, equals, expression, file, fileParam, filePath, fingerprint, frameOptions, freeStyle, freeStyleJob, fromScm, fromSource, git, gitBranchDiscovery, gitHubBranchDiscovery, gitHubBranchHeadAuthority, gitHubExcludeArchivedRepositories, gitHubForkDiscovery, gitHubPullRequestDiscovery, gitHubSshCheckout, gitHubTagDiscovery, gitHubTrustContributors, gitHubTrustEveryone, gitHubTrustNobody, gitHubTrustPermissions, gitLabConnection, gitTagDiscovery, github, githubPush, gitlab, gradle, headRegexFilter, headWildcardFilter, hostPathVolume, hostPathWorkspaceVolume, hyperlink, hyperlinkToModels, inheriting, inheritingGlobal, installSource, isRestartedRun, jdk, jdkInstaller, jgit, jgitapache, jnlp, jobBuildDiscarder, jobName, kubernetes, label, lastCompleted, lastDuration, lastFailure, lastGrantedAuthorities, lastStable, lastSuccess, lastSuccessful, lastWithArtifacts, latestSavedBuild, legacy, legacySCM, list, local, location, logRotator, loggedInUsersCanDoAnything, masterBuild, maven, maven3Mojos, mavenErrors, mavenMojos, mavenWarnings, merge, modernSCM, myView, never, newContainerPerStage, nfsVolume, nfsWorkspaceVolume, node, nodeProperties, nonInheriting, none, not, onFailure, override, overrideIndexTriggers, paneStatus, parallelsAlwaysFailFast, parameters, password, pattern, permalink, permanent, persistentVolumeClaim, persistentVolumeClaimWorkspaceVolume, pipeline-model, pipeline-model-docker, pipelineTriggers, plainText, plugin, podAnnotation, podEnvVar, podLabel, pollSCM, portMapping, preserveStashes, projectNamingStrategy, proxy, queueItemAuthenticator, quietPeriod, rateLimitBuilds, recipients, requestor, resourceRoot, run, runParam, sSHLauncher, schedule, scmRetryCount, scriptApproval, scriptApprovalLink, search, secretEnvVar, secretVolume, security, shell, simpleBuildDiscarder, skipDefaultCheckout, skipStagesAfterUnstable, slave, sourceRegexFilter, sourceWildcardFilter, specific, ssh, sshPublisher, sshPublisherDesc, sshTransfer, sshUserPrivateKey, stackTrace, standard, status, string, stringParam, swapSpace, tag, teamSlugFilter, text, textParam, timezone, tmpSpace, toolLocation, triggeredBy, unsecured, upstream, upstreamDevelopers, userSeed, usernameColonPassword, usernamePassword, viewsTabBar, weather, withAnt, workspace, zfs, zip] or globals [currentBuild, docker, env, params, pipeline, scm, ssh]



[2020-08-30T00:12:11.966Z] WorkflowScript: 177: Only one of "matrix", "parallel", "stages", or "steps" allowed for stage "test" @ line 177, column 7.

WorkflowScript: 269: Unknown stage section "stage". Starting with version 0.5, steps in a stage must be in a ‘steps’ block. @ line 269, column 7.


withCredentials: Masking supported pattern matches of $identity or $ or $userName 覆盖$identity或$或$userName支持的模式匹配

```groovy
checkout([$class: 'GitSCM', branches: [[name: '*/master'], [name: '*/heng/jenkins']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: 'fca8630e-c882-4adc-aecb-b56f631cb6b6', url: 'http://gitlab.esgyn.cn/esgyn/om.git']]])

git credentialsId: 'fca8630e-c882-4adc-aecb-b56f631cb6b6', url: 'http://gitlab.esgyn.cn/heng.sun/shared-library.git'


checkout([$class: 'GitSCM',
  branches: [[name: '*/branch_name']],
  doGenerateSubmoduleConfigurations: false,
  extensions: [[$class: 'RelativeTargetDirectory',
    relativeTargetDir: 'different_directory']],
  submoduleCfg: [],
  userRemoteConfigs: [[url: 'git@github.domain:org/repo.git']]])
```

Deleting pods in bad state
kubectl get pods -o name --selector=jenkins=slave --all-namespaces  | xargs -I {} kubectl delete {}
kubectl get pods -o name --selector=jenkins=slave --all-namespaces  | xargs -I {} kubectl delete {}


docker buildx build -t  reg.esgyn.cn/om-test/docker-rc:v1 -t  reg.esgyn.cn/om-test/docker-rc:latest --platform  linux/arm64,linux/amd64 --push .


docker buildx build -t reg.esgyn.cn/om-test/mgmt-exporter:7cadaca8 -t reg.esgyn.cn/om-test/mgmt-exporter:latest --cache-to=type=inline --platform  linux/arm64,linux/amd64 -f ./Dockerfile.multi --push .

 when {
              expression {
                currentBuild.result == null || currentBuild.result == 'SUCCESS' 
              }
            }