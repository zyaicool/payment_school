package utilities

func ChangeStringSortBySchoolClass(sortBy string) string {
	switch sortBy {
	case "invoiceNumber":
		sortBy = "invoiceNumber"
	case "studentName":
		sortBy = "studentName"
	case "paymentDate":
		sortBy = "paymentDate"
	case "paymentMethod":
		sortBy = "paymentMethod"
	case "username":
		sortBy = "username"
	case "schoolGradeName":
		sortBy = "schoolGradeName"
	case "schoolClassName":
		sortBy = "schoolClassName"
	case "totalAmount":
		sortBy = "totalAmount"
	case "transactionStatus":
		sortBy = "transactionStatus"
	default:
	}
	return sortBy
}

func ChangeStringSortByPaymentReport(sortBy string) string {
	switch sortBy {
	case "unit":
		sortBy = "schoolGradeId"
	case "prefixClass":
		sortBy = "prefixClassId"
	case "schoolMajor":
		sortBy = "schoolMajorId"
	default:
	}
	return sortBy
}

func ChangeStringSortByBillingHistory(sortBy string) string {
	switch sortBy {
	case "studentName":
		return "student_name"
	case "totalAmount":
		return "total_amount"
	case "createdDate":
		return "created_date"

	default:
		return sortBy
	}
}

func ChangeStringSortByBillingStudent(sortBy string) string {
	switch sortBy {
	case "schoolGrade":
		sortBy = "schoolGradeId"
	case "schoolClass":
		sortBy = "schoolClassId"
	case "studentName":
		sortBy = "s.full_name"
	case "billingDetailName":
		sortBy = "detailBillingName"
	default:
	}
	return sortBy
}

func ChangeStringSortByUser(sortBy string) string {
	switch sortBy {
	case "roleName":
		return "roles.name"
	case "createdDate":
		return "users.created_at"
	case "status":
		return "users.is_block"

	default:
		return sortBy
	}
}

func ChangeStringSortByAnnouncement(sortBy string) string {
	switch sortBy {
	case "createdAt":
		return "created_at"
	case "createdBy":
		return "created_by"
	case "updatedAt":
		return "updated_at"
	case "updatedBy":
		return "updated_by"
	default:
		return sortBy
	}
}

func ChangeStringSortByBillingReport(sortBy string) string {
	switch sortBy {
	case "detailBillingName":
		sortBy = "detail_billing_name"
	case "billingType":
		sortBy = "billing_type"
	case "studentName":
		sortBy = "student_name"
	case "schoolGradeName":
		sortBy = "school_grade_name"
	case "schoolClassName":
		sortBy = "school_class_name"
	case "schoolYearName":
		sortBy = "school_year_name"
	case "bankAccountName":
		sortBy = "bank_name"
	case "paymentStatus":
		sortBy = "payment_status"
	default:
	}
	return sortBy
}

func ChangeStringSortByStudent(sortBy string) string {
	switch sortBy {
	case "schoolClass":
		sortBy = "school_class_name"
	default:
	}
	return sortBy
}
