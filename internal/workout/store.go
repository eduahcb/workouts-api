package workout

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutStore interface {
	Create(workout *Workout) (*Workout, error)
	GetByID(id int64) (*Workout, error)
	Update(workout *Workout) error
	Delete(id int64) error
}

type PostgresWorkoutStore struct {
	db *pgxpool.Pool
}

func NewPostgresWorkoutStore(db *pgxpool.Pool) *PostgresWorkoutStore {
	return &PostgresWorkoutStore{db: db}
}

func (pg *PostgresWorkoutStore) Create(workout *Workout) (*Workout, error) {
	tx, err := pg.db.Begin(context.Background())

	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `
	  INSERT INTO workouts (title, description, duration_minutes, calories_burned)
	  VALUES ($1, $2, $3, $4)
	  RETURNING id
	`

	err = tx.QueryRow(context.Background(), query, workout.Title, workout.Description, workout.DurationMinutes, workout.CaloriesBurned).Scan(&workout.ID)

	if err != nil {
		return nil, err
	}

	for i := range workout.Entries {
		entry := &workout.Entries[i]

		query := `
		  INSERT INTO workouts_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		  RETURNING id
		`

		err = tx.QueryRow(context.Background(), query, workout.ID, entry.ExerciseName, entry.Sets, entry.Reps, entry.DurationSeconds, entry.Weight, entry.Notes, entry.OrderIndex).Scan(&entry.ID)

		if err != nil {
			return nil, err
		}
	}
	err = tx.Commit(context.Background())

	if err != nil {
		return nil, err
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) GetByID(id int64) (*Workout, error) {
	var workout = &Workout{}

	query := `
	  SELECT id, title, description, duration_minutes, calories_burned FROM workouts WHERE id = $1
	`

	err := pg.db.QueryRow(context.Background(), query, id).Scan(&workout.ID, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned)

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		println("entrou aqui")
		return nil, err
	}

	entryQuery := `
	  SELECT id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index FROM workouts_entries WHERE workout_id = $1 ORDER BY order_index
	`

	rows, err := pg.db.Query(context.Background(), entryQuery, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var workoutEntry WorkoutEntry

		err = rows.Scan(&workoutEntry.ID, &workoutEntry.ExerciseName, &workoutEntry.Sets, &workoutEntry.Reps, &workoutEntry.DurationSeconds, &workoutEntry.Weight, &workoutEntry.Notes, &workoutEntry.OrderIndex)

		if err != nil {
			return nil, err
		}

		workout.Entries = append(workout.Entries, workoutEntry)
	}

	return workout, nil
}

func (pg *PostgresWorkoutStore) Update(workout *Workout) error {
	tx, err := pg.db.Begin(context.Background())

	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `
	  UPDATE workouts
	  SET title = $1, description = $2, duration_minutes = $3, calories_burned = $4
	  WHERE id = $5
	`

	result, err := tx.Exec(context.Background(), query, &workout.Title, &workout.Description, &workout.DurationMinutes, &workout.CaloriesBurned, &workout.ID)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM workouts_entries WHERE workout_id = $1", workout.ID)

	if err != nil {
		return nil
	}

	for _, entry := range workout.Entries {
		query := `
		  INSERT INTO workouts_entries (workout_id, exercise_name, sets, reps, duration_seconds, weight, notes, order_index)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`

		_, err := tx.Exec(context.Background(), query, &workout.ID, &entry.ExerciseName, &entry.Sets, &entry.Reps, &entry.DurationSeconds, &entry.Weight, &entry.Notes, &entry.OrderIndex)

		if err != nil {
			return err
		}
	}

	return tx.Commit(context.Background())
}

func (pg *PostgresWorkoutStore) Delete(id int64) error {
	tx, err := pg.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	result, err := tx.Exec(context.Background(), "DELETE FROM workouts WHERE id = $1", id)

	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return pgx.ErrNoRows
	}

	return tx.Commit(context.Background())
}
