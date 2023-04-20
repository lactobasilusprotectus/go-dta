package db

import (
	"fmt"
	"go-dts/pkg/common/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DatabaseConnection is the connection to the database
type DatabaseConnection struct {
	Master *gorm.DB
	Slave  *gorm.DB
}

// NewDatabaseConnection constructs new DatabaseConnection
func NewDatabaseConnection(config config.DatabaseConfig) (*DatabaseConnection, error) {
	connStr := ConnStr(config)

	master, err := connect(connStr, config.MaxIdleConnections, config.MaxOpenConnections)
	if err != nil {
		return nil, err
	}

	slave, err := connect(connStr, config.MaxIdleConnections, config.MaxOpenConnections)
	if err != nil {
		return nil, err
	}

	return &DatabaseConnection{
		Master: master,
		Slave:  slave,
	}, nil
}

// ConnStr construct db connection string
func ConnStr(config config.DatabaseConfig) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		config.Username, config.Password, config.Database, config.Host, config.Port)
}

// connect connects to database, given the configuration
func connect(address string, maxIdle, maxOpen int) (*gorm.DB, error) {
	conn, err := gorm.Open(postgres.Open(address), &gorm.Config{})

	if err != nil {
		log.Println("Error while connecting to database at", address, err.Error())
		return nil, err
	}

	//Set max idle and max open connections
	sqlDB, err := conn.DB()

	if err != nil {
		log.Println("Error while getting sql db", err.Error())
		return nil, err
	}

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(maxOpen)

	return conn, nil
}

// CloseConnections closes database connections
func CloseConnections(dbConn *DatabaseConnection) {
	sqlDBMaster, err := dbConn.Master.DB()
	if err != nil {
		log.Println("Error while getting sql db", err.Error())
	}

	if err = sqlDBMaster.Close(); err != nil {
		log.Println("Error while closing sql db", err.Error())
	}

	sqlDBSlave, err := dbConn.Slave.DB()
	if err != nil {
		log.Println("Error while getting sql db", err.Error())
	}

	if err = sqlDBSlave.Close(); err != nil {
		log.Println("Error while closing sql db", err.Error())
	}
}
