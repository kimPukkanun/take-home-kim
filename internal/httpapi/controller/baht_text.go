package controller

import (
	"encoding/json"
	"net/http"

	"github.com/shopspring/decimal"
)

type ThaiBahtTextConverter func(amount decimal.Decimal) (string, error)

type BahtTextController struct {
	convert ThaiBahtTextConverter
}

func NewBahtTextController(convert ThaiBahtTextConverter) *BahtTextController {
	return &BahtTextController{convert: convert}
}

type convertRequest struct {
	Amount string `json:"amount"`
}

type convertResponse struct {
	Amount        string `json:"amount"`
	RoundedAmount string `json:"roundedAmount"`
	Baht          string `json:"baht"`
	Satang        string `json:"satang"`
	Text          string `json:"text"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (c *BahtTextController) HandleConvert(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req convertRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}
	if req.Amount == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "amount is required"})
		return
	}

	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "amount must be a valid decimal string"})
		return
	}

	rounded := amount.Round(2)
	baht := rounded.Floor()
	satang := rounded.Sub(baht).Mul(decimal.NewFromInt(100)).Round(0)

	text, err := c.convert(amount)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, errorResponse{Error: "conversion failed"})
		return
	}

	resp := convertResponse{
		Amount:        amount.String(),
		RoundedAmount: rounded.StringFixed(2),
		Baht:          baht.String(),
		Satang:        satang.String(),
		Text:          text,
	}
	writeJSON(w, http.StatusOK, resp)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
