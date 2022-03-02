package postgres

import (
	"database/sql"
	"time"
	pb "two_services/task-services/genproto/task"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type taskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *taskRepo {
	return &taskRepo{db: db}
}

var (
	isupdated sql.NullString
	isnull    sql.NullString
)

func (r *taskRepo) Create(task *pb.CreateTaskReq) (*pb.TaskRes, error) {
	newTaskRes := pb.TaskRes{}
	query := `INSERT INTO task (
		id,
		assignee_id,
		title,
		deadline,
		status,
		created_at
		) VALUES ($1,$2,$3,$4,$5,$6) 
		RETURNING 
			id,
			assignee_id,
			title,
			deadline,
			status, 
			created_at`
	uid := uuid.New().String()
	err := r.db.QueryRow(query, uid, task.AssigneeId, task.Title, task.Deadline, task.Status, time.Now().Format(time.RFC3339)).Scan(
		&newTaskRes.Id,
		&newTaskRes.AssigneeId,
		&newTaskRes.Title,
		&newTaskRes.Deadline,
		&newTaskRes.Status,
		&newTaskRes.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newTaskRes, nil
}

func (r *taskRepo) Get(task *pb.GetAndDeleteTask) (*pb.TaskRes, error) {
	query := `SELECT 
		id,
		assignee_id,
		title,
		deadline,
		status,
		created_at,
		updated_at 
		FROM task WHERE 
			id=$1 AND deleted_at IS NULL`
	newTaskRes := pb.TaskRes{}
	err := r.db.QueryRow(query, task.Id).Scan(
		&newTaskRes.Id,
		&newTaskRes.AssigneeId,
		&newTaskRes.Title,
		&newTaskRes.Deadline,
		&newTaskRes.Status,
		&newTaskRes.CreatedAt,
		&isnull,
	)
	if isnull.Valid {
		newTaskRes.UpdatedAt = isnull.String
	}
	if err != nil {
		return nil, err
	}
	return &newTaskRes, nil
}

func (r *taskRepo) Update(task *pb.UpdateTaskReq) (*pb.TaskRes, error) {
	query := `UPDATE task SET 
		assignee_id = $1,
		title = $2,
		deadline = $3,
		status = $4,
		updated_at = $5
		WHERE id=$6 AND deleted_at IS NULL RETURNING id,assignee_id,title,deadline,status,created_at,updated_at`
	newTaskRes := pb.TaskRes{}
	err := r.db.QueryRow(query, task.AssigneeId, task.Title, task.Deadline, task.Status, time.Now().Format(time.RFC3339), task.Id).Scan(
		&newTaskRes.Id,
		&newTaskRes.AssigneeId,
		&newTaskRes.Title,
		&newTaskRes.Deadline,
		&newTaskRes.Status,
		&newTaskRes.CreatedAt,
		&newTaskRes.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &newTaskRes, nil
}

func (r *taskRepo) Delete(task *pb.GetAndDeleteTask) (*pb.ErrOrStatus, error) {
	query := `UPDATE task SET deleted_at=$1 WHERE id=$2`
	_, err := r.db.Exec(query, time.Now().Format(time.RFC3339), task.Id)
	if err != nil {
		return &pb.ErrOrStatus{
			Message: "Cant Deleted",
		}, err
	}
	return &pb.ErrOrStatus{
		Message: "Ok!",
	}, nil
}

func (r *taskRepo) List(task *pb.LimAndPage) (*pb.TasksList, error) {
	start := (task.Page - 1) * task.Limit
	query := `SELECT id,assignee_id,title,deadline,status,created_at,updated_at FROM task WHERE deleted_at IS NULL ORDER BY title OFFSET $1 LIMIT $2`
	newTasksList := pb.TasksList{}
	rows, err := r.db.Query(query, start, task.Limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		newTaskRes := pb.TaskRes{}
		err = rows.Scan(
			&newTaskRes.Id,
			&newTaskRes.AssigneeId,
			&newTaskRes.Title,
			&newTaskRes.Deadline,
			&newTaskRes.Status,
			&newTaskRes.CreatedAt,
			&isupdated,
		)
		if isupdated.Valid {
			newTaskRes.UpdatedAt = isupdated.String
		}
		if err != nil {
			return nil, err
		}
		newTasksList.Tasks = append(newTasksList.Tasks, &newTaskRes)
	}
	return &newTasksList, nil
}

func (r *taskRepo) ListOverdue(task *pb.Empty) (*pb.TasksList, error) {
	query := `SELECT id,assignee_id,title,deadline,status,created_at,updated_at FROM task WHERE deadline<$1 AND deleted_at IS NULL`
	newTasksList := pb.TasksList{}
	row, err := r.db.Query(query, time.Now().Format(time.RFC3339))
	if err != nil {
		return nil, err
	}

	for row.Next() {
		newTaskRes := pb.TaskRes{}

		err = row.Scan(
			&newTaskRes.Id,
			&newTaskRes.AssigneeId,
			&newTaskRes.Title,
			&newTaskRes.Deadline,
			&newTaskRes.Status,
			&newTaskRes.CreatedAt,
			&isupdated,
		)
		if isupdated.Valid {
			newTaskRes.UpdatedAt = isupdated.String
		}
		if err != nil {
			return nil, err
		}
		newTasksList.Tasks = append(newTasksList.Tasks, &newTaskRes)
	}
	return &newTasksList, nil
}
