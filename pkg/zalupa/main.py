from openpyxl import load_workbook
from pathlib import Path
from fastapi import FastAPI
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from datetime import datetime
import uvicorn
import os
from typing import List, Dict, Any

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "API работает", "endpoints": ["/generate", "/generate-multiple"]}


@app.get("/health")
async def health():
    return {"status": "ok", "timestamp": datetime.now().isoformat()}


class RequestItem(BaseModel):
    source_name: str
    applicant_name: str
    address: str
    dogs_count: int
    behavior: str
    urgency: str
    contact_person: str


class RequestFull(BaseModel):
    number: str
    source_name: str
    applicant_name: str
    address: str
    dogs_count: int
    behavior: str
    urgency: str
    contact_person: str


class RequestMultiple(BaseModel):
    number: str
    requests: List[RequestItem]


def generate_single_document(
    template_path: str, output_path: str, data: Dict[str, Any]
):
    """
    Генерация Excel-документа на основе шаблона для одной заявки
    :param template_path: путь к шаблону (xlsx)
    :param output_path: путь для сохранения нового файла
    :param data: словарь с данными {ячейка: значение}
    """
    wb = load_workbook(template_path)
    ws = wb.active

    # Подставляем данные
    print(f"Всего данных для записи: {len(data)} элементов")
    for cell_address, value in data.items():
        print(f"Попытка записи в {cell_address} = {value} (тип: {type(value)})")
        try:
            # Находим все объединенные диапазоны
            merged_ranges = list(ws.merged_cells.ranges)
            target_cell = None

            # Проверяем, находится ли ячейка в объединенном диапазоне
            for merged_range in merged_ranges:
                if cell_address in merged_range:
                    min_col, min_row, max_col, max_row = merged_range.bounds
                    target_cell = ws.cell(row=min_row, column=min_col).coordinate
                    print(
                        f"Ячейка {cell_address} в объединенном диапазоне, используем {target_cell}"
                    )
                    break

            if target_cell:
                ws[target_cell] = value
                print(f"✓ Записано в объединенную ячейку {target_cell}: {value}")
            else:
                ws[cell_address] = value
                print(f"✓ Записано в ячейку {cell_address}: {value}")

        except Exception as e:
            print(f"✗ Ошибка при записи в ячейку {cell_address}: {e}")
            try:
                ws[cell_address].value = value
                print(f"✓ Записано напрямую в {cell_address}: {value}")
            except Exception as e2:
                print(f"✗ Не удалось записать в ячейку {cell_address}: {e2}")

    shared_dir = Path("/app/shared")
    wb.save(f"{shared_dir}/{output_path}")


def generate_multiple_sheets_document(
    template_path: str, output_path: str, number: str, requests: List[RequestItem]
):
    """
    Генерация Excel-документа с отдельными листами для каждой заявки
    :param template_path: путь к шаблону (xlsx)
    :param output_path: путь для сохранения нового файла
    :param number: номер заявки
    :param requests: список заявок
    """
    from copy import copy

    wb = load_workbook(template_path)
    template_sheet = wb.active

    # Удаляем все существующие листы кроме шаблона
    sheets_to_remove = []
    for sheet in wb.worksheets:
        if sheet != template_sheet:
            sheets_to_remove.append(sheet)
    for sheet in sheets_to_remove:
        wb.remove(sheet)

    # Создаем листы для каждой заявки
    for i, request_item in enumerate(requests):
        print(f"Создаем лист для заявки {i + 1}")

        if i == 0:
            # Первая заявка - используем существующий лист
            current_sheet = template_sheet
            current_sheet.title = f"Заявка {i + 1}"
        else:
            # Копируем шаблон для новых заявок
            current_sheet = wb.copy_worksheet(template_sheet)
            current_sheet.title = f"Заявка {i + 1}"

        # Данные для текущей заявки
        sheet_data = {
            "B15": f"Заявка № {number}-{i + 1}",
            "B12": datetime.now().strftime("%d.%m.%Y"),
            "B20": request_item.source_name,
            "D20": request_item.applicant_name,
            "E20": request_item.address,
            "F20": request_item.dogs_count,
            "G20": request_item.behavior,
            "H20": request_item.urgency,
            "I20": request_item.contact_person,
        }

        # Записываем данные на текущий лист
        for cell_address, value in sheet_data.items():
            try:
                # Проверяем объединенные ячейки
                merged_ranges = list(current_sheet.merged_cells.ranges)
                target_cell = None

                for merged_range in merged_ranges:
                    if cell_address in merged_range:
                        min_col, min_row, max_col, max_row = merged_range.bounds
                        target_cell = current_sheet.cell(
                            row=min_row, column=min_col
                        ).coordinate
                        break

                if target_cell:
                    current_sheet[target_cell] = value
                    print(f"✓ Записано в объединенную ячейку {target_cell}: {value}")
                else:
                    current_sheet[cell_address] = value
                    print(f"✓ Записано в ячейку {cell_address}: {value}")

            except Exception as e:
                print(f"✗ Ошибка при записи в ячейку {cell_address}: {e}")
                try:
                    current_sheet[cell_address].value = value
                    print(f"✓ Записано напрямую в {cell_address}: {value}")
                except Exception as e2:
                    print(f"✗ Не удалось записать в ячейку {cell_address}: {e2}")

    shared_dir = Path(".")
    wb.save(f"{shared_dir}/{output_path}")
    print(f"Документ с {len(requests)} листами создан: {output_path}")


@app.post("/generate")
async def generate_doc(req: RequestFull):
    try:
        print(f"Получен запрос: {req}")
        filename = f"zayavka_{req.number}_{datetime.now().year}.xlsx"
        template = "tamplate.xlsx"

        # Проверяем существование шаблона
        if not os.path.exists(template):
            return JSONResponse(
                {"error": f"Template file {template} not found"}, status_code=500
            )

        data = {
            "B15": f"Заявка № {req.number}",
            "B12": datetime.now().strftime("%d.%m.%Y"),
            "B20": req.source_name,
            "D20": req.applicant_name,
            "E20": req.address,
            "F20": req.dogs_count,
            "G20": req.behavior,
            "H20": req.urgency,
            "I20": req.contact_person,
        }

        generate_single_document(template, filename, data)
        print("Документ успешно создан:", filename)
        shared_dir = Path("/app/shared")
        if os.path.exists(f"{shared_dir}/{filename}"):
            return JSONResponse({"status": "ok", "filename": filename}, status_code=200)
        return JSONResponse({"error": "File not created"}, status_code=500)
    except Exception as e:
        print(f"Ошибка при создании документа: {e}")
        return JSONResponse({"error": str(e)}, status_code=500)


@app.post("/generate-multiple")
async def generate_multiple_doc(req: RequestMultiple):
    try:
        print(f"Получен множественный запрос: {req}")

        if not req.requests:
            return JSONResponse({"error": "No requests provided"}, status_code=400)

        filename = f"zayavka_{req.number}_{datetime.now().year}.xlsx"
        template = "tamplate.xlsx"

        # Проверяем существование шаблона
        if not os.path.exists(template):
            return JSONResponse(
                {"error": f"Template file {template} not found"}, status_code=500
            )

        # Создаем документ с отдельными листами
        generate_multiple_sheets_document(template, filename, req.number, req.requests)
        print(f"Документ с {len(req.requests)} заявками успешно создан:", filename)
        shared_dir = Path("/app/shared")
        if os.path.exists(f"{shared_dir}/{filename}"):
            return JSONResponse(
                {
                    "status": "ok",
                    "requests_count": len(req.requests),
                    "filename": filename,
                },
                status_code=200,
            )
        return JSONResponse({"error": "File not created"}, status_code=500)
    except Exception as e:
        print(f"Ошибка при создании документа: {e}")
        return JSONResponse({"error": str(e)}, status_code=500)


if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)
