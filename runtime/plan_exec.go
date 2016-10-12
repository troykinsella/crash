package runtime

type PlanExec struct {
	plan *Plan
}

func NewPlanExec(plan *Plan) *PlanExec {
	return &PlanExec{
		plan: plan,
	}
}
