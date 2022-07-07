package config

import "leafgame/pb"

type Config struct {
	Skills []*pb.Skill
	Roles  []*pb.Role
}

func New(config *pb.Config) *Config {
	return &Config{Skills: config.Skills, Roles: config.Roles}
}
