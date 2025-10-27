CREATE TABLE volunteers (
    id VARCHAR(36) PRIMARY KEY,
    nama_lengkap VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    no_whatsapp VARCHAR(20) NOT NULL,
    tanggal_lahir DATE NOT NULL,
    jenis_kelamin ENUM('Laki-laki', 'Perempuan') NOT NULL,
    domisili VARCHAR(255) NOT NULL,
    status ENUM('Mahasiswa', 'Siswa', 'Pekerja') NOT NULL,
    minat_lingkungan VARCHAR(255) NOT NULL,
    nama_sertifikat VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);