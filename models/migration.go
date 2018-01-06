package models

import "time"

type Migrations struct {
	Indexes  MigrationInfo `json:"indexes"`
	LastTime time.Time     `json:"date"`
}

type MigrationInfo struct {
	Migration map[string]IndexInfo `json:"migration"`
}

type IndexInfo struct {
	Info map[string]string `json:"info"`
}

// Stores data about migration Upgrade
/*type upgradeMigration struct {
	Cr
}*/
