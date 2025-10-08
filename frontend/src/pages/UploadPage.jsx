import React, { useState } from "react";
import axios from "axios";

function UploadPage() {
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState("");

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
  };

  const handleUpload = async () => {
    if (!file) {
      setMessage("Please select a file first.");
      return;
    }

    if (file.type !== "application/pdf") {
      setMessage("Only PDF files are allowed.");
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      const endpoint= `http://${window.location.hostname}:8080`
      const response = await axios.post(
        `${endpoint}/api/document`,
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      );

      setMessage(response.data.message || "File uploaded successfully!");
    } catch (error) {
      console.error(error);
      setMessage("Failed to upload file.");
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", alignItems: "center", marginTop: "50px", gap: "20px" }}>
      <h2>Upload Page</h2>
      <input type="file" accept="application/pdf" onChange={handleFileChange} />
      <button onClick={handleUpload} style={{ padding: "10px 20px", fontSize: "16px" }}>
        Upload PDF
      </button>
      {message && <p>{message}</p>}
    </div>
  );
}

export default UploadPage;
