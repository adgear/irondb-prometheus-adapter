package handlers

import (
	"github.com/labstack/echo"
	dto "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

// PrometheusWrite2_0 - the application handler which converts a prometheus
// MetricFamily message to a MetricList for ingestion into IRONdb
func PrometheusWrite2_0(ctx echo.Context) error {
	var (
		// create our prometheus format decoder
		dec          = expfmt.NewDecoder(ctx.Request().Body, expfmt.Format(ctx.Request().Header.Get("Content-Type")))
		metricFamily = new(dto.MetricFamily)
		err          error
		data         []byte
	)
	// close request body
	defer ctx.Request().Body.Close()

	// decode the metrics into the metric family
	err = dec.Decode(metricFamily)
	if err != nil {
		ctx.Logger().Errorf("failed to decode: %s", err.Error())
		return err
	}
	ctx.Logger().Debugf("parsed metric-family: %+v\n", metricFamily)

	data, err = MakeMetricList(metricFamily, ctx.Param("account"),
		ctx.Param("check_name"), ctx.Param("check_uuid"))
	if err != nil {
		ctx.Logger().Errorf("failed to convert to flatbuffer: %s", err.Error())
		return err
	}

	// call snowth with flatbuffer data
	ctx.Logger().Debugf("converted flatbuffer: %+v\n", data)

	return nil
}

// PrometheusRead2_0 - the application handler which converts a prometheus
// read message to an IRONdb read message, and returns the results converted
// back into prometheus output
func PrometheusRead2_0(ctx echo.Context) error {
	return nil
}
