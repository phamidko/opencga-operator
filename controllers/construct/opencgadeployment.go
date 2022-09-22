package construct

import (

	// "github.com/phamidko/opencga-operator/pkg/automationconfig"
	// "github.com/phamidko/opencga-operator/pkg/kube/container"
	// "github.com/phamidko/opencga-operator/pkg/kube/persistentvolumeclaim"
	// "github.com/phamidko/opencga-operator/pkg/kube/podtemplatespec"
	// "github.com/phamidko/opencga-operator/pkg/kube/probes"
	// "github.com/phamidko/opencga-operator/pkg/kube/resourcerequirements"
	// "github.com/phamidko/opencga-operator/pkg/kube/statefulset"
	// "github.com/phamidko/opencga-operator/pkg/util/envvar"
	// "github.com/phamidko/opencga-operator/pkg/util/scale"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"

	ocbv1 "github.com/phamidko/opencga-operator/api/v1"
	"github.com/phamidko/opencga-operator/pkg/kube/statefulset"
	"github.com/phamidko/opencga-operator/pkg/util/scale"
)

const (
	AgentName   = "opencga-agent"
	opencgaName = "opencga"

	versionUpgradeHookName            = "opencga-posthook"
	ReadinessProbeContainerName       = "opencga-agent-readinessprobe"
	readinessProbePath                = "/opt/scripts/readinessprobe"
	agentHealthStatusFilePathEnv      = "AGENT_STATUS_FILEPATH"
	clusterFilePath                   = "/var/lib/automation/config/cluster-config.json"
	opencgaDatabaseServiceAccountName = "opencga-database"
	agentHealthStatusFilePathValue    = "/var/log/opencga-mms-automation/healthstatus/agent-health-status.json"

	opencgaRepoUrl = "opencga_REPO_URL"

	headlessAgentEnv           = "HEADLESS_AGENT"
	podNamespaceEnv            = "POD_NAMESPACE"
	automationConfigEnv        = "AUTOMATION_CONFIG_MAP"
	AgentImageEnv              = "AGENT_IMAGE"
	OpencgaImageEnv            = "OPENCGA_IMAGE"
	VersionUpgradeHookImageEnv = "VERSION_UPGRADE_HOOK_IMAGE"
	ReadinessProbeImageEnv     = "READINESS_PROBE_IMAGE"
	ManagedSecurityContextEnv  = "MANAGED_SECURITY_CONTEXT"

	automationMongodConfFileName = "automation-opencga.conf"
	keyfileFilePath              = "/var/lib/opencga-mms-automation/authentication/keyfile"

	automationAgentOptions = " -skipMongoStart -noDaemonize -useLocalOpencgaTools"

	OpencgaUserCommand = `current_uid=$(id -u)
AGENT_API_KEY="$(cat /opencga-automation/agent-api-key/agentApiKey)"
declare -r current_uid
if ! grep -q "${current_uid}" /etc/passwd ; then
sed -e "s/^opencga:/builder:/" /etc/passwd > /tmp/passwd
echo "opencga:x:$(id -u):$(id -g):,,,:/:/bin/bash" >> /tmp/passwd
export NSS_WRAPPER_PASSWD=/tmp/passwd
export LD_PRELOAD=libnss_wrapper.so
export NSS_WRAPPER_GROUP=/etc/group
fi
`
)

// OpenCGADeploymentOwner is an interface which any resource which generates openCGA should implement.
type OpenCGADeploymentOwner interface {
	// ServiceName returns the name of the K8S service the operator will create.
	ServiceName() string
	// GetName returns the name of the resource.
	GetName() string
	// GetNamespace returns the namespace the resource is defined in.
	GetNamespace() string
	// GetOpenCGAVersion returns the version of OpenCGA to be used for this resource
	GetOpenCGAVersion() string
	// AutomationConfigSecretName returns the name of the secret which will contain the automation config.
	AutomationConfigSecretName() string
	// GetUpdateStrategyType returns the UpdateStrategyType of the statefulset.
	GetUpdateStrategyType() appsv1.StatefulSetUpdateStrategyType
	// HasSeparateDataAndLogsVolumes returns whether or not the volumes for data and logs would need to be different.
	HasSeparateDataAndLogsVolumes() bool
	// GetAgentScramKeyfileSecretNamespacedName returns the NamespacedName of the secret which stores the keyfile for the agent.
	GetAgentKeyfileSecretNamespacedName() types.NamespacedName
	// DataVolumeName returns the name that the data volume should have
	DataVolumeName() string
	// LogsVolumeName returns the name that the data volume should have
	LogsVolumeName() string

	// GetOpenCGAConfiguration returns the OpenCGA configuration for each member.
	GetOpenCGAConfiguration() ocbv1.OpenCGAConfiguration

	// NeedsAutomationConfigVolume returns whether the statefuslet needs to have a volume for the automationconfig.
	NeedsAutomationConfigVolume() bool
}

// BuildMongoDBReplicaSetStatefulSetModificationFunction builds the parts of the replica set that are common between every resource that implements
// MongoDBStatefulSetOwner.
// It doesn't configure TLS or additional containers/env vars that the statefulset might need.
func BuildOpenCGABReplicaSetDeploymentModificationFunction(ocb OpenCGADeploymentOwner, scaler scale.ReplicaSetScaler) statefulset.Modification {
	labels := map[string]string{
		"app": ocb.ServiceName(),
	}
	// healthStatusVolume :=

	return labels
}
