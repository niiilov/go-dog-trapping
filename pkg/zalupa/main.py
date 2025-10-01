from openpyxl import load_workbook





from pathlib import Path

from fastapi import FastAPI
from fastapi.responses import  JSONResponse
from pydantic import BaseModel

from datetime import datetime
import uvicorn
import os












app = FastAPI()


class RequestFull(BaseModel):
    number: str
    source_name: str
    applicant_name: str
    address: str
    dogs_count: int
    behavior: str
    urgency: str
    contact_person: str



@app.post("/generate")
async def generate_doc(req: RequestFull):
    filename = f"zayavka_{req.number}.xlsx"
    template = "tamplate.xlsx"
    # Передайте нужные параметры в функцию создания документа

    data = {
        "B15": f"Заявка № {req.number}",
        "B12": datetime.now().strftime('%d.%m.%Y'),
        "B20": req.source_name,       
        "D20": req.applicant_name,         
        "E20": req.address,               
        "F20": req.dogs_count,      
        "G20": req.behavior,    
        "H20": req.urgency,                 
        "I20": req.contact_person
    }
    generate_document(template, filename, data)
    print("Документ успешно создан:", filename)
    shared_dir = Path("/app/shared")
    if os.path.exists(f"{shared_dir}/zayavka_{req.number}.xlsx"):
        return JSONResponse({"status": "ok"}, status_code=200)
    return JSONResponse({"error": "File not created"}, status_code=500)






if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)




def generate_document(template_path, output_path, data):
    """
    Генерация Excel-документа на основе шаблона
    :param template_path: путь к шаблону (xlsx)
    :param output_path: путь для сохранения нового файла
    :param data: словарь с данными {ячейка: значение}
    """
    # Загружаем шаблон
    wb = load_workbook(template_path)
    ws = wb.active  # если один лист, берём его

    # Подставляем данные
    for cell, value in data.items():
        ws[cell] = value

    shared_dir = Path("/app/shared")
    # Сохраняем результат
    wb.save(f"{shared_dir}/{output_path}")


