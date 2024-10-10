# POST /login
# Host: localhost:8080
# Content-Type: application/json
#
# {
#      "email": "test@test.com",
#      "password": "password123"
# }
# {% 
#
# local body = context.json_decode(context.result.body)
# context.set_env("TOKEN", body.token)
#
# %}

POST http://localhost:42069/api/v1/login

{
  "email": "andresc@incogres.com",
  "password": "fjcl1229"
}
