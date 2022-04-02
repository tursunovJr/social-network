package repository

import (
	"database/sql"

	"app/internal/model"

	"github.com/lib/pq"
)

type postRepoImpl struct {
	db *sql.DB
}

func NewPostRepo(db *sql.DB) PostRepo {
	return &postRepoImpl{db}
}

func scanPost(row MultiScanner, p *model.Post) error {
	var arr pq.Int64Array
	err := row.Scan(
		&p.Id,
		&p.UserId,
		&p.Created,
		&p.Tags,
		&p.Content,
		&p.AtchType,
		&p.AtchId,
		&p.AtchUrl,
		&arr,
		&p.CmtCount,
	)
	p.Reaction = arr
	return err
}

func (r postRepoImpl) Insert(p *model.Post) (id int, err error) {
	query := `insert into Post(userId, tags, content, atchType, atchId, atchUrl)
		values ($1, $2, $3, $4, $5, $6) returning id`
	row := r.db.QueryRow(query, p.UserId, p.Tags, p.Content, p.AtchType, p.AtchId, p.AtchUrl)
	err = row.Scan(&id)
	return
}

func (r postRepoImpl) Select(postId int) (post model.Post, err error) {
	row := r.db.QueryRow("select * from Post where id=$1 limit 1", postId)
	err = scanPost(row, &post)
	return
}

func (r postRepoImpl) SelectByUserId(userId int) (posts []int64, err error) {
	row := r.db.QueryRow("select array(select id from Post where userId=$1 order by created desc)", userId)

	var arr pq.Int64Array
	err = row.Scan(&arr)
	posts = arr

	return
}
