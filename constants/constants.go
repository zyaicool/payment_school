package constants

//---------------parameter--------------------------
const ImageJPEG = "image/jpeg"
const ImagePNG = "image/png"
const ContentType = "Content-Type"
const DateFormatYYYYMMDD = "2006-01-02"
const DateFormatDDMMMYYYhhmm = "02 January 2006 15:04"
const MessageUserCantAccessPage = "User can't access this page"
const ContentTypeExcel = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
const MessageErrorFirebaseClientNotInitialized = "Firebase client not initialized"
const MessageErrorInvalidArgument = "Invalid argument error: %v"
const MessageErrorSenderIdMismatch = "SenderId mismatch error: %v"
const MessageErrorSendingMessage = "General error sending message: %v"
const MessageSuccessSendMessage = "Successfully sent message: %s\n"
const MessageBillingName = "Billing name"
const DummyEmail = "contoh@contoh.com"
const EmailConfirmationText = "Konfirmasi Email"

//---------------message error----------------------
const ErrorTextMessage = "Error: "
const UnsupportedImageFormatMessage = "Unsupported image format: "
const ErrorEncodingImageMessage = "Error encoding image: "
const InvalidBIllingStudentIdMessage = "Invalid billingStudentId"
const DataNotFoundMessage = "Data Not Found"
const CannotParseJsonMessage = "Cannot parse JSON"
const FailedToParseRequestBodyMessage = "Failed to parse request body"

//--------------query---------------------------------
//billing report
const FilterByUscSchoolId = " AND usc.school_id = %d "
const FilterByUsSchoolClassId = " AND sc.id = %d "
const FilterBySgSchoolGradeId = " AND sg.id = %d "
const FilterByBsPaymentStatus = " AND bs.payment_status = '%d' "
const FilterBySySchoolYearId = " AND sy.id = %d "
const FilterByBankAccountId = " AND b.bank_account_id = %d "
const FilterByBillingType = " AND b.billing_type = '%s' "
const FilterByStudentId = " AND s.id = %d "

//dashboard
const JoinUserSchoolsToUserStudents = "JOIN user_schools ON user_schools.user_id = user_students.user_id"

//students
const JoinUserStudentsToStudents = "JOIN user_students ON user_students.student_id = s.id "
const JoinUsersToUserStudents = "JOIN users ON users.id = user_students.user_id "
const JoinUserStudentsToStudentsAndFilterDeletedAt = "JOIN user_students ON user_students.student_id = students.id AND user_students.deleted_at is null"
const JoinUsersToUserStudentsAndFilterDeletedAt = "JOIN users ON users.id = user_students.user_id AND users.deleted_at is null"
const JoinUserSchoolsToUsersAndFilterDeletedAt = "JOIN user_schools ON user_schools.user_id = users.id AND user_schools.deleted_at is null"
const JoinSchoolsToUserSChoolsAndFilterDeletedAt = "JOIN schools ON schools.id = user_schools.school_id AND schools.deleted_at is null"
const FilterByUsersId = "users.id = ?"

//student parent
const JoinUsersToStudentParenstAndFilterDeletedAt = "JOIN users ON users.id = student_parents.user_id AND users.deleted_at is null"

//transaction
const FilterOrderId = "order_id = ?"
const MessageErrorGenerateInvoiceNumber = "failed to generate invoice number: %v"
const MessageErrorConvertTotalPaymenr = "gagal mengonversi TotalPayment ke big.Int, nilai: %s"

//user
const PreloadUserSchoolToSchool = "UserSchool.School"
const PreloadRoleToRoleMatrix = "Role.RoleMatrix"
const MessageErrorLinkExpired = "Link has expired, please try again."

//login
const EmailPasswordSalahMessage = "Email/Username atau password yang Anda masukkan salah. Silahkan coba lagi"

//payment report
const PaymentReportParam = "Laporan Pembayaran"

//schedule
const JsonFirebaseConfigFile = "./data/firebase_config_file.json"
const ErrorFirebaseClientMessage = "Error initializing Firebase client: %v"

//student parent
const ErrorMessageLoginAsParent = "Silahkan login dengan akun orang tua atau wali siswa."

//dashboard
const JoinUserStudentsToStudentsDashboard = "JOIN user_students ON user_students.student_id = students.id"
