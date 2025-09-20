
FROM python:3.13.7


WORKDIR /app



COPY main.py .
COPY requirements.txt .

RUN pip install -r requirements.txt



EXPOSE 8000

CMD ["python", "main.py"]