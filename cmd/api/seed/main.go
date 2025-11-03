package main

import (
	"log"

	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/auth/repository"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/config"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/database"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/domain/entity"
	programRepo "github.com/Eksekutif-Mahasiswa-UB/terra-server/internal/program/repository"
	"github.com/Eksekutif-Mahasiswa-UB/terra-server/pkg/hash"
	"github.com/google/uuid"
)

func main() {
	log.Println("Starting database seeder...")

	// Load configuration
	config.LoadConfig()

	// Initialize database connection
	db, err := database.NewMySQLConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	authRepo := repository.NewAuthRepository(db)
	programRepository := programRepo.NewProgramRepository(db)

	// Seed Admin User
	log.Println("Seeding admin user...")
	if err := seedAdminUser(authRepo, *config.AppConfig); err != nil {
		log.Printf("Failed to seed admin user: %v", err)
	} else {
		log.Println("✓ Admin user seeded successfully")
	}

	// Seed Programs
	log.Println("Seeding programs...")
	if err := seedPrograms(programRepository); err != nil {
		log.Printf("Failed to seed programs: %v", err)
	} else {
		log.Println("✓ Programs seeded successfully")
	}

	log.Println("Database seeding completed!")
}

// seedAdminUser creates a dummy admin user
func seedAdminUser(authRepo repository.AuthRepository, cfg config.Config) error {
	// Hash password from config
	hashedPassword, err := hash.HashPassword(cfg.SEEDER_ADMIN_PASS)
	if err != nil {
		return err
	}

	// Create admin user
	adminUser := &entity.User{
		ID:         uuid.NewString(),
		FullName:   "Admin Terra",
		Email:      "admin@terra.com",
		Password:   hashedPassword,
		Role:       "admin",
		AuthMethod: "email",
		GoogleID:   nil,
	}

	// Save to database
	if err := authRepo.CreateUser(adminUser); err != nil {
		return err
	}

	log.Printf("  → Admin user created with email: %s", adminUser.Email)
	return nil
}

// seedPrograms creates dummy programs
func seedPrograms(programRepo programRepo.ProgramRepository) error {
	programs := []entity.Program{
		{
			ID:          uuid.NewString(),
			Title:       "Waste Bank (Bank Sampah)",
			Description: "Program pengelolaan sampah yang bertujuan untuk mengurangi volume sampah dan memberikan nilai ekonomis bagi masyarakat. Melalui program ini, warga dapat menukarkan sampah yang telah dipilah dengan poin yang dapat ditukar dengan uang atau barang kebutuhan sehari-hari.",
			ImageURL:    "",
		},
		{
			ID:          uuid.NewString(),
			Title:       "Reforestation Program",
			Description: "Program penanaman pohon untuk merestorasi hutan dan lahan kritis. Kegiatan ini melibatkan masyarakat lokal dalam upaya pelestarian lingkungan dan peningkatan tutupan hijau di wilayah sekitar.",
			ImageURL:    "",
		},
		{
			ID:          uuid.NewString(),
			Title:       "Clean Water Initiative",
			Description: "Inisiatif penyediaan akses air bersih untuk komunitas yang membutuhkan. Program ini mencakup pembangunan sumur, sistem filtrasi air, dan edukasi tentang pentingnya air bersih bagi kesehatan.",
			ImageURL:    "",
		},
	}

	for _, program := range programs {
		if err := programRepo.CreateProgram(&program); err != nil {
			return err
		}
		log.Printf("  → Program created: %s", program.Title)
	}

	return nil
}
