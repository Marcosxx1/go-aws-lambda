// Package mysqlservice provides functionality for interacting with MySQL database.
package mysqlservice

import (
	"database/sql"
	"fmt"
	"test/lambda/interfaces"
	mysqlconfig "test/lambda/services/mysql-service/config"
	"time"
)

// MysqlTabloideRepository represents a repository for interacting with a MySQL database.
type MysqlTabloideRepository struct {
	connection *sql.DB // The underlying SQL database connection.
	tableName  string  // The name of the table in the database.
}

// NewMysqlTabloideRepository creates a new instance of MysqlTabloideRepository.
// It initializes a connection to the MySQL database using the configuration from mysql_config package.
// It returns a pointer to the MysqlTabloideRepository.
func NewMysqlTabloideRepository() *MysqlTabloideRepository {
	c := mysqlconfig.NewMysqlDatabase().GetConn()
	tableName := "tabloide"
	return &MysqlTabloideRepository{
		connection: c,
		tableName:  tableName,
	}
}

// InsertTabloid inserts a new tabloid record into the database.
// It takes name, regionID, startValidityDate, endValidityDate as input parameters.
// It also takes a transaction object for performing the insert operation as part of a larger transaction.
// It returns the ID of the newly inserted record or an error if the operation fails.
//
// Example:
//
//	repository := NewMysqlTabloideRepository()
//	transaction, err := repository.connection.Begin()
//	if err != nil {
//	    log.Fatalf("Failed to start transaction: %v", err)
//	}
//	defer transaction.Rollback()
//
//	name := "Sample Tabloid"
//	regionID := 1
//	startValidityDate := time.Now()
//	endValidityDate := time.Now().AddDate(0, 1, 0) // Valid for 1 month
//
//	lastID, err := repository.InsertTabloid(name, regionID, startValidityDate, endValidityDate, transaction)
//	if err != nil {
//	    log.Fatalf("Failed to insert tabloid: %v", err)
//	}
//	fmt.Printf("Tabloid inserted successfully with ID: %d\n", lastID)
//
//	err = transaction.Commit()
//	if err != nil {
//	    log.Fatalf("Failed to commit transaction: %v", err)
//	}
func (r *MysqlTabloideRepository) InsertTabloid(name string, regionID int, startValidityDate, endValidityDate time.Time, transaction *sql.Tx) (int64, error) {
	query := `INSERT INTO ` + r.tableName + `
        (nome, regiao_id, dt_inicio_vigencia, dt_fim_vigencia, ativo, dt_cadastro) 
    VALUES 
        (?, ?, ?, ?, 1, NOW())
    `

	result, err := transaction.Exec(query, name, regionID, startValidityDate, endValidityDate)
	if err != nil {
		return 0, fmt.Errorf("failed to execute query: %v", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get last insert ID: %v", err)
	}

	return lastID, nil
}

// InsertTabloidImage inserts an image URL associated with a tabloid into the database.
// It takes imageURL, tabloidID, order as input parameters and a transaction object for performing the insert operation as part of a larger transaction.
// It returns an error if the operation fails.
//
// Example:
//
//	repository := NewMysqlTabloideRepository()
//	transaction, err := repository.connection.Begin()
//	if err != nil {
//	    log.Fatalf("Failed to start transaction: %v", err)
//	}
//	defer transaction.Rollback()
//
//	imageURL := "https://example.com/image.jpg"
//	tabloidID := 1
//	order := 1
//
//	err := repository.InsertTabloidImage(imageURL, tabloidID, order, transaction)
//	if err != nil {
//	    log.Fatalf("Failed to insert tabloid image: %v", err)
//	}
//	fmt.Println("Tabloid image inserted successfully.")
//
//	err = transaction.Commit()
//	if err != nil {
//	    log.Fatalf("Failed to commit transaction: %v", err)
//	}
func (r *MysqlTabloideRepository) InsertTabloidImage(imageURL string, tabloidID int64, order int, transaction *sql.Tx) error {
	query :=
		`INSERT INTO imagem_tabloide 
		(imagem_url, tabloide_id, ordem, dt_cadastro) 
		VALUES ( ?, ?, ?, NOW())`

	result, err := transaction.Exec(query, imageURL, tabloidID, order)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}
	fmt.Println("Rows affected:", rowsAffected)
	return nil
}

// GetRegionById retrieves a region from the database by its ID.
// It takes regionID as input parameter and returns the corresponding region object or an error if the operation fails.
//
// Example:
//
//	repository := NewMysqlTabloideRepository()
//
//	regionID := 1
//	region, err := repository.GetRegionById(regionID)
//	if err != nil {
//	    log.Fatalf("Failed to retrieve region: %v", err)
//	}
//	fmt.Printf("Region details - ID: %d, Name: %s, Creation Date: %s, Last Updated: %s\n",
//	    region.ID, region.Nome, region.Dt_cadastro.Format("2006-01-02"), region.Dt_alteracao.Format("2006-01-02"))
func (r *MysqlTabloideRepository) GetRegionById(regionID int) (*interfaces.Region, error) {
	var region interfaces.Region
	var dtCadastro, dtAlteracao []uint8

	query := "SELECT id, nome, dt_cadastro, dt_alteracao FROM regiao WHERE id = ? LIMIT 1"

	err := r.connection.QueryRow(query, regionID).Scan(&region.ID, &region.Nome, &dtCadastro, &dtAlteracao)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}

	dtCadastroTime, err := time.Parse("2006-01-02 15:04:05", string(dtCadastro))
	if err != nil {
		return nil, fmt.Errorf("failed to parse dt_cadastro: %v", err)
	}
	region.Dt_cadastro = dtCadastroTime

	dtAlteracaoTime, err := time.Parse("2006-01-02 15:04:05", string(dtAlteracao))
	if err != nil {
		return nil, fmt.Errorf("failed to parse dt_alteracao: %v", err)
	}
	region.Dt_alteracao = dtAlteracaoTime

	return &region, nil
}

func (r *MysqlTabloideRepository) GetTransaction() (*sql.Tx, error) {
	tx, err := r.connection.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	return tx, nil
}

func (r *MysqlTabloideRepository) CommitTransaction(transaction *sql.Tx) error {
	return transaction.Commit()
}

func (r *MysqlTabloideRepository) RollbackTransaction(transaction *sql.Tx) error {
	return transaction.Rollback()
}
