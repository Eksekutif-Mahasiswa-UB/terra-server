package repository

import (
	"time"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	"github.com/jmoiron/sqlx"
)

// ProgramRepository defines the interface for program data operations
type ProgramRepository interface {
	CreateProgram(program *entity.Program) error
	GetProgramByID(id string) (*entity.Program, error)
	GetAllPrograms() ([]entity.Program, error)
	UpdateProgram(program *entity.Program) error
	DeleteProgram(id string) error
}

// programRepository is the concrete implementation of ProgramRepository
type programRepository struct {
	db *sqlx.DB
}

// NewProgramRepository creates a new instance of ProgramRepository
func NewProgramRepository(db *sqlx.DB) ProgramRepository {
	return &programRepository{db: db}
}

// CreateProgram inserts a new program into the database
func (r *programRepository) CreateProgram(program *entity.Program) error {
	program.CreatedAt = time.Now()
	program.UpdatedAt = time.Now()

	query := `INSERT INTO programs (id, title, description, image_url, created_at, updated_at) 
			  VALUES (:id, :title, :description, :image_url, :created_at, :updated_at)`

	_, err := r.db.NamedExec(query, program)
	return err
}

// GetProgramByID retrieves a program by its ID
func (r *programRepository) GetProgramByID(id string) (*entity.Program, error) {
	var program entity.Program
	query := `SELECT * FROM programs WHERE id = ?`

	err := r.db.Get(&program, query, id)
	if err != nil {
		return nil, err
	}

	return &program, nil
}

// GetAllPrograms retrieves all programs from the database
func (r *programRepository) GetAllPrograms() ([]entity.Program, error) {
	var programs []entity.Program
	query := `SELECT * FROM programs ORDER BY created_at DESC`

	err := r.db.Select(&programs, query)
	if err != nil {
		return nil, err
	}

	return programs, nil
}

// UpdateProgram updates an existing program in the database
func (r *programRepository) UpdateProgram(program *entity.Program) error {
	program.UpdatedAt = time.Now()

	query := `UPDATE programs 
			  SET title = :title, description = :description, image_url = :image_url, updated_at = :updated_at 
			  WHERE id = :id`

	_, err := r.db.NamedExec(query, program)
	return err
}

// DeleteProgram deletes a program from the database
func (r *programRepository) DeleteProgram(id string) error {
	query := `DELETE FROM programs WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}


