package models

// Role
// @Description role model
type Role struct {
	// id of the role
	ID string `json:"id"`
	// name of the role
	Name string `json:"name"`
	// if the role is mentionable
	Mentionable bool `json:"mentionable"`
} // @name Discord.Role

// Channel
// @Description Discord Channel Model
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`
} // @name Discord.Channel

type GuildInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
