def windows = 'windows'
def linux = 'linux'
def git_commit = ''
def git_branch = ''

def getBasePath = { String filepath ->
	def filename = filepath.lastIndexOf(File.separator)
	return filepath.substring(0, filename)
}

def sh = { String cmd ->
	wrap([$class: 'AnsiColorBuildWrapper', 'colorMapName': 'XTerm']) {
		sh cmd
	}
}

def pscmd = { String cmd ->
	"powershell -NoProfile -ExecutionPolicy Bypass -Command \"${cmd}\""
}

def bootstrap = { String node ->
	bootstrapCmd = node == windows ? bat(pscmd('.\\make bootstrap')) : make('bootstrap')

	try {
		bootstrapCmd
	} catch(error) {
		echo "bootstrap failed; wiping 'vendor' directory and trying again..."
		retry(1) {
			dir('vendor') { deleteDir() }
			bootstrapCmd
		}
	}
}

node(windows) {
	def gopath = pwd() + "\\gopath"
	env.GOPATH = gopath
	def workdir = gopath + "\\src\\github.com\\deis\\workflow-cli"

	dir(workdir) {
		stage 'Checkout Windows'
			checkout scm
		stage 'Install Windows'
			bootstrap windows
		stage 'Test Windows'
			bat pscmd('.\\make test')
	}
}

stage 'Git Info'
node(linux) {
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

def test_image = "quay.io/deisci/workflow-cli-dev:${git_commit.take(7)}"
def mutable_image = 'quay.io/deisci/workflow-cli-dev:latest'

node(linux) {
		stage 'Build and push test container'
			checkout scm
			def quayUsername = "deisci+jenkins"
			def quayEmail = "deis+jenkins@deis.com"
			withCredentials([[$class: 'StringBinding',
												credentialsId: 'c67dc0a1-c8c4-4568-a73d-53ad8530ceeb',
									 			variable: 'QUAY_PASSWORD']]) {

				sh """
					docker login -e="${quayEmail}" -u="${quayUsername}" -p="\${QUAY_PASSWORD}" quay.io
					docker build -t ${test_image} .
					docker push ${test_image}
				"""

				if (git_branch == "remotes/origin/master") {
					sh """
						docker tag ${test_image} ${mutable_image}
						docker push ${mutable_image}
					"""
				}
			}
}


stage 'Lint and test container'
parallel(
	lint: {
		node(linux) {
			sh "docker run --rm ${test_image} lint"
		}
	},
	test: {
		node(linux) {
			withCredentials([[$class: 'StringBinding',
												credentialsId: '995d99a7-466b-4beb-bf75-f3ba91cbbc18',
												variable: 'CODECOV_TOKEN']]) {
				sh "docker run -e CODECOV_TOKEN=\${CODECOV_TOKEN} --rm ${test_image} sh -c 'test-cover.sh && curl -s https://codecov.io/bash | bash'"
			}
		}
	}
)

stage 'Build and Upload Artifacts'

def gcs_bucket = "gs://workflow-cli"

def upload_artifacts = { String dist_dir, boolean cache ->
	headers  = "-h 'x-goog-meta-git-branch:${git_branch}' "
	headers += "-h 'x-goog-meta-git-sha:${git_commit}' "
	headers += "-h 'x-goog-meta-ci-job:${env.JOB_NAME}' "
	headers += "-h 'x-goog-meta-ci-number:${env.BUILD_NUMBER}' "
	headers += "-h 'x-goog-meta-ci-url:${env.BUILD_URL}'"
	if(!cache) {
		headers += ' -h "Cache-Control:no-cache"'
	}

	script = "sh -c 'echo \${GCS_KEY_JSON} | base64 -d - > /tmp/key.json "
	script += "&& gcloud auth activate-service-account -q --key-file /tmp/key.json "
	script += "&& gsutil -mq ${headers} cp -a public-read -r /upload/* ${gcs_bucket} "
	script += "&& rm -rf /upload/*'"

	withCredentials([[$class: 'StringBinding',
										credentialsId: '6561701c-b7b4-4796-83c4-9d87946799e4',
										variable: 'GCSKEY']]) {
		sh "docker run ${dist_dir} -e GCS_KEY_JSON=\"\${GCSKEY}\" --rm ${test_image} ${script}"
	}
	sh "rm -rf ${tmp_dir}"
}

def mktmp = {
	// Create tmp directory to store files
	sh 'mktemp -d > tmp_dir'
	tmp = readFile('tmp_dir').trim()
	echo "Storing binaries in ${tmp}"
	sh 'rm tmp_dir'
	return tmp
}

parallel(
	revision: {
		node(linux) {

			flags = ""
			if (git_branch != "remotes/origin/master") {
				echo "Skipping build of 386 binaries to shorten CI for Pull Requests"
				flags += "-e BUILD_ARCH=amd64"
			}

			tmp_dir = mktmp()
			dist_dir = "-e DIST_DIR=/upload -v ${tmp_dir}:/upload"
			sh "docker run ${flags} -e REVISION=${git_commit.take(7)} ${dist_dir} --rm ${test_image} make build-revision"


			upload_artifacts(dist_dir, true)
		}
	},
	latest: {
		node(linux) {
			if (git_branch == "remotes/origin/master") {
				tmp_dir = mktmp()
				dist_dir = "-e DIST_DIR=/upload -v ${tmp_dir}:/upload"
				sh "docker run ${dist_dir} --rm ${test_image} make build-latest"

				upload_artifacts(dist_dir, false)
			} else {
				echo "Skipping build of latest artifacts because this build is not on the master branch (branch: ${git_branch})"
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

		node(linux) {
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
