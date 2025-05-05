package common

func ConvertToOffsetLimit(page, size int) (offset, limit int) {
	if page <= 0 || size <= 0 {
		return
	}

	offset = (page - 1) * size
	limit = size
	return
}
