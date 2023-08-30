package repository

import (
	"avito/config"
	models "avito/internal"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Repository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func InitRepository(c *config.Config) *Repository {
	ctx := context.Background()
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		*c.Postgres.Host,
		*c.Postgres.Port,
		*c.Postgres.User,
		*c.Postgres.Password,
		*c.Postgres.DbName)
	fmt.Println(connectionUrl)
	pool, err := pgxpool.New(ctx, connectionUrl)

	if err != nil {

		log.Panic("couldn't connect database service:" + err.Error())
	} else {
		log.Println("database is connected")
	}
	rep := Repository{pool: pool, ctx: ctx}
	err = rep.TableCreation()
	if err != nil {
		log.Panic("data objects cannot initialized:" + err.Error())
	}
	return &rep
}

func (rep *Repository) TableCreation() error {

	_, err := rep.pool.Exec(rep.ctx, `CREATE TABLE IF NOT EXISTS segments
		(
			id   SERIAL PRIMARY KEY ,
			segmentName TEXT NOT NULL UNIQUE,
			percents INT NOT NULL
		);`)
	if err != nil {
		return err
	}

	_, err = rep.pool.Exec(rep.ctx, `CREATE TABLE IF NOT EXISTS subscription
		(
			userId    BIGSERIAL,
			segment   SERIAL REFERENCES segments(id),
			PRIMARY KEY (userId, segment),
			timeToDie date NOT NULL
		);`)
	if err != nil {
		return err
	}
	return nil
}

func (rep *Repository) GetSubs(id int64) ([]models.UserSubscription, error) {
	var result []models.UserSubscription
	var tmp models.UserSubscription
	rows, err := rep.pool.Query(rep.ctx, "SELECT segment, timeToDie FROM subscription WHERE userId=$1", id)
	if err != nil {
		return nil, errors.New("subscription not found:" + err.Error())
	}
	for rows.Next() {
		if err := rows.Scan(&tmp.Name, &tmp.TimeOut); err != nil {
			return []models.UserSubscription{}, err
		}
		result = append(result, tmp)
	}
	if err != nil {
		return nil, errors.New("user are not exist:" + err.Error())
	}
	return result, nil
}

func (rep *Repository) CreateSegment(body models.SegmentBody) (models.SegmentBody, error) {
	var result models.SegmentBody
	_, err := rep.pool.Exec(rep.ctx, "INSERT INTO segments (segmentName, percents) VALUES ($1, $2)", body.Name, body.Percent)
	if err != nil {
		return models.SegmentBody{}, errors.New("cannot create segment:" + err.Error())
	}
	err = rep.pool.QueryRow(rep.ctx, "SELECT * FROM segments WHERE segmentName = $1", body.Name).Scan(&result.Id, &result.Name, &result.Percent)
	if err != nil {
		return models.SegmentBody{}, err
	}
	if body.Percent != 0 {
		rows, err := rep.pool.Query(rep.ctx, "SELECT userId FROM subscription ORDER BY random() LIMIT (SELECT count(userId)*$1/100 FROM subscription)", body.Percent)
		if err != nil {
			return models.SegmentBody{}, errors.New("cannot make percent subs:" + err.Error())
		}
		var ids []int64
		var tmp int64
		for rows.Next() {
			if err := rows.Scan(&tmp); err != nil {
				return models.SegmentBody{}, err
			}
			ids = append(ids, tmp)
		}
		fmt.Println(ids)
		for i := 0; i < len(ids); i++ {
			_, err := rep.pool.Exec(rep.ctx, "INSERT INTO subscription (userId, segment, timeToDie) VALUES ($1, $2, $3)", ids[i], result.Id, time.Time{})
			if err != nil {
				return models.SegmentBody{}, errors.New("cannot create note about new subscription:" + err.Error())
			}
		}
	}
	return result, nil
}

func (rep *Repository) DeleteSegment(segmentName string) error {
	_, err := rep.pool.Exec(rep.ctx, "DELETE FROM subscription WHERE segment = ("+
		"SELECT id FROM segments WHERE segmentName = $1)", segmentName)
	if err != nil {
		return errors.New("cannot delete user subscription:" + err.Error())
	}
	_, err = rep.pool.Exec(rep.ctx, "DELETE FROM segments WHERE segmentName = $1", segmentName)
	if err != nil {
		return errors.New("cannot delete segment:" + err.Error())
	}
	return nil
}

func (rep *Repository) GetSegments() ([]models.SegmentBody, error) {
	var result []models.SegmentBody
	var tmp models.SegmentBody
	rows, err := rep.pool.Query(rep.ctx, "SELECT * FROM segments")
	for rows.Next() {
		if err := rows.Scan(&tmp.Id, &tmp.Name, &tmp.Percent); err != nil {
			return []models.SegmentBody{}, err
		}
		result = append(result, tmp)
	}
	if err != nil {
		return nil, errors.New("cannot get info about segments:" + err.Error())
	}
	return result, nil
}

func (rep *Repository) Subscriber(subscriber models.Subscriber) error {
	for i := 0; i < len(subscriber.Add); i++ {
		_, err := rep.pool.Exec(rep.ctx, "INSERT INTO subscription (userId, segment, timeToDie) VALUES($2, (SELECT id FROM segments WHERE segmentName=$1), $3)", subscriber.Add[i], subscriber.UserId, time.Time{})
		if err != nil {
			return errors.New(fmt.Sprintf("cannot add the new subscription for user %d\n", subscriber.UserId) + err.Error())
		}
	}
	for i := 0; i < len(subscriber.Delete); i++ {
		_, err := rep.pool.Exec(rep.ctx, "DELETE FROM subscription WHERE userId = $1 AND segment = (SELECT id FROM segments WHERE segmentName = $2)", subscriber.UserId, subscriber.Delete[i])
		if err != nil {
			return errors.New(fmt.Sprintf("cannot delete the subscription for user %d\n", subscriber.UserId) + err.Error())
		}
	}
	return nil
}

func (rep *Repository) SubWIthTimeout(id int64, name string, timeToDie time.Time) error {
	_, err := rep.pool.Exec(rep.ctx, "INSERT INTO subscription (userId, segment, timeToDie) VALUES($2, (SELECT id FROM segments WHERE segmentName=$1), $3) ON CONFLICT (userId, segment) DO UPDATE SET timeToDie = $3", name, id, timeToDie)
	if err != nil {
		return errors.New("cannot create timeouts subscription: " + err.Error())
	}
	return nil
}

func (rep *Repository) TimeOutDeleter(sec chan int) {
	for {
		select {
		case <-sec:
			return
		default:
			time.Sleep(time.Second * 20)
			rows, err := rep.pool.Query(rep.ctx, "SELECT userId, segment FROM subscription WHERE timeToDie<$1 AND timeToDie != $2", time.Now(), time.Time{})
			if err != nil {
				fmt.Println(err)
			}
			var (
				user    int
				segment int
			)
			for rows.Next() {
				err := rows.Scan(&user, &segment)
				if err != nil {
					fmt.Println(err.Error())
				}
				_, err = rep.pool.Exec(context.Background(), "DELETE FROM subscription WHERE userId=$1 AND segment=$2", user, segment)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}

}
