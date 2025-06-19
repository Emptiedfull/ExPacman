from fastapi import FastAPI
from fastapi.staticfiles import StaticFiles
from fastapi.responses import HTMLResponse,FileResponse

app = FastAPI()

app.mount("/static", StaticFiles(directory="static"), name="static")

@app.get("/")
async def read_root():
    return HTMLResponse(content=open("static/index.html").read())

@app.get("/lobby/{wildcard_path:path}")
async def read_lobby(wildcard_path: str):
    return HTMLResponse(content=open("static/lobby.html").read())
