# Setup/Deployment #

Change to the root/top-most directory of the repo 
Update parameters in automation/setup.sh for your local environment
Variables in this script are:

```console
DOCKER_IMAGE  : desired name of Docker image for API server
DOCKER_REPO : address of your Docker repository
IMAGE_VERSION : a version to tag your Docker image with
KUBECONFIG_PATH : path to kubeconfig credentials file for your Kubernetes cluster
```

## Changing external IP or port of API server ##

Change these in `Service` section of Kubernetes manifest file `manifests/apiserver.yml`.  (Parameters `port` and IPs under `externalIPs`)

Oncee the script is updated with required variables / external IP/port details, run `automation/setup.sh`, this script must be executed from the root/top-most directory of the repo

## Testing/Usage of the API service ##

URL format:
http://192.168.99.100:8082/query?Person=Jata&Relation=PaternalUncle

Where `Person` is the name of the person queried, and `Relation` is the type of relationship information required.  `Relation` can be one of the following values:

* `paternaluncle`
* `paternalaunt`
* `maternalaunt`
* `maternaluncle`
* `brotherinlaw` or `brother-in-law`
* `sisterinlaw` or `sister-in-law`
* `cousins`

## Other details ##

* Part 1 - Problem 1 code in `src/`
* Part 2 - Dockerfile in `docker/`
* Part 3 - Kube Configuration file / manifest in `manifests/`
* Part 4 - Deployment script at `automation/setup.sh`
