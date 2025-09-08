package constants

// Auth Messages
const (
	LoginSuccess         = "Login berhasil"
	LoginFailed          = "Email/username atau password salah"
	RefreshSuccess       = "Token berhasil diperbarui"
	RefreshFailed        = "Token refresh tidak valid"
	LogoutSuccess        = "Logout berhasil"
	LogoutFailed         = "Logout gagal"
	OTPSent              = "Jika email terdaftar, kode OTP telah dikirim"
	OTPVerified          = "Kode OTP berhasil diverifikasi"
	OTPInvalid           = "Kode OTP tidak valid atau telah kedaluwarsa"
	PasswordResetSuccess = "Password berhasil direset"
	PasswordResetFailed  = "Gagal mereset password"
	TokenInvalid         = "Token tidak valid"
	TokenExpired         = "Token telah kedaluwarsa"
	UnauthorizedAccess   = "Akses tidak diizinkan"
)

// User Messages
const (
	UserListSuccess    = "Data pengguna berhasil diambil"
	UserDetailSuccess  = "Detail pengguna berhasil diambil"
	UserCreateSuccess  = "Pengguna berhasil dibuat"
	UserUpdateSuccess  = "Pengguna berhasil diperbarui"
	UserDeleteSuccess  = "Pengguna berhasil dihapus"
	UserNotFound       = "Pengguna tidak ditemukan"
	UserCreateFailed   = "Gagal membuat pengguna"
	UserUpdateFailed   = "Gagal memperbarui pengguna"
	UserDeleteFailed   = "Gagal menghapus pengguna"
	EmailAlreadyExists = "Email sudah terdaftar"
	UsernameExists     = "Username sudah digunakan"
)

// School Messages
const (
	SchoolListSuccess   = "Data sekolah berhasil diambil"
	SchoolDetailSuccess = "Detail sekolah berhasil diambil"
	SchoolCreateSuccess = "Sekolah berhasil dibuat"
	SchoolUpdateSuccess = "Sekolah berhasil diperbarui"
	SchoolDeleteSuccess = "Sekolah berhasil dihapus"
	SchoolNotFound      = "Sekolah tidak ditemukan"
	SchoolCreateFailed  = "Gagal membuat sekolah"
	SchoolUpdateFailed  = "Gagal memperbarui sekolah"
	SchoolDeleteFailed  = "Gagal menghapus sekolah"

	// Majority Messages
	MajorityListSuccess   = "Data jurusan berhasil diambil"
	MajorityDetailSuccess = "Detail jurusan berhasil diambil"
	MajorityCreateSuccess = "Jurusan berhasil dibuat"
	MajorityUpdateSuccess = "Jurusan berhasil diperbarui"
	MajorityDeleteSuccess = "Jurusan berhasil dihapus"
	MajorityNotFound      = "Jurusan tidak ditemukan"

	// Class Messages
	ClassListSuccess   = "Data kelas berhasil diambil"
	ClassDetailSuccess = "Detail kelas berhasil diambil"
	ClassGetSuccess    = "Data kelas berhasil diambil"
	ClassGetAllSuccess = "Data semua kelas berhasil diambil"
	ClassCreateSuccess = "Kelas berhasil dibuat"
	ClassUpdateSuccess = "Kelas berhasil diperbarui"
	ClassDeleteSuccess = "Kelas berhasil dihapus"
	ClassNotFound      = "Kelas tidak ditemukan"

	// Partner Messages
	PartnerListSuccess   = "Data mitra berhasil diambil"
	PartnerDetailSuccess = "Detail mitra berhasil diambil"
	PartnerGetSuccess    = "Data mitra berhasil diambil"
	PartnerGetAllSuccess = "Data semua mitra berhasil diambil"
	PartnerCreateSuccess = "Mitra berhasil dibuat"
	PartnerUpdateSuccess = "Mitra berhasil diperbarui"
	PartnerDeleteSuccess = "Mitra berhasil dihapus"
	PartnerNotFound      = "Mitra tidak ditemukan"
)

// RBAC Messages
const (
	RoleListSuccess        = "Data role berhasil diambil"
	RoleDetailSuccess      = "Detail role berhasil diambil"
	RoleCreateSuccess      = "Role berhasil dibuat"
	RoleUpdateSuccess      = "Role berhasil diperbarui"
	RoleDeleteSuccess      = "Role berhasil dihapus"
	RoleNotFound           = "Role tidak ditemukan"
	PermissionListSuccess  = "Data permission berhasil diambil"
	PermissionNotFound     = "Permission tidak ditemukan"
	MenuListSuccess        = "Data menu berhasil diambil"
	MenuNotFound           = "Menu tidak ditemukan"
	UserRoleAssigned       = "Role berhasil diberikan kepada pengguna"
	UserRoleRevoked        = "Role berhasil dicabut dari pengguna"
	InsufficientPermission = "Anda tidak memiliki izin untuk mengakses resource ini"
)

// General Messages
const (
	InternalServerError = "Terjadi kesalahan pada server"
	BadRequest          = "Permintaan tidak valid"
	ValidationError     = "Data yang dikirim tidak valid"
	NotFound            = "Data tidak ditemukan"
	ConflictError       = "Data sudah ada atau konflik"
	Success             = "Operasi berhasil dilakukan"
)
