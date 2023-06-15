API

Create Table
POST /table

Create Player
POST /player

Player joins Table
PUT /table/<table ID>/players/<player ID>

Start Hand
POST /table/<table ID>/hand

Get Hand
GET /hand/<hand ID>

Player Folds/Checks/Calls/Raises
PUT /hand/<hand ID>/player/<player ID>/moves/<fold/check/call/raise>
