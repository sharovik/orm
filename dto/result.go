package dto

//ResultInterface the interface for the result of query execution
type ResultInterface interface {
	Items() []ModelInterface
	AddItem(ModelInterface)
	Error() error
	SetError(error)
	LastInsertID() int64
	SetLastInsertID(int64)
}

type BaseResult struct {
	rows     []ModelInterface
	InsertID int64
	Err      error
}

func (r *BaseResult) Items() []ModelInterface {
	return r.rows
}

func (r *BaseResult) AddItem(model ModelInterface) {
	r.rows = append(r.rows, model)
	return
}

func (r *BaseResult) Error() error {
	return r.Err
}

func (r *BaseResult) SetError(err error) {
	r.Err = err
}

func (r *BaseResult) LastInsertID() int64 {
	return r.InsertID
}

func (r *BaseResult) SetLastInsertID(id int64) {
	r.InsertID = id
}
