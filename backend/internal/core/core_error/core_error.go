package coreerror

import "errors"

var (
	ErrNotFound          = errors.New("not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrInvalidInputSafe  = errors.New("invalid input : contains invalid characters")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrConflict          = errors.New("conflict")
	ErrDbError           = errors.New("db error")
	ErrUserHasSession    = errors.New("user has session")
	ErrInvalidDateFormat = errors.New("invalid date format")

	ErrInvalidToken          = errors.New("invalid token")
	ErrInvalidCreds          = errors.New("invalid credentials")
	ErrNoTokenProvided       = errors.New("no token provided")
	ErrSessionExpired        = errors.New("session expired")
	ErrHasSession            = errors.New("user already has session")
	ErrNoSessionId           = errors.New("no session id")
	ErrInvalidActionActivity = errors.New("invalid action activity")
	ErrInvalidIpAddress      = errors.New("invalid ip address")

	ErrMissingId    = errors.New("missing id")
	ErrMissingField = errors.New("missing field")

	ErrNoFieldToUpdate       = errors.New("no field to update")
	ErrUsernameTaken         = errors.New("username taken")
	ErrUsernameLengthInvalid = errors.New("username length invalid")
	ErrEmailTaken            = errors.New("email taken")
	ErrNoHpTaken             = errors.New("phone number taken")
	ErrNipTaken              = errors.New("NIP taken")
	ErrNisnTaken             = errors.New("NISN taken")
	ErrInvalidStatusAkun     = errors.New("invalid status akun")

	ErrTingkatKelasExist = errors.New("tingkat kelas already exist")
	ErrNamaKelasExist    = errors.New("nama kelas already exist")
	ErrKodeMapelExist    = errors.New("kode mapel already exist")

	ErrKodeRuangUjianExist = errors.New("kode ruang ujian already exist")
	ErrSesiUjianExist      = errors.New("sesi ujian already exist")

	ErrFileTooLarge      = errors.New("file too large")
	ErrInvalidFileFormat = errors.New("invalid file format")

	ErrDeleteRestricted = errors.New("delete restricted")

	ErrInvalidRequestBody = errors.New("invalid request body")
	ErrMustBePointer      = errors.New("must be pointer")

	ErrPathTraversal = errors.New("path traversal")

	ErrImportJobNotFound = errors.New("import job not found")
	ErrInvalidDocxFormat = errors.New("invalid docx format")
	ErrParsingFailed     = errors.New("parsing failed")
	ErrBankSoalNotFound  = errors.New("bank soal not found")

	ErrContentTypeMustJson      = errors.New("content type must be application/json")
	ErrContentTypeMustMultipart = errors.New("content type must be multipart/form-data")
	ErrInvalidMultipartForm     = errors.New("invalid multipart form")

	ErrTimeEmpty = errors.New("time is empty")

	ErrPesertaInvalid = errors.New("Siswa not found in peserta ujian")
	ErrWaktuAttemptPesertaInvalid = errors.New("Expired attempt time (waktu attempt lewat dari waktu selesai ujian)")
	ErrTokenUjianInvalid = errors.New("Invalid token ujian")
	ErrPesertaNotAllowedToAttemptJadwal = errors.New("Peserta ujian is not allowed to attempt this ujian")
	ErrMissingTokenUjian = errors.New("Missing token ujian")
)
