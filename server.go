package revoltgo

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
)

// Server struct.
type Server struct {
	Client    *Client
	CreatedAt time.Time

	Id                 string                 `json:"_id"`
	Nonce              string                 `json:"nonce"`
	OwnerId            string                 `json:"owner"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	ChannelIds         []string               `json:"channels"`
	Categories         []*ServerCategory      `json:"categories"`
	SystemMessages     *SystemMessages        `json:"system_messages"`
	Roles              map[string]interface{} `json:"roles"`
	DefaultPermissions []interface{}          `json:"default_permissions"`
	Icon               *Attachment            `json:"icon"`
	Banner             *Attachment            `json:"banner"`
}

// Server categories struct.
type ServerCategory struct {
	Id         string   `json:"id"`
	Title      string   `json:"title"`
	ChannelIds []string `json:"channels"`
}

// System messages struct.
type SystemMessages struct {
	UserJoined string `json:"user_joined,omitempty"`
	UserLeft   string `json:"user_left,omitempty"`
	UserKicked string `json:"user_kicker,omitempty"`
	UserBanned string `json:"user_banned,omitempty"`
}

// Server member struct.
type Member struct {
	Informations struct {
		ServerId string `json:"server"`
		UserId   string `json:"user"`
	} `json:"_id"`
	Nickname string      `json:"nickname"`
	Avatar   *Attachment `json:"avatar"`
	Roles    []string    `json:"roles"`
}

// Fetched members struct.
type FetchedMembers struct {
	Members []*Member `json:"members"`
	Users   []*User   `json:"users"`
}

// Calculate creation date and edit the struct.
func (s *Server) CalculateCreationDate() error {
	ulid, err := ulid.Parse(s.Id)

	if err != nil {
		return err
	}

	s.CreatedAt = time.UnixMilli(int64(ulid.Time()))
	return nil
}

// Edit server.
func (s Server) Edit(es *EditServer) error {
	data, err := json.Marshal(es)

	if err != nil {
		return err
	}

	_, err = s.Client.Request("PATCH", "/servers/"+s.Id, data)

	if err != nil {
		return err
	}

	return nil
}

// Delete / leave server.
// If the server not created by client, it will leave.
// Otherwise it will be deleted.
func (s Server) Delete() error {
	_, err := s.Client.Request("DELETE", "/servers/"+s.Id, []byte{})

	if err != nil {
		return err
	}

	return nil
}

// Create a new text-channel.
func (s Server) CreateTextChannel(name, description string) (*Channel, error) {
	channel := &Channel{}
	channel.Client = s.Client

	data, err := s.Client.Request("POST", "/servers/"+s.Id+"/channels", []byte("{\"type\":\"Text\",\"name\":\""+name+"\",\"description\":\""+description+"\",\"nonce\":\""+genULID()+"\"}"))

	if err != nil {
		return channel, err
	}

	err = json.Unmarshal(data, channel)

	if err != nil {
		return channel, err
	}

	return channel, nil
}

// Create a new voice-channel.
func (s Server) CreateVoiceChannel(name, description string) (*Channel, error) {
	channel := &Channel{}
	channel.Client = s.Client

	data, err := s.Client.Request("POST", "/servers/"+s.Id+"/channels", []byte("{\"type\":\"Voice\",\"name\":\""+name+"\",\"description\":\""+description+"\",\"nonce\":\""+genULID()+"\"}"))

	if err != nil {
		return channel, err
	}

	err = json.Unmarshal(data, channel)

	if err != nil {
		return channel, err
	}

	return channel, nil
}

// Fetch a member from Server.
func (s Server) FetchMember(id string) (*Member, error) {
	member := &Member{}

	data, err := s.Client.Request("GET", "/servers/"+s.Id+"/members/"+id, []byte{})

	if err != nil {
		return member, err
	}

	err = json.Unmarshal(data, member)

	if err != nil {
		return member, err
	}

	return member, nil
}

// Fetch all of the members from Server.
func (s Server) FetchMembers() (*FetchedMembers, error) {
	members := &FetchedMembers{}

	data, err := s.Client.Request("GET", "/servers/"+s.Id+"/members", []byte{})

	if err != nil {
		return members, err
	}

	err = json.Unmarshal(data, members)

	if err != nil {
		return members, err
	}

	// Add client to the user
	for _, i := range members.Users {
		i.Client = s.Client
	}

	return members, nil
}

// Edit a member.
func (s Server) EditMember(id string, em *EditMember) error {
	data, err := json.Marshal(em)

	if err != nil {
		return err
	}

	_, err = s.Client.Request("PATCH", "/servers/"+s.Id+"/members/"+id, data)

	if err != nil {
		return err
	}

	return nil
}

// Kick a member from server.
func (s Server) KickMember(id string) error {
	_, err := s.Client.Request("DELETE", "/servers/"+s.Id+"/members/"+id, []byte{})

	if err != nil {
		return err
	}

	return nil
}

// Ban a member from server.
func (s Server) BanMember(id, reason string) error {
	_, err := s.Client.Request("PUT", "/servers/"+s.Id+"/bans/"+id, []byte("{\"reason\":\""+reason+"\"}"))

	if err != nil {
		return err
	}

	return nil
}

// Unban a member from server.
func (s Server) UnbanMember(id string) error {
	_, err := s.Client.Request("DELETE", "/servers/"+s.Id+"/bans/"+id, []byte{})

	if err != nil {
		return err
	}

	return nil
}

// // Fetch all server invites.
// func (s Server) FetchInvites() {
// 	data, _ := s.Client.Request("GET", "/servers/"+s.Id+"/invites", []byte{})

// 	fmt.Println("\n\n" + string(data) + "\n\n")
// }
