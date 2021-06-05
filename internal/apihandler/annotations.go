package apihandler

import (
	grafana_json "github.com/clambin/grafana-json"
	"github.com/clambin/sciensano/internal/vaccines"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func (handler *Handler) Annotations(name, query string, args *grafana_json.AnnotationRequestArgs) (annotations []grafana_json.Annotation, err error) {
	log.WithFields(log.Fields{
		"name":    name,
		"query":   query,
		"endTime": args.Range.To,
	}).Info("annotations")

	var batches []vaccines.Batch
	if batches, err = handler.Vaccines.GetBatches(); err == nil {
		for _, batch := range batches {
			annotations = append(annotations, grafana_json.Annotation{
				Time: time.Time(batch.Date),
				// Title: batch.Manufacturer,
				Text: "Amount: " + strconv.Itoa(batch.Amount),
			})
		}
	}
	return
}
