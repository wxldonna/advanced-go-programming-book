package tasks

import (
	"sync"
	"time"

	"github.xiaoliang.graphql.users/graph/model"
)

type Tasks struct {
	sync.Mutex

	tasks  map[int]*model.Task
	nextId int
}

func New() *Tasks {
	ts := &Tasks{}
	ts.tasks = make(map[int]*model.Task)
	ts.nextId = 0
	return ts
}

func NewAttachments() *Attachments {
	as := &Attachments{}
	as.atts = make(map[int]*model.Attachment)
	as.nextId = 0
	return as
}

// CreateTask creates a new task in the store.
func (ts *Tasks) CreateTask(text string, tags []string, due time.Time, attachments []*model.Attachment) *model.Task {
	ts.Lock()
	defer ts.Unlock()

	task := &model.Task{
		ID:          ts.nextId,
		Text:        text,
		Due:         due,
		Attachments: attachments}
	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task
}

// GetAllTasks returns all the tasks in the store, in arbitrary order.
func (ts *Tasks) GetAllTasks() []*model.Task {
	ts.Lock()
	defer ts.Unlock()

	allTasks := make([]*model.Task, 0, len(ts.tasks))
	for _, task := range ts.tasks {
		allTasks = append(allTasks, task)
	}
	return allTasks
}

type Attachments struct {
	sync.Mutex

	atts   map[int]*model.Attachment
	nextId int
}

func (as *Attachments) GetAttchment(id int) []*model.Attachment {
	ats := make([]*model.Attachment, 0)
	ats = append(ats, as.atts[id])

	return ats
}

func (as *Attachments) Create(name string, date time.Time, content string) *model.Attachment {
	as.Lock()
	defer as.Unlock()

	attachment := &model.Attachment{
		ID:       as.nextId,
		Name:     name,
		Date:     date,
		Contents: content,
	}

	as.atts[as.nextId] = attachment
	as.nextId++
	return attachment
}

func (as *Attachments) GetAllAtts() []*model.Attachment {
	as.Lock()
	defer as.Unlock()

	allTasks := make([]*model.Attachment, 0, len(as.atts))
	for _, at := range as.atts {
		allTasks = append(allTasks, at)
	}
	return allTasks
}
