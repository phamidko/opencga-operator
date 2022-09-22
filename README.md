# OpenCGA Kubernetes Operator (opencga-operator) Architecture

OpenCGA Kubernetes Operator is a Custom Resource Definition (OpenCGACommunity, CRD) and a Controller

### Table of Contexts
- [Cluster Configuration](#cluster-configuration)


## Cluster Configuration
By deploying OpenCGACommunity resource definition, the Operator:
1. Creates a StatefulSet that contains one pod for each REST member.
2. Writes the Automation configuration as a Secret and mounts it to each pod.

3. Creates one init container and two containers in each pod:
    -  An init container which copies the `cmd/versionhook` binary to the main `opencga-REST` container. This is run before `opencga-REST` starts to handle version upgrades
    -  A container of `opencga-REST` is jetty server webapp, It handles data queries.
    -  A container of `opencga-agent`. The Automation function of the OpenCGA Agent handles configuring, stopping, and restarting the `opencga-REST` process. The OpenCGA Agent periodically polls `opencga-REST` to determine status and can deploy changes as needed

3. Creates one init container and two containers in each pod for the purpose of `opencga-client` (MASTER)
    -  An init container which copies the `cmd/versionhook` binary to the main `opencga-client` container. This is run before `opencga-client` starts to handle version upgrades
    -  A container of `opencga-client` is a shell command-line access, It handles data queries and operations such as migration and indexing.
    -  A container of `opencga-agent`. The Automation function of the OpenCGA Agent handles configuring, stopping, and restarting the `opencga-client` process. The OpenCGA Agent periodically polls `opencga-client` to determine status and can deploy changes as needed

4. Creates several volumes:
    -  `data-volume` which are persistent and mount such as folder suchas (sessions, variants, log) to /data on both the server and agent containers. Stores server data as well as automation-opencga.conf written by the agent and some locks the agent needs.
    -  `automation-config` which is mounted from the previously generated Secret to both the server and agent. Only lives as long as the pod.
    -  `healthstatus` which contains the agent's current status. This is shared with the `opencga-REST` and `opencga-client` container where it's used by the pre-stop hook. Only lives as long as the pod.

5. Initiates `opencga-agent`, which in turn creates the database configuration and launches the `opencga-REST` and `opencga-client` process according to the OpenCGACommunity resource definition

## HOW-TO Steps

Install the CRD
Install the necessary roles and role-bindings
Install the Operator


By default, the operator will creates three pods, each of them automatically linked to a new persistent volume claim bounded to a new persistent volume also created by the operator

