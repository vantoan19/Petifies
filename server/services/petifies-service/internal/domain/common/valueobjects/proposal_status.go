package valueobjects

type PetifiesProposalStatus string

const (
	PetifiesProposalStatusWaitingForAcceptance PetifiesProposalStatus = "PETIFIES_PROPOSAL_STATUS_WAITING_FOR_ACCEPTANCE"
	PetifiesProposalStatusAccepted             PetifiesProposalStatus = "PETIFIES_PROPOSAL_STATUS_ACCEPTED"
	PetifiesProposalStatusCancelled            PetifiesProposalStatus = "PETIFIES_PROPOSAL_STATUS_CANCELLED"
	PetifiesProposalStatusRejected             PetifiesProposalStatus = "PETIFIES_PROPOSAL_STATUS_REJECTED"
	PetifiesProposalStatusSessionClosed        PetifiesProposalStatus = "PETIFIES_PROPOSAL_STATUS_SESSION_CLOSED"
)
