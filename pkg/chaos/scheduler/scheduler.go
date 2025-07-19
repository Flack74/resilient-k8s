package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/flack/chaos-engineering-as-a-platform/pkg/chaos/executor"
	"github.com/flack/chaos-engineering-as-a-platform/pkg/storage"
)

// ScheduleType defines the type of schedule
type ScheduleType string

const (
	ScheduleOneTime ScheduleType = "one-time"
	ScheduleCron    ScheduleType = "cron"
)

// Schedule represents a schedule for a chaos experiment
type Schedule struct {
	ID           string       `json:"id"`
	ExperimentID string       `json:"experiment_id"`
	Type         ScheduleType `json:"type"`
	CronExpression string     `json:"cron_expression,omitempty"`
	ExecuteAt    time.Time    `json:"execute_at,omitempty"`
	Enabled      bool         `json:"enabled"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// Scheduler schedules chaos experiments
type Scheduler struct {
	db       *storage.Database
	executor *executor.Executor
	schedules map[string]*Schedule
	mutex    sync.RWMutex
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

// NewScheduler creates a new experiment scheduler
func NewScheduler(db *storage.Database, executor *executor.Executor) *Scheduler {
	return &Scheduler{
		db:       db,
		executor: executor,
		schedules: make(map[string]*Schedule),
		stopCh:   make(chan struct{}),
	}
}

// Start starts the scheduler
func (s *Scheduler) Start() error {
	log.Println("Starting experiment scheduler...")

	// Start the scheduler loop
	s.wg.Add(1)
	go s.schedulerLoop()

	return nil
}

// Stop stops the scheduler
func (s *Scheduler) Stop() error {
	log.Println("Stopping experiment scheduler...")

	// Signal the scheduler loop to stop
	close(s.stopCh)

	// Wait for the scheduler loop to finish
	s.wg.Wait()

	return nil
}

// schedulerLoop is the main loop of the scheduler
func (s *Scheduler) schedulerLoop() {
	defer s.wg.Done()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.stopCh:
			log.Println("Scheduler loop stopping...")
			return
		case <-ticker.C:
			s.checkSchedules()
		}
	}
}

// checkSchedules checks for scheduled experiments to run
func (s *Scheduler) checkSchedules() {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	now := time.Now()
	for id, schedule := range s.schedules {
		if !schedule.Enabled {
			continue
		}

		if schedule.Type == ScheduleOneTime && !schedule.ExecuteAt.IsZero() && now.After(schedule.ExecuteAt) {
			// Execute one-time schedule
			go s.executeSchedule(id)
		} else if schedule.Type == ScheduleCron {
			// In a real implementation, this would check if the cron expression matches the current time
			// For simplicity, we'll skip this for now
		}
	}
}

// executeSchedule executes a scheduled experiment
func (s *Scheduler) executeSchedule(scheduleID string) {
	s.mutex.RLock()
	schedule, exists := s.schedules[scheduleID]
	s.mutex.RUnlock()

	if !exists {
		log.Printf("Schedule %s not found", scheduleID)
		return
	}

	log.Printf("Executing scheduled experiment %s", schedule.ExperimentID)

	// Execute the experiment
	_, err := s.executor.ExecuteExperiment(schedule.ExperimentID)
	if err != nil {
		log.Printf("Failed to execute scheduled experiment %s: %v", schedule.ExperimentID, err)
	}

	// If it's a one-time schedule, disable it
	if schedule.Type == ScheduleOneTime {
		s.mutex.Lock()
		schedule.Enabled = false
		s.mutex.Unlock()
	}
}

// AddSchedule adds a new schedule
func (s *Scheduler) AddSchedule(schedule *Schedule) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.schedules[schedule.ID] = schedule
}

// RemoveSchedule removes a schedule
func (s *Scheduler) RemoveSchedule(scheduleID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.schedules, scheduleID)
}

// GetSchedule gets a schedule by ID
func (s *Scheduler) GetSchedule(scheduleID string) (*Schedule, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	schedule, exists := s.schedules[scheduleID]
	return schedule, exists
}

// ListSchedules lists all schedules
func (s *Scheduler) ListSchedules() []*Schedule {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	schedules := make([]*Schedule, 0, len(s.schedules))
	for _, schedule := range s.schedules {
		schedules = append(schedules, schedule)
	}

	return schedules
}