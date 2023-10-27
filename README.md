# RefineMetric Exporter

You can use metric configuration or writing an exporter by registering and then implement the Update method to introducing more metrics.

## Metric Configuration Guide
This file provides guidance on the metric configuration for collecting data from various data sources.

### Metric Format:
Metrics are defined in the following structured format:
```yaml
- name: <metric_name>
  description: <description>
  data_source_type: <data_source_type>
  query: <promql_query>
  expected_labels:
    - <label_name 1>
    ...
    - <label_name N>
  aggregation:
    labels:
      - <label_name 1>
      ...
      - <label_name M>
    method: <aggregation_method>
```

### Explanation of Fields:
- $<metric\_name>$: Name of the metric that will be used in the output file.

- $description$: just description

- $<data\_source\_type>$: Type of data source used to collect the metric. Available data source types include: $promql$

- $query$: query to fetch the data

- $<promql\_query>$: The promQL query utilized to fetch the metric.


- $expected\_labels$: List of labels expected to be present in the metric. The $<promql\_query>$ must return metrics with these specified labels.

- $aggregation$: Method used to aggregate metrics. Available aggregation methods are:
    $average$,
    $max$,
    $min$,
    $sum$,

### Aggregation Example

Consider the following scenario:

If the job label is in the expected_labels but not in the aggregation labels and the promQL query returns the metrics as:

    job="job1", instance="instance1", value=1
    job="job1", instance="instance2", value=2
    job="job2", instance="instance1", value=3
    job="job2", instance="instance2", value=4
First, filter metrics with labels in expected_labels but not in aggregation labels, resulting in:

    job="job1", instance="instance1", value=1
    job="job1", instance="instance2", value=2
Then, use the aggregation method to aggregate the metrics (assuming average is used):

    job="job1", value=1.5
The same process would be applied for job="job2, leading to:

    job="job2", value=3.5
If no labels exist in expected_labels but not in the aggregation labels, the aggregation method will apply to all metrics returned by the promQL query. This results in the metric name being aggregated to a single value without any labels. 

If the aggregation labels are not specified, no aggregation is done; only aliasing the promQL query to the metric name occurs.

