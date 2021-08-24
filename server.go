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
