package main

import "time"

const EMERGENCY_TIME = time.Hour * 48

type UserInfo struct {
	LastHeartBeartTime time.Time
	Username           string
	Contacts           []Contact
	EmergencyMessage   string
	Key                string
}

func (u *UserInfo) GetTimeout(now time.Time) time.Duration {
	ne := u.LastHeartBeartTime.Add(EMERGENCY_TIME)
	return ne.Sub(now)
}
