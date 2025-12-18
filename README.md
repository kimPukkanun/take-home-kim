# Thai Baht Text (Golang)

Converts a `decimal.Decimal` amount (using `github.com/shopspring/decimal`) into Thai currency text.

Examples:

- `1234` → `หนึ่งพันสองร้อยสามสิบสี่บาทถ้วน`
- `33333.75` → `สามหมื่นสามพันสามร้อยสามสิบสามบาทเจ็ดสิบห้าสตางค์`

## Usage

- Conversion function: `thaibahttext.ToThaiBahtText(amount decimal.Decimal) (string, error)`
- HTTP API server: `cmd/api/main.go`

## HTTP API + Swagger

Start the API server:

```powershell
go run ./cmd/api
```

Endpoints:

- `POST http://localhost:8080/v1/baht-text`
- Swagger UI: `http://localhost:8080/docs`

Example request:

```powershell
curl -Method Post http://localhost:8080/v1/baht-text -ContentType application/json -Body '{"amount":"33333.75"}'
```

## Notes / assumptions

- Negative values are prefixed with `ลบ`.
- The input is rounded to 2 decimal places before rendering satang.
- Supports large numbers by grouping digits with repeated `ล้าน`.
- This project use mvc for easier integration.
- This is all I can create from the time I have left as I stated in the previous email. I have been working late for a couple of night in a rows. Thank you for your understanding.
