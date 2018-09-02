# sbsw-apiserver-exercise

Change to the root/top-most directory of the repo 
Update parameters in automation/setup.sh for your local environment
Variables in this script are:

```console
DOCKER_IMAGE  : desired name of Docker image for API server
DOCKER_REPO : address of your Docker repository
IMAGE_VERSION : a version to tage your Docker image with
```

One the script is updated, run `automation/setup.sh`, this script must be executed from the root/top-most directory of the repo