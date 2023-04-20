package repository

type RepositoryFizz struct {
	*Repository
}

func NewFizzRepository(repo *Repository) *RepositoryFizz {
	return &RepositoryFizz{
		repo,
	}
}

// AuditRequest
// func (r *RepositoryFizz) AuditRequest(checksum string) {
// 	res, err := r.db.Exec("DELETE FROM auth_brand WHERE code = ?", code)
// 	if err != nil {
// 		return 0, nil
// 	}
// 	return res.RowsAffected()
// }
