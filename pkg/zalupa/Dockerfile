
FROM python:3.13.7



WORKDIR /app



COPY main.py .
COPY requirements.txt .
COPY template.xlsx .
COPY templateV1.xlsx .

RUN pip install -r requirements.txt



EXPOSE 8001

CMD ["python", "main.py"]
