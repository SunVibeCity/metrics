FROM prom/prometheus:v2.8.0
CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--storage.tsdb.retention.time=15d", \
             "--web.console.libraries=/etc/prometheus/console_libraries", \
             "--web.console.templates=/etc/prometheus/consoles" ]
