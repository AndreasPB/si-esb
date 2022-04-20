import json
from fastapi import FastAPI, HTTPException


app = FastAPI()


messages = {
    "1": [
        {"id": "1", "message": "m1", "access": "*"},
        {"id": "2", "message": "m2", "access": "*"},
        {"id": "3", "message": "m3", "access": "*"},
        {"id": "4", "message": "m4", "access": "*"},
    ]
}


@app.get("/")
async def index():
    return "(◡ ‿ ◡ ✿)"


@app.get("/provider/{id}/from/{last_message_id}/limit/{limit}/token/{token}")
async def _(id, last_message_id, limit: int, token):
    try:
        return json.dumps(messages[id][:limit])
    except Exception as e:
        return HTTPException(status_code=204)
