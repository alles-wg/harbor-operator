package database

import (
	"fmt"

	"github.com/goharbor/harbor-cluster-operator/controllers/database/api"

	//pg "github.com/zalando/postgres-operator/pkg/apis/acid.zalan.do/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	DefaultDatabaseReplica = 3
	DefaultDatabaseMemory  = "1Gi"
	DefaultDatabaseVersion = "12"
)

func (p *PostgreSQLReconciler) GetDatabases() map[string]string {
	databases := map[string]string{
		CoreDatabase: "zalando",
	}

	if p.HarborCluster.Spec.Notary != nil {
		databases[NotaryServerDatabase] = "zalando"
		databases[NotarySignerDatabase] = "zalando"
	}

	return databases
}

// GetDatabaseConn is getting database connection
func (p *PostgreSQLReconciler) GetDatabaseConn(secretName string) (*Connect, error) {
	secret, err := p.GetSecret(secretName)
	if err != nil {
		return nil, err
	}

	conn := &Connect{
		Host:     string(secret["host"]),
		Port:     string(secret["port"]),
		Password: string(secret["password"]),
		Username: string(secret["username"]),
		Database: string(secret["database"]),
	}

	return conn, nil
}

// GetSecret returns the database connection Secret
func (p *PostgreSQLReconciler) GetSecret(secretName string) (map[string][]byte, error) {
	secret := &corev1.Secret{}
	err := p.Client.Get(types.NamespacedName{Name: secretName, Namespace: p.HarborCluster.Namespace}, secret)
	if err != nil {
		return nil, err
	}
	data := secret.Data
	return data, nil
}

// GetPostgreResource returns postgres resource
func (p *PostgreSQLReconciler) GetPostgreResource() api.Resources {
	resources := api.Resources{}

	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec == nil {
		resources.ResourceRequests = api.ResourceDescription{
			CPU:    "1",
			Memory: "1Gi",
		}
		resources.ResourceRequests = api.ResourceDescription{
			CPU:    "2",
			Memory: "2Gi",
		}
		return resources
	}

	cpu := p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Resources.Requests.Cpu()
	mem := p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Resources.Requests.Memory()

	request := api.ResourceDescription{}
	if cpu != nil {
		request.CPU = cpu.String()
	}
	if mem != nil {
		request.Memory = mem.String()
	}
	resources.ResourceRequests = request
	resources.ResourceLimits = request

	return resources
}

// GetRedisServerReplica returns postgres replicas
func (p *PostgreSQLReconciler) GetPostgreReplica() int32 {
	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec == nil {
		return DefaultDatabaseReplica
	}

	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Replicas == 0 {
		return DefaultDatabaseReplica
	}
	return int32(p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Replicas)
}

// GetPostgreStorageSize returns Postgre storage size
func (p *PostgreSQLReconciler) GetPostgreStorageSize() string {
	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec == nil {
		return DefaultDatabaseMemory
	}

	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Storage == "" {
		return DefaultDatabaseMemory
	}
	return p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Storage
}

func (p *PostgreSQLReconciler) GetPostgreVersion() string {
	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec == nil {
		return DefaultDatabaseVersion
	}

	if p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Version == "" {
		return DefaultDatabaseVersion
	}

	return p.HarborCluster.Spec.InClusterDatabase.PostgresSQLSpec.Version
}

// GenDatabaseUrl returns database connection url
func (c *Connect) GenDatabaseUrl() string {
	databaseURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", c.Username, c.Password, c.Host, c.Port, c.Database)
	return databaseURL
}
