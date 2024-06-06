package memStorage

const (
	MyTypeGauge   string = "gauge"
	MyTypeCounter string = "counter"
)

type Metric struct {
	Mtype  string
	MName  string
	MValue *float64
}
