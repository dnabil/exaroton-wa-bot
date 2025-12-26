package dto

type ExarotonAccountInfo struct {
	// Name represents the account's name.
	Name string `json:"name"`

	// Email represents the account's email.
	Email string `json:"email"`

	// Verified represents whether the account is verified.
	Verified bool `json:"verified"`

	// Credits represents the account's credits.
	Credits float64 `json:"credits"`
}

type ExarotonServerInfo struct {
	// ID represents the unique server ID.
	ID string `json:"id"`

	// Name represents the server name.
	Name string `json:"name"`

	// Address represents the full server address.
	Address string `json:"address"`

	// MOTD represents the MOTD of the server.
	Motd string `json:"motd"`

	// Status represents the current status of the server.
	Status ServerStatus `json:"status"`

	// Host represents the host machine the server is running on.
	// Only available if the server is online.
	Host *string `json:"host"`

	// Port represents the port the server is listening on.
	// Only available if the server is online.
	Port *int `json:"port"`

	// Players represents the players-related information.
	Players ExarotonServerPlayers `json:"players"`

	// Software represents the software-related information.
	Software *ServerSoftware `json:"software"`

	// Shared represents whether the server is accessed via the Share Access feature.
	Shared bool `json:"shared"`
}

// ExarotonServerPlayers represents the information about players on a server.
type ExarotonServerPlayers struct {
	// Max represents the maximum player count (slots).
	Max int `json:"max"`
	// Online represents the current player count.
	Count int `json:"count"`
	// List represents the current player list.
	// Only available if the server is online.
	List []string `json:"list"`
}

const (
	// ServerStatusOffline represents a server that is offline.
	ServerStatusOffline ServerStatus = 0

	// ServerStatusOnline represents a server that is online.
	ServerStatusOnline ServerStatus = 1

	// ServerStatusStarting represents a server that is starting.
	ServerStatusStarting ServerStatus = 2

	// ServerStatusStopping represents a server that is stopping.
	ServerStatusStopping ServerStatus = 3

	// ServerStatusRestarting represents a server that is restarting.
	ServerStatusRestarting ServerStatus = 4

	// ServerStatusSaving represents a server that is saving.
	ServerStatusSaving ServerStatus = 5

	// ServerStatusLoading represents a server that is loading.
	ServerStatusLoading ServerStatus = 6

	// ServerStatusCrashed represents a server that has crashed.
	ServerStatusCrashed ServerStatus = 7

	// ServerStatusPending represents a server that is pending.
	ServerStatusPending ServerStatus = 8

	// ServerStatusTransferring represents a server that is transferring.
	ServerStatusTransferring ServerStatus = 9

	// ServerStatusPreparing represents a server that is preparing.
	ServerStatusPreparing ServerStatus = 10
)

// ServerStatus represents the status of a server.
type ServerStatus uint8

// String returns the string representation of the server status.
func (s ServerStatus) String() string {
	switch s {
	case ServerStatusOffline:
		return "offline"
	case ServerStatusOnline:
		return "online"
	case ServerStatusStarting:
		return "starting"
	case ServerStatusStopping:
		return "stopping"
	case ServerStatusRestarting:
		return "restarting"
	case ServerStatusSaving:
		return "saving"
	case ServerStatusLoading:
		return "loading"
	case ServerStatusCrashed:
		return "crashed"
	case ServerStatusPending:
		return "pending"
	case ServerStatusTransferring:
		return "transferring"
	case ServerStatusPreparing:
		return "preparing"
	default:
		return "unknown server status"
	}
}

// ServerSoftware represents the information about installed server software.
type ServerSoftware struct {
	// ID represents the unique software ID.
	ID string `json:"id"`
	// Name represents the software name.
	Name string `json:"name"`
	// Version represents the software version.
	Version string `json:"version"`
}
