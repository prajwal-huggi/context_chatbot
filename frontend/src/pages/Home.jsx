import React from "react";
import { useNavigate } from "react-router-dom";

function Home() {
  const navigate = useNavigate();

  const goToUpload = () => {
    navigate("/upload"); // Redirect to Upload Document page
  };

  const goToChat = () => {
    navigate("/chat"); // Redirect to Chat page
  };

  return (
    <div style={{ display: "flex", gap: "20px", justifyContent: "center", marginTop: "50px" }}>
      <button onClick={goToUpload} style={{ padding: "10px 20px", fontSize: "16px" }}>
        Upload Document
      </button>
      <button onClick={goToChat} style={{ padding: "10px 20px", fontSize: "16px" }}>
        Chat
      </button>
    </div>
  );
}

export default Home;
