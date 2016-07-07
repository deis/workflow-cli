def workpath_linux = "/src/github.com/deis/workflow-cli"
def keyfile = "tmp/key.json"

def upload_artifacts = { String filepath ->
	withCredentials([[$class: 'FileBinding', credentialsId: 'e80fd033-dd76-4d96-be79-6c272726fb82', variable: 'GCSKEY']]) {
		sh "cat \"\${GCSKEY}\" > ${filepath}"
		sh "make upload-gcs"
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
			sh 'make bootstrap'
		stage 'Test Linux'
			sh 'make test'
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

					// HACK: Recommended approach for getting command output is writing to and then reading a file.
					sh 'mkdir -p tmp'
					sh 'git describe --all --exact-match > tmp/GIT_BRANCH'
					sh 'git tag -l --contains HEAD > tmp/GIT_TAG'
					def git_branch = readFile('tmp/GIT_BRANCH')
					def git_tag = readFile('tmp/GIT_TAG')

					if (git_branch != "heads/master" && git_tag == "") {
						echo "Skipping build of 386 binaries to shorten CI for Pull Requests"
						env.BUILD_ARCH = "amd64"
					}

					sh 'make bootstrap'
					sh 'make build-revision'

					upload_artifacts(keyfile)
			}
		}
	},
	latest: {
		node('linux') {
			def gopath = gopath_linux()
			def workdir = workdir_linux(gopath)

			dir(workdir) {
					checkout scm

					// HACK: Recommended approach for getting command output is writing to and then reading a file.
					sh 'mkdir -p tmp'
					sh 'git describe --all --exact-match > tmp/GIT_BRANCH'
					def git_branch = readFile('tmp/GIT_BRANCH')

					if (git_branch == "heads/master") {
						sh 'make bootstrap'
						sh 'make build-latest'

						upload_artifacts(keyfile)
					} else {
						echo "Skipping build of latest artifacts because this build is not on the master branch (branch: ${git_branch})"
					}
			}
		}
	}
)
