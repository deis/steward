package mode

type BindRequest struct {
	ServiceID  string     `json:"service_id"`
	PlanID     string     `json:"plan_id"`
	Parameters JSONObject `json:"parameters"`
}
