from fastapi import FastAPI, UploadFile, File
import shutil
from typing import Optional

from langchain_google_genai import ChatGoogleGenerativeAI, GoogleGenerativeAIEmbeddings
from langchain.prompts import PromptTemplate
from langchain.text_splitter import RecursiveCharacterTextSplitter

from langchain_core.runnables import RunnableParallel, RunnablePassthrough, RunnableLambda

from langchain.output_parsers import PydanticOutputParser
from pydantic import BaseModel, Field

from langchain_chroma import Chroma
import chromadb

from dotenv import load_dotenv
import os

from PyPDF2 import PdfReader

app= FastAPI()

load_dotenv()

# Creating the embeddings
embeddings= GoogleGenerativeAIEmbeddings(
    # model="gemini-embedding-001"
    model="models/gemini-embedding-exp-03-07",   
)

#Creating the LLM
llm= ChatGoogleGenerativeAI(
    model="gemini-2.5-flash-lite",
    temperature=0.3,
    max_tokens=500,
    timeout=None,
    max_retries=2,
    api_key= os.getenv("GOOGLE_API_KEY")
)

#Creating the vector store chromaDB
dbPath= os.path.join(os.getcwd(), 'pdfDB')

client= chromadb.Client(
    settings=chromadb.config.Settings(is_persistent= True, persist_directory= dbPath)
)

collections= client.get_or_create_collection("Pdfs")

# Building the Retriever to get the relevant documents
vectorestore= Chroma(
    client= client,
    collection_name="Pdfs",
    embedding_function=embeddings
)

retriever= vectorestore.as_retriever(
    search_type= "mmr",
    search_kwargs= {"k": 3, "fetch_k": 10}
)

# Making the parser
class RAGAnswer(BaseModel):
    mode: str = Field(description="Either 'context+reasoning' or 'reasoning'")
    answer: str = Field(description="The answer to the question")
    
parser= PydanticOutputParser(pydantic_object=RAGAnswer)

# Building the parser
prompt= PromptTemplate(
    template="""
    You are a helpful assistant. Decide how to answer based on the context:

    - If the provided context is useful, combine it with your reasoning. 
    In that case set "mode" = "context+reasoning".
    - If the context is missing or irrelevant, rely only on your reasoning. 
    In that case set "mode" = "reasoning".
    - Always return valid JSON.

    Context:
    {context}

    Question:
    {question}

    {format_instructions}
    """,
    input_variables= ["context", "question"],
    partial_variables={"format_instructions": parser.get_format_instructions()}   
)

def format_docs(retrieved_docs):
    context_text= "\n\n".join(doc.page_content for doc in retrieved_docs)
    return context_text

parallel_chain= RunnableParallel({
    "question": RunnablePassthrough(),
    "context": retriever | RunnableLambda(format_docs)
})

final_chain= parallel_chain | prompt | llm | parser

@app.get('/')
async def greetings():
    print("Welcome to ragSystem")
    return {"message":"Welcome to ragSystem"}

# Fast API endpoints
@app.post('/add_pdf')
async def add_pdf(file: UploadFile= File(...)):
    if file.content_type!= "application/pdf":
        return {"error":"Only PDFs are supported"}
    
    temp_path= os.path.join("/tmp", file.filename)
    with open(temp_path, "wb") as f:
        shutil.copyfileobj(file.file, f)
        
    reader= PdfReader(temp_path)
    pdfContent= [page.extract_text() for page in reader.pages]
    
    splitter = RecursiveCharacterTextSplitter(chunk_size=500, chunk_overlap=100)
    chunks = splitter.create_documents(pdfContent)
    
    collections.add(
        ids=[f"doc_{i}" for i in range(len(chunks))],
        documents=[i.page_content for i in chunks],
        embeddings=embeddings.embed_documents([i.page_content for i in chunks])
    )
    
    return {"message":f"{file.filename} added successfully"}

class QueryRequest(BaseModel):
    question: str
    
@app.post("/query")
async def query_rag(req: QueryRequest):
    result = final_chain.invoke(req.question)
    final_result = result.model_dump()
    return final_result

@app.post("/reset")
async def reset_rag():
    global retriever, final_chain, vectorestore, collections, parallel_chain
    
    # recreate empty collection
    client.delete_collection("Pdfs")
    collections = client.get_or_create_collection("Pdfs")
    
    vectorestore = Chroma(
        client=client,
        collection_name="Pdfs",
        embedding_function=embeddings
    )
    retriever = vectorestore.as_retriever(
        search_type="mmr",
        search_kwargs={"k": 3, "fetch_k": 10}
    )
    
    # Rebuild chain so it points to the fresh retriever
    parallel_chain = RunnableParallel({
        "question": RunnablePassthrough(),
        "context": retriever | RunnableLambda(format_docs)
    })
    final_chain = parallel_chain | prompt | llm | parser

    return {"message": "RAG database reset successfully"}
