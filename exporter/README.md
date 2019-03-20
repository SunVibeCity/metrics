# Client for Growatt API

Credit for [Sjoerd Langkemper](https://github.com/Sjord/growatt_api_client) - who reverse engineered the API from the mobile app.

Unfortunately the [Official API Documentation](https://raw.githubusercontent.com/SunVibeCity/metrics/master/exporter/Growatt-Server-Open-API-protocol-standards.pdf) does not complete, neither works.

This exporter is logging into Growatt's server, fetches monitoring data and print it Prometheus metrics format.


## Simple API Documentation

### Login
Request curl
```text
curl -v -X POST \
  https://server.growatt.com/LoginAPI.do \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/x-www-form-urlencoded' \
  -d 'userName=JohnDoe&password=1234567890abcdef1234567890abcdef'
```
Request http
```text
POST /LoginAPI.do HTTP/1.1
Host: server.growatt.com
Content-Type: application/x-www-form-urlencoded
Cache-Control: no-cache

userName=JohnDoe&password=1234567890abcdef1234567890abcdef
```
Response
```text
Server : [Tengine]
Date : [Wed, 20 Mar 2019 04:18:12 GMT]
Content-Type : [application/json;charset=UTF-8]
Set-Cookie : [JSESSIONID=...; Path=/; HttpOnly SERVERID=...;Path=/]

{
  "back":{
    "userId":123456,
    "userLevel":1,
    "success":true,
  }
}
```

### PlantList

Request curl
```text
curl -v -X GET \
  https://server.growatt.com/PlantListAPI.do \
  -H 'cache-control: no-cache' \
  -H 'content-type: application/json' \
  --cookie "JSESSIONID=...; SERVERID=..."
```
Request http
```text
GET /PlantListAPI.do HTTP/1.1
Host: server.growatt.com
Content-Type: application/json
Cache-Control: no-cache
Cookie: JSESSIONID=...; SERVERID=...
```

Response
```text
Server : [Tengine]
Date : [Wed, 20 Mar 2019 07:40:48 GMT]
Content-Type : [application/json;charset=UTF-8]
Set-Cookie : [SERVERID=...;Path=/]

{
   "back":{
      "data":[
         {
            "plantMoneyText":"59.1 ",
            "plantName":"Fancy Plant",
            "plantId":"123456",
            "isHaveStorage":"false",
            "todayEnergy":"4.1 kWh",
            "totalEnergy":"632.5 kWh",
            "currentPower":"963.2 W"
         }
      ],
      "totalData":{
         "currentPowerSum":"963.2 W",
         "CO2Sum":"632.5 T",
         "isHaveStorage":"false",
         "eTotalMoneyText":"59.1 ",
         "todayEnergySum":"4.1 kWh",
         "totalEnergySum":"632.5 kWh"
      },
      "success":true
   }
}
```

### PlantDetail

Request http
```text
GET /PlantDetailAPI.do?plantId=123456&amp;type=1&amp;date=2019-03-20 HTTP/1.1
Host: server.growatt.com
Cache-Control: no-cache
Cookie: JSESSIONID=...; SERVERID=...
```

Response
```text
Server : [Tengine]
Date : [Wed, 20 Mar 2019 07:40:48 GMT]
Content-Type : [application/json;charset=UTF-8]
Set-Cookie : [SERVERID=...;Path=/]

{
   "back":{
      "plantData":{
         "plantMoneyText":"0.8 ",
         "plantName":"Fancy Plant",
         "plantId":"123456",
         "currentEnergy":"8.6 kWh"
      },
      "data":{
         "05:00":"0",
         "07:30":"471.82",
         "08:30":"961.1",
         "09:00":"830.73",
         "03:00":"0",
         "02:00":"0",
         "01:00":"0",
         "04:00":"0",
         "06:30":"73.4",
         "07:00":"246.22",
         "05:30":"0",
         "09:30":"1288.13",
         "08:00":"655.92",
         "03:30":"0",
         "01:30":"0",
         "04:30":"0",
         "06:00":"0",
         "00:30":"0",
         "02:30":"0"
      },
      "success":true
   }
}
```