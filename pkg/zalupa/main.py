from docx import Document
from docx.shared import Pt
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.section import WD_ORIENT

from pathlib import Path

from fastapi import FastAPI
from fastapi.responses import FileResponse, JSONResponse
from pydantic import BaseModel
from typing import Optional
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
    filename = f"zayavka_{req.number}.docx"
    # Передайте нужные параметры в функцию создания документа
    create_document(
        territorial_office=req.source_name,
        applicant=req.applicant_name ,
        address=req.address,
        dogs_count=str(req.dogs_count),
        behavior=req.behavior,
        urgency=req.urgency,
        contact_person=req.contact_person,
        number= req.number
        
    )
    shared_dir = Path("/app/shared")
    if os.path.exists(f"{shared_dir}/zayavka_{req.number}.docx"):
        return JSONResponse({"status": "ok"}, status_code=200)
    return JSONResponse({"error": "File not created"}, status_code=500)






def create_document(number, territorial_office, applicant, address, dogs_count, behavior, urgency, contact_person):
        
    # Создаем документ
    doc = Document()

    section = doc.sections[0]
    section.orientation = WD_ORIENT.LANDSCAPE
    section.page_width, section.page_height = section.page_height, section.page_width


    # Верхняя часть (две колонки: слева Администрация, справа ИП)
    table_header = doc.add_table(rows=1, cols=2)
    table_header.autofit = True

    # Левая ячейка
    cell_left = table_header.cell(0, 0)
    p_left = cell_left.paragraphs[0]
    p_left.alignment = WD_ALIGN_PARAGRAPH.LEFT
    run_left = p_left.add_run(
        "АДМИНИСТРАЦИЯ\nЛИПЕЦКОГО\nМУНИЦИПАЛЬНОГО ОКРУГА\nЛИПЕЦКОЙ ОБЛАСТИ\n"
        "398037, г. Липецк, Воевод проезд 30\nтел./факс 34-19-19\nalmo@admr.lipetsk.ru"
    )
    run_left.bold = True
    run_left.font.size = Pt(10)

    # Правая ячейка
    cell_right = table_header.cell(0, 1)
    p_right = cell_right.paragraphs[0]
    p_right.alignment = WD_ALIGN_PARAGRAPH.RIGHT
    run_right = p_right.add_run("ИП Ульшин Г.Н.\nenter100@list.ru")
    run_right.font.size = Pt(10)

    doc.add_paragraph(f"\n{datetime.now().strftime('%d.%m.%Y')} №______", style='Normal')

    # Заголовок
    title = doc.add_paragraph(f"Заявка №{number}")
    title.alignment = WD_ALIGN_PARAGRAPH.CENTER
    title.runs[0].bold = True

    subtitle = doc.add_paragraph(
        "на оказание услуг по организации проведения мероприятий по отлову и содержанию животных без владельцев\n"
        "на территории Липецкого муниципального округа"
    )
    subtitle.alignment = WD_ALIGN_PARAGRAPH.CENTER

    # Таблица 4 столбца
    table = doc.add_table(rows=2, cols=7)
    table.style = 'Table Grid'

    headers = [
        "Территориальный отдел",
        "Сведения о заявителе",
        "Место нахождения животного",
        "Сведения о количестве животных без владельцев",
        "Поведение",
        "Срочность",
        "Контактное лицо в территориальном отделе"
    ]

    for i, text in enumerate(headers):
        cell = table.cell(0, i)
        cell.text = text
        cell.paragraphs[0].alignment = WD_ALIGN_PARAGRAPH.CENTER

    data = [
        territorial_office,
        applicant,
        address,
        dogs_count,
        behavior,
        urgency,
        contact_person
    ]

    for i, text in enumerate(data):
        cell = table.cell(1, i)
        cell.text = text
        cell.paragraphs[0].alignment = WD_ALIGN_PARAGRAPH.CENTER

    # Контакты

    # Подпись
    doc.add_paragraph("\nЗаместитель председателя комитета энергетики и ЖКХ\nадминистрации Липецкого муниципального округа                                      Е.А. Челидина")
    doc.add_paragraph("34-57-29")

    # Сохраняем документ
    shared_dir = Path("/app/shared")
    doc.save(f"{shared_dir}/zayavka_{number}.docx")
    print(f"Документ создан: zayavka_{number}.docx")

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)



