package repository

import (
	"database/sql"
	"fmt"
	"governance-indexer/internal/models"
	"strings"
	"time"
)

type ProposalPostgres struct {
	db *sql.DB
}

func NewProposalPostgres(db *sql.DB) *ProposalPostgres {
	return &ProposalPostgres{db: db}
}

func (p ProposalPostgres) AddProposal(proposals []models.Proposals) error {

	if len(proposals) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(proposals))
	args := make([]interface{}, 0, len(proposals)*4)

	i := 1
	for _, t := range proposals {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", i, i+1, i+2, i+3))
		args = append(args, t.ID, time.Unix(t.Created, 0), t.State, t.Author)
		i += 4
	}

	query := `
        INSERT INTO proposal(hex_id, created_at, state, author)
        VALUES ` + strings.Join(placeholders, ", ") + `
        ON CONFLICT (hex_id) DO NOTHING
    `

	//fmt.Println(query)
	//fmt.Println(args...)

	_, err := p.db.Exec(query, args...)
	fmt.Println(err)
	if err != nil {
		return err
	}

	for _, p := range proposals {
		createdTime := time.Unix(p.Created, 0).Format("2006-01-02 15:04:05 UTC")

		fmt.Printf(
			"- [%s] %s\n  Space: %s (%s)\n  Author: %s\n  State: %s\n\n",
			createdTime,
			p.Title,
			p.Space.Name,
			p.Space.ID,
			p.Author,
			p.State,
		)
	}
	println(proposals)
	return nil
}

//type ProposalsRepo struct {
//	Id      int
//	HexId   string
//	Created models.NullableTime
//	State   int
//	Space   int64
//	Author  string
//}

//func (p ProposalPostgres) GetDiff(proposals []models.Proposals) ([]models.Proposals, error) {
//
//	query := fmt.Sprintf(`
//		SELECT id, hex_id FROM proposals
//		FROM unnest(ARRAY[$1]) AS x(id)
//		LEFT JOIN users u ON u.id = x.id
//		WHERE u.id IS NULL;
//
//		`, usersTable)
//
//	err := p.db.QueryRow(query).Scan(&proposals)
//
//	return nil, nil
//}

//type SpaceRepo struct {
//	ID   string
//	Name string
//}

//func (p ProposalPostgres) GetHexId() ([]ProposalsHexId, error) {
//
//	var proposals []ProposalsHexId
//	query := "SELECT id, hex_id FROM proposals"
//	err := p.db.QueryRow(query).Scan(&proposals)
//	if err != nil {
//		return nil, err
//	}
//	return proposals, nil
//}
