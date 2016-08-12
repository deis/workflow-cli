def workpath_linux = "/src/github.com/deis/workflow-cli"

def getBasePath = { String filepath ->
	def filename = filepath.lastIndexOf(File.separator)
	return filepath.substring(0, filename)
}

def make = { String target ->
	try {
		sh "make ${target} fileperms"
	} catch(error) {
		sh "make fileperms"
		throw error
	}
}

def gcs_cleanup_cmd = "sh -c 'rm -rf /.config/*'"
def gcs_bucket = "gs://workflow-cli"
def gcs_key = "tmp/key.json"

def gcs_cmd = { String cmd ->
	gcs_cmd = "docker run --rm -v  ${pwd()}/tmp:/.config -v ${pwd()}/_dist:/upload google/cloud-sdk:latest "
	try {
		sh(gcs_cmd + cmd)
	} catch(error) {
		sh(gcs_cmd + gcs_cleanup_cmd)
		throw error
	}
}

def upload_artifacts = { boolean cache ->
	withCredentials([[$class: 'FileBinding', credentialsId: 'e80fd033-dd76-4d96-be79-6c272726fb82', variable: 'GCSKEY']]) {
		sh "mkdir -p ${getBasePath(gcs_key)}"
		sh "cat \"\${GCSKEY}\" > ${gcs_key}"
		gcs_cmd 'gcloud auth activate-service-account -q --key-file /.config/key.json'
		headers = ""
		if(!cache) {
			headers += '-h "Cache-Control:no-cache"'
		}
		gcs_cmd "gsutil -mq ${headers} cp -a public-read -r /upload/* ${gcs_bucket}"
		gcs_cmd gcs_cleanup_cmd
	}
}

def gopath_linux = {
	def gopath = pwd() + "/gopath"
	env.GOPATH = gopath
	gopath
}

def workdir_linux = { String gopath ->
	gopath + workpath_linux
}

def sh = { String cmd ->
	wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm']) {
		sh cmd
	}
}

node('windows') {
	def gopath = pwd() + "\\gopath"
	env.GOPATH = gopath
	def workdir = gopath + "\\src\\github.com\\deis\\workflow-cli"

	def pscmd = { String cmd ->
		"powershell -NoProfile -ExecutionPolicy Bypass -Command \"${cmd}\""
	}

	dir(workdir) {
		stage 'Checkout Windows'
			checkout scm
		stage 'Install Windows'
			bat pscmd('.\\make bootstrap')
		stage 'Test Windows'
			bat pscmd('.\\make test')
	}
}

node('linux') {
	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath)

	dir(workdir) {
		stage 'Checkout Linux'
			checkout scm
		stage 'Install Linux'
			make 'bootstrap'
		stage 'Test Linux'
			make 'test-cover'
		stage 'Upload to Codeconv'
			withCredentials([[$class: 'StringBinding', credentialsId: '995d99a7-466b-4beb-bf75-f3ba91cbbc18', variable: 'CODECOV_TOKEN']]) {
				sh 'curl -s https://codecov.io/bash | bash'
			}
	}
}

def git_commit = ''
def git_branch = ''

stage 'Git Info'
node('linux') {

	def gopath = gopath_linux()
	def workdir = workdir_linux(gopath)

	dir(workdir) {
		checkout scm

		// HACK: Recommended approach for getting command output is writing to and then reading a file.
		sh 'mkdir -p tmp'
		sh 'git describe --all > tmp/GIT_BRANCH'
		sh 'git rev-parse HEAD > tmp/GIT_COMMIT'
		git_branch = readFile('tmp/GIT_BRANCH').trim()
		git_commit = readFile('tmp/GIT_COMMIT').trim()

		if (git_branch != "remotes/origin/master") {
			// Determine actual PR commit, if necessary
			sh 'git rev-parse HEAD | git log --pretty=%P -n 1 --date-order > tmp/MERGE_COMMIT_PARENTS'
			sh 'cat tmp/MERGE_COMMIT_PARENTS'
			merge_commit_parents = readFile('tmp/MERGE_COMMIT_PARENTS').trim()
			if (merge_commit_parents.length() > 40) {
				echo 'More than one merge commit parent signifies that the merge commit is not the PR commit'
				echo "Changing git_commit from '${git_commit}' to '${merge_commit_parents.take(40)}'"
				git_commit = merge_commit_parents.take(40)
			} else {
				echo 'Only one merge commit parent signifies that the merge commit is also the PR commit'
				echo "Keeping git_commit as '${git_commit}'"
			}
		}
	}
}

stage 'Build and Upload Artifacts'

parallel(
	revision: {
		node('linux') {
			def gopath = gopath_linux()
			def workdir = workdir_linux(gopath)

			dir(workdir) {
					checkout scm

					if (git_branch != "remotes/origin/master") {
						echo "Skipping build of 386 binaries to shorten CI for Pull Requests"
						env.BUILD_ARCH = "amd64"
					}

					make 'bootstrap'
					env.REVISION = git_commit.take(7)
					make 'build-revision'

					upload_artifacts(true)
			}
		}
	},
	latest: {
		node('linux') {
			def gopath = gopath_linux()
			def workdir = workdir_linux(gopath)

			dir(workdir) {
					checkout scm

					if (git_branch == "remotes/origin/master") {
						make 'bootstrap'
						make 'build-latest'

						upload_artifacts(false)
					} else {
						echo "Skipping build of latest artifacts because this build is not on the master branch (branch: ${git_branch})"
					}
			}
		}
	}
)

stage 'Trigger e2e tests'

// If build is on master, trigger workflow-test, otherwise, assume build is a PR and trigger workflow-test-pr
waitUntil {
	try {
		if (git_branch == "remotes/origin/master") {
			build job: '/workflow-test', parameters: [[$class: 'StringParameterValue', name: 'WORKFLOW_CLI_SHA', value: git_commit],
				[$class: 'StringParameterValue', name: 'COMPONENT_REPO', value: 'workflow-cli']]
		} else {
			build job: '/workflow-test-pr', parameters: [[$class: 'StringParameterValue', name: 'WORKFLOW_CLI_SHA', value: git_commit],
				[$class: 'StringParameterValue', name: 'COMPONENT_REPO', value: 'workflow-cli']]
		}
		true
	} catch(error) {
		if (git_branch == "remotes/origin/master") {
			throw error
		}

		node('linux') {
			withCredentials([[$class: 'StringBinding', credentialsId: '8a727911-596f-4057-97c2-b9e23de5268d', variable: 'SLACKEMAIL']]) {
				mail body: """<!DOCTYPE html>
<html>
<head>
<meta content='text/html; charset=UTF-8' http-equiv='Content-Type' />
</head>
<body>
<div>Author: ${env.CHANGE_AUTHOR}<br/>
Branch: ${env.BRANCH_NAME}<br/>
Commit: ${env.CHANGE_TITLE}<br/>
<a href="${env.BUILD_URL}console">Click here</a> to view logs.</p>
<a href="${env.BUILD_URL}input/">Click here</a> to restart e2e.</p>
</div>
</html>
""", from: 'jenkins@ci.deis.io', subject: 'Workflow CLI E2E Test Failure', to: env.SLACKEMAIL, mimeType: 'text/html'
			}
			input "Retry the e2e tests?"
		}
		false
	}
}
