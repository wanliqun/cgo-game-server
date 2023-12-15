package stress

type Task interface {
	Do(robot *Rotbot, iteration int) error
}

type SimpleTask struct{}

func (t *SimpleTask) Do(robot *Rotbot, iteration int) error {
	return robot.client.Info()
}
