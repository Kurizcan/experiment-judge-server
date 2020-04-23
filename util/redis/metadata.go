package redis

import (
	"fmt"
)

// filed: problemId val: status
func GetGroupSidKey(groupId, studentId int) string {
	return fmt.Sprintf("experiment_%d_student_%d", groupId, studentId)
}

// runId_userId:status
func GetRunIdStatusKey(runId, studentId int) string {
	return fmt.Sprintf("run_%d_student_%d", runId, studentId)
}
