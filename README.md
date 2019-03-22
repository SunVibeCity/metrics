# SunVibe Energy Production Monitoring

Energy Production Monitoring for SunVibe.city

## Infrastructure Stack
- [Docker](https://www.docker.com/) provides a way to run applications securely isolated in a container, packaged with all its dependencies and libraries.
- [Docker Compose](https://docs.docker.com/compose/) is a tool for defining and running multi-container Docker applications.
- Exporter Application is written in Go.
  - [exporter](./exporter/README.md)
  - [jupyter notebook](./exporter/growatt.ipynb)
- [Prometheus](https://prometheus.io/) monitoring and alerting toolkit.
- [Alertmanager](https://prometheus.io/docs/alerting/alertmanager/) handles alerts sent by client applications such as the Prometheus server. It takes care of deduplicating, grouping, and routing them to the correct receiver integration such as email, PagerDuty, or OpsGenie.
- [Grafana](https://grafana.com) is for beautiful monitoring and metric analytics & dashboards for Graphite, InfluxDB & Prometheus & More

## Build and maintain the containers
Git clone:
```
~$git clone git@github.com:SunVibeCity/metrics.git
~$ cd metrics
```
To get the containers up and running, run this command:
```
cb-monitoring$ docker-compose up -d --build
```
To check running containers, run this command:
```
cb-monitoring$ watch "sudo docker ps --format='table{{.Image}}\t{{.Names}}\t{{.Status}}\t{{.Ports}}'"
```
For stop services, run this command:
```
cb-monitoring$ docker-compose stop
```

### Set up Grafana on 1st run
http://localhost:3000/
- Username: admin
- Password: 5ecret

#### 1. Add datasource
- On http://localhost:3000/datasources page add:
- Name: Prometheus
- Type: Prometheus
- Url: http://prometheus:9090
- Access: proxy

![Add datasource](https://raw.githubusercontent.com/SunVibeCity/metrics/master/add-data-source.png "Add datasource")

#### 2. Import growatt.json dashboard
- Main menu (top left Grafana logo) -> Dashboards -> Import
- Upload `grafana/growatt.json`

#### 3. Create guest user
- Main menu (top left Grafana logo) -> admin Main Org. -> Users
- Invite a user
- Activate Pending Invitations

![Activate Pending Invitations](https://raw.githubusercontent.com/SunVibeCity/metrics/master/pending-invitation.png "Activate Pending Invitations")

## Endpoints of services
- Exporter: fetches data from Growatt's API and returns metrics
  - http://localhost:5000/metrics
- Prometheus: scraps metrics and store it in time series
  - http://localhost:9090/graph
- Alertmanger
  - http://localhost:9093/#/alerts
- Grafana Access (admin/5ecret):
  - http://localhost:3000/dashboard/db/growatt?refresh=5m&orgId=1
