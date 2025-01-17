package models
import (
    "errors"
)
// Provides operations to manage the collection of agreementAcceptance entities.
type CallRecordingStatus int

const (
    SUCCESS_CALLRECORDINGSTATUS CallRecordingStatus = iota
    FAILURE_CALLRECORDINGSTATUS
    INITIAL_CALLRECORDINGSTATUS
    CHUNKFINISHED_CALLRECORDINGSTATUS
    UNKNOWNFUTUREVALUE_CALLRECORDINGSTATUS
)

func (i CallRecordingStatus) String() string {
    return []string{"success", "failure", "initial", "chunkFinished", "unknownFutureValue"}[i]
}
func ParseCallRecordingStatus(v string) (interface{}, error) {
    result := SUCCESS_CALLRECORDINGSTATUS
    switch v {
        case "success":
            result = SUCCESS_CALLRECORDINGSTATUS
        case "failure":
            result = FAILURE_CALLRECORDINGSTATUS
        case "initial":
            result = INITIAL_CALLRECORDINGSTATUS
        case "chunkFinished":
            result = CHUNKFINISHED_CALLRECORDINGSTATUS
        case "unknownFutureValue":
            result = UNKNOWNFUTUREVALUE_CALLRECORDINGSTATUS
        default:
            return 0, errors.New("Unknown CallRecordingStatus value: " + v)
    }
    return &result, nil
}
func SerializeCallRecordingStatus(values []CallRecordingStatus) []string {
    result := make([]string, len(values))
    for i, v := range values {
        result[i] = v.String()
    }
    return result
}
