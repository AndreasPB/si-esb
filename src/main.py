import json
from fastapi import FastAPI, HTTPException


app = FastAPI()


messages = {"1": [{"id": "10", "message": "Hi", "access": "*"}]}


@app.get("/")
async def index():
    return "(◡ ‿ ◡ ✿)"

@app.get("/provider/{id}/from/{last_message_id}/limit/{limit}/token/{token}")
async def idk(id, last_message_id, limit, token):
    try:
        return json.dumps(messages[id])
    except Exception as e:
        return HTTPException(status_code=204)
