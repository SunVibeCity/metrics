# Prometheus

Sample Prometheus to scrape data and fire alerting

# Getting started
## To update prometheus configuration
Edit /config/prometheus.yml
run
```
bash reload.sh local
```

## To add or modify alerting/recording rule
Create or modify a rule file in ./config following https://prometheus.io/docs/querying/rules/

check rules by running
```
promtool check-rules <rule-file-name>
```

reload config
```
bash reload.sh local
```
