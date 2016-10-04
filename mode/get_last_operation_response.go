package mode

const (
  StateSucceeded = "succeeded"
  StateFailed = "failed"
  StateInProgress = "in progress"
)

type GetLastOperationResponse struct {
  State string `json:"state"`
}
