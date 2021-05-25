package chrono

import (
	"sync"
	"time"
)

type ScheduledExecutor interface {
	Schedule(task Task, delay time.Duration) ScheduledTask
	ScheduleWithFixedDelay(task Task, initialDelay time.Duration, delay time.Duration) ScheduledTask
	ScheduleAtRate(task Task, initialDelay time.Duration, period time.Duration) ScheduledTask
	Shutdown()
}

type ScheduledTaskExecutor struct {
	nextSequence          int
	nextSequenceMu        sync.RWMutex
	timer                 *time.Timer
	taskQueue             ScheduledTaskQueue
	newTaskChannel        chan *ScheduledRunnableTask
	removeTaskChannel     chan *ScheduledRunnableTask
	rescheduleTaskChannel chan *ScheduledRunnableTask
	taskWaitGroup         sync.WaitGroup
	taskRunner            TaskRunner
}

func NewDefaultScheduledExecutor() ScheduledExecutor {
	return NewScheduledTaskExecutor(NewDefaultTaskRunner())
}

func NewScheduledTaskExecutor(runner TaskRunner) *ScheduledTaskExecutor {
	if runner == nil {
		runner = NewDefaultTaskRunner()
	}

	executor := &ScheduledTaskExecutor{
		timer:                 time.NewTimer(1 * time.Hour),
		taskQueue:             make(ScheduledTaskQueue, 0),
		newTaskChannel:        make(chan *ScheduledRunnableTask),
		rescheduleTaskChannel: make(chan *ScheduledRunnableTask),
		removeTaskChannel:     make(chan *ScheduledRunnableTask),
		taskRunner:            runner,
	}

	executor.timer.Stop()

	go executor.run()

	return executor
}

func (executor *ScheduledTaskExecutor) Schedule(task Task, delay time.Duration) ScheduledTask {
	if task == nil {
		panic("task cannot be nil")
	}

	executor.nextSequenceMu.Lock()
	executor.nextSequence++
	scheduledTask := NewScheduledRunnableTask(executor.nextSequence, task, executor.calculateTriggerTime(delay), 0, false)
	executor.nextSequenceMu.Unlock()

	executor.addNewTask(scheduledTask)

	return scheduledTask
}

func (executor *ScheduledTaskExecutor) ScheduleWithFixedDelay(task Task, initialDelay time.Duration, delay time.Duration) ScheduledTask {
	if task == nil {
		panic("task cannot be nil")
	}

	executor.nextSequenceMu.Lock()
	executor.nextSequence++
	scheduledTask := NewScheduledRunnableTask(executor.nextSequence, task, executor.calculateTriggerTime(initialDelay), delay, false)
	executor.nextSequenceMu.Unlock()

	executor.addNewTask(scheduledTask)

	return scheduledTask
}

func (executor *ScheduledTaskExecutor) ScheduleAtRate(task Task, initialDelay time.Duration, period time.Duration) ScheduledTask {
	if task == nil {
		panic("task cannot be nil")
	}

	executor.nextSequenceMu.Lock()
	executor.nextSequence++
	scheduledTask := NewScheduledRunnableTask(executor.nextSequence, task, executor.calculateTriggerTime(initialDelay), period, true)
	executor.nextSequenceMu.Unlock()

	executor.addNewTask(scheduledTask)

	return scheduledTask
}

func (executor *ScheduledTaskExecutor) Shutdown() {

}

func (executor *ScheduledTaskExecutor) calculateTriggerTime(delay time.Duration) time.Time {
	if delay < 0 {
		delay = 0
	}

	return time.Now().Add(delay)
}

func (executor *ScheduledTaskExecutor) addNewTask(task *ScheduledRunnableTask) {
	executor.newTaskChannel <- task
}

func (executor *ScheduledTaskExecutor) run() {

	for {
		executor.taskQueue.SorByTriggerTime()

		if len(executor.taskQueue) == 0 {
			executor.timer.Stop()
		} else {
			executor.timer.Reset(executor.taskQueue[0].getDelay())
		}

		for {
			select {
			case clock := <-executor.timer.C:
				executor.timer.Stop()

				taskIndex := -1
				for index, scheduledTask := range executor.taskQueue {

					if scheduledTask.triggerTime.After(clock) || scheduledTask.triggerTime.IsZero() {
						break
					}

					taskIndex = index

					if scheduledTask.IsCancelled() {
						continue
					}

					if scheduledTask.isPeriodic() && scheduledTask.isFixedRate() {
						scheduledTask.triggerTime = scheduledTask.triggerTime.Add(scheduledTask.period)
					}

					executor.startTask(scheduledTask)
				}

				executor.taskQueue = executor.taskQueue[taskIndex+1:]
			case newScheduledTask := <-executor.newTaskChannel:
				executor.timer.Stop()
				executor.taskQueue = append(executor.taskQueue, newScheduledTask)
			case rescheduledTask := <-executor.rescheduleTaskChannel:
				executor.timer.Stop()
				executor.taskQueue = append(executor.taskQueue, rescheduledTask)
			case task := <-executor.removeTaskChannel:
				executor.timer.Stop()

				taskIndex := -1
				for index, scheduledTask := range executor.taskQueue {
					if scheduledTask.id == task.id {
						taskIndex = index
						break
					}
				}

				executor.taskQueue = append(executor.taskQueue[:taskIndex], executor.taskQueue[taskIndex+1:]...)
			}

			break
		}

	}

}

func (executor *ScheduledTaskExecutor) startTask(scheduledRunnableTask *ScheduledRunnableTask) {
	executor.taskWaitGroup.Add(1)

	go func() {
		defer func() {
			executor.taskWaitGroup.Done()
			scheduledRunnableTask.triggerTime = executor.calculateTriggerTime(scheduledRunnableTask.period)

			if !scheduledRunnableTask.isPeriodic() {
				scheduledRunnableTask.Cancel()
			} else {
				if !scheduledRunnableTask.isFixedRate() {
					executor.rescheduleTaskChannel <- scheduledRunnableTask
				}
			}
		}()

		if scheduledRunnableTask.isPeriodic() && scheduledRunnableTask.isFixedRate() {
			executor.rescheduleTaskChannel <- scheduledRunnableTask
		}

		executor.taskRunner.Run(scheduledRunnableTask.task)
	}()
}
