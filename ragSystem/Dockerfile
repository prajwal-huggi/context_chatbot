# Use an official Python base image
FROM python:3.12-slim

# Set working directory inside the container
WORKDIR /app

ARG VERSION=ragsystem-v7

# Copy requirements first for caching
COPY ./requirementsPython.txt ./

# Install dependencies
RUN pip install --no-cache-dir --upgrade pip \
    && pip install --no-cache-dir -r requirementsPython.txt

# Copy the rest of the application code
COPY . .

# Expose the port Flask will run on
EXPOSE 8081

# Set environment variables (optional, can also use .env)
# ENV FLASK_APP=app.py
# ENV FLASK_RUN_HOST=0.0.0.0
# ENV PYTHONUNBUFFERED=1

# Run the Flask app
CMD ["uvicorn", "main_rag:app", "--host", "0.0.0.0", "--port", "8081", "--reload"]
