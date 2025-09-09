from docx import Document
from docx.shared import Pt
from docx.enum.text import WD_ALIGN_PARAGRAPH
from docx.enum.section import WD_ORIENT

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

doc.add_paragraph("\n18.08.2025 №______", style='Normal')

# Заголовок
title = doc.add_paragraph("Заявка №190")
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
    "Введенский территориальный отдел",
    "И.о. начальника Введенского территориального отдела И.В.Денисов",
    "Липецкий округ, с.Воскресеновка, ул.Наречная в районе д.53",
    "1", 
    "Агрессивное", 
    "Срочно",
    "И.о. начальника Введенского территориального отдела Денисов И.В. т. 75-60-21"
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
doc.save("zayavka.docx")
print("Документ создан: zayavka.docx")
