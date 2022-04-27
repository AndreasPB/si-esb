from pydantic import BaseModel


class Message(BaseModel):
    name: str
    content: str
