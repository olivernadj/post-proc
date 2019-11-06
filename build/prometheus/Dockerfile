FROM prom/prometheus:v2.8.0
MAINTAINER Oliver Nadj <mr.oliver.nadj@gmail.com>

ADD config/* /etc/prometheus/

CMD        [ "--config.file=/etc/prometheus/prometheus.yml", \
             "--storage.tsdb.path=/prometheus", \
             "--storage.tsdb.retention.time=1780d", \
             "--web.console.libraries=/etc/prometheus/console_libraries", \
             "--web.console.templates=/etc/prometheus/consoles" ]
