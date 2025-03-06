package telegraf

import "log"

type Output interface {
	PluginDescriber

	// Connect to the Output; connect is only called once when the plugin starts
	Connect() error
	// Close any connections to the Output. Close is called once when the output
	// is shutting down. Close will not be called until all writes have finished,
	// and Write() will not be called once Close() has been, so locking is not
	// necessary.
	Close() error
	// Write takes in group of points to be written to the Output
	Write(metrics []Metric) error
}

// AggregatingOutput adds aggregating functionality to an Output.  May be used
// if the Output only accepts a fixed set of aggregations over a time period.
// These functions may be called concurrently to the Write function.
type AggregatingOutput interface {
	Output

	// Add the metric to the aggregator
	Add(in Metric)
	// Push returns the aggregated metrics and is called every flush interval.
	Push() []Metric
	// Reset signals that the aggregator period is completed.
	Reset()
}

// Enhanced Write function to track errors
func (o *Output) Write(metrics []Metric) error {
	err := o.Write(metrics)
	if err != nil {
		log.Println("Error writing metrics:", err)
		// Increment error counter
		metrics[0].AddField("error_count", 1)
	}
	return err
}
