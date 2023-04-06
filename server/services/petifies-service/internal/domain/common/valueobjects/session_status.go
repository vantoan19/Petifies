package valueobjects

type PetifiesSessionStatus string

const (
	PetifiesSessionStatusWaitingForProposal PetifiesSessionStatus = "PETIFIES_SESSION_STATUS_WAITING_FOR_PROPOSAL"
	PetifiesSessionStatusProposalAccepted   PetifiesSessionStatus = "PETIFIES_SESSION_STATUS_PROPOSAL_ACCEPTED"
	PetifiesSessionStatusOnGoing            PetifiesSessionStatus = "PETIFIES_SESSION_STATUS_ON_GOING"
	PetifiesSessionStatusEnded              PetifiesSessionStatus = "PETIFIES_SESSION_STATUS_ENDED"
)
