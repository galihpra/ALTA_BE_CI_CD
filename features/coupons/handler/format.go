package handler

type AddRequest struct {
	NamaProgram string `json:"nama_program" form:"nama_program" validate:"required"`
	Kode        string `json:"kode" form:"kode" validate:"required"`
	Link        string `json:"link" form:"link" validate:"url"`
	Gambar      string `json:"gambar" form:"gambar"`
}

type AddResponse struct {
	ID          uint   `json:"id"`
	NamaProgram string `json:"nama_program"`
	Kode        string `json:"kode"`
	Link        string `json:"link"`
	Gambar      string `json:"gambar"`
}

type GetResponse struct {
	ID          uint   `json:"id"`
	NamaProgram string `json:"nama_program"`
	Kode        string `json:"kode"`
	Link        string `json:"link"`
	Gambar      string `json:"gambar"`
	UserID      uint   `json:"user_id"`
}
