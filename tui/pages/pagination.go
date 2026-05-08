package pages

import "fmt"

const (
	quoteListPageSize   = 20
	historyListPageSize = 20
)

func clampPageOffset(totalCount int64, pageSize, offset int) int {
	if totalCount <= 0 || pageSize <= 0 {
		return 0
	}
	if offset < 0 {
		offset = 0
	}
	maxOffset := int((totalCount - 1) / int64(pageSize) * int64(pageSize))
	if offset > maxOffset {
		return maxOffset
	}
	return offset
}

func totalPages(totalCount int64, pageSize int) int {
	if totalCount <= 0 || pageSize <= 0 {
		return 0
	}
	return int((totalCount + int64(pageSize) - 1) / int64(pageSize))
}

func currentPage(offset, pageSize int) int {
	if pageSize <= 0 {
		return 1
	}
	return offset/pageSize + 1
}

func formatPagedTitle(base string, totalCount int64, pageSize, offset int) string {
	if totalCount <= 0 {
		return fmt.Sprintf("%s (0)", base)
	}
	return fmt.Sprintf("%s (%d total, page %d/%d)", base, totalCount, currentPage(offset, pageSize), totalPages(totalCount, pageSize))
}
