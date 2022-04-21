package packet

import (
	"time"

	"github.com/MagalixTechnologies/policy-core/domain"
	"github.com/MagalixTechnologies/uuid-go"
)

const (
	ProtocolMajorVersion = 2
	ProtocolMinorVersion = 4
)

// PacketKind represents supported packet kinds
type PacketKind string

const (
	PacketKindHello                 PacketKind = "hello"
	PacketKindAuthorizationRequest  PacketKind = "authorization/request"
	PacketKindAuthorizationAnswer   PacketKind = "authorization/answer"
	PacketPolicyValidationAudit     PacketKind = "policyValidation/audit"
	PacketPolicyValidationAdmission PacketKind = "policyValidation/admission"
	PacketKindPing                  PacketKind = "ping"
)

func (kind PacketKind) String() string {
	return string(kind)
}

// PacketHello packet sent in the hello phase
type PacketHello struct {
	Major            uint      `json:"major"`
	Minor            uint      `json:"minor"`
	Build            string    `json:"build"`
	AccountID        uuid.UUID `json:"account_id"`
	ClusterID        uuid.UUID `json:"cluster_id"`
	K8sServerVersion string    `json:"k8s_server_version"`
	ClusterProvider  string    `json:"cluster_provider"`
}

// PacketAuthorizationRequest packet sent to request authorization
type PacketAuthorizationRequest struct {
	AccountID uuid.UUID `json:"account_id"`
	ClusterID uuid.UUID `json:"cluster_id"`
}

// PacketAuthorizationQuestion packet sent from backend when requesting authorization
type PacketAuthorizationQuestion struct {
	Token []byte `json:"token"`
}

// PacketAuthorizationAnswer packet sent as identifier for authentication
type PacketAuthorizationAnswer struct {
	Token []byte `json:"token"`
}

// PacketAuthorizationSuccess sent when authentication succeeded
type PacketAuthorizationSuccess struct{}

// PacketPing sent for a ping request
type PacketPing struct {
	Number  int       `json:"number,omitempty"`
	Started time.Time `json:"started"`
}

// PacketPong recieved during a ping request
type PacketPong struct {
	Number  int       `json:"number,omitempty"`
	Started time.Time `json:"started"`
}

// PacketPolicyValidationResults contains the policy validation to be sent to backend
type PacketPolicyValidationResults struct {
	Items     []domain.PolicyValidation `json:"items"`
	Timestamp time.Time                 `json:"timestamp"`
}
